package repo

import (
	"context"
	"time"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/article"
	"core-server/internal/infras/repo/ent/comment"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"

	"entgo.io/ent/dialect/sql"
)

type CommentRepo struct {
	*EntClient
}

func NewCommentRepo(client *EntClient) *CommentRepo {
	return &CommentRepo{client}
}

// CreateComment 新建一级评论（校验文章，文章评论数 +1），返回评论 ID
func (r *CommentRepo) CreateComment(ctx context.Context, data *entity.Comment) (uint64, error) {
	var id uint64
	err := r.WithTransaction(ctx, func(ctx context.Context) error {
		if err := r.ensureArticleExists(ctx, data.ArticleID); err != nil {
			return err
		}

		created, err := r.DB(ctx).Comment.Create().
			SetArticleID(data.ArticleID).
			SetPublishUserID(data.PublishUserID).
			SetContent(data.Comment).
			Save(ctx)
		if err != nil {
			return err
		}

		if err = r.DB(ctx).Article.UpdateOneID(data.ArticleID).
			AddCommentCount(1).
			Exec(ctx); err != nil {
			return err
		}

		id = created.ID
		*data = *toEntityComment(created)
		return nil
	})
	return id, err
}

// CreateReply 新建回复（校验根评论和被回复评论，文章评论数 +1、根评论子评论数 +1），返回回复 ID
func (r *CommentRepo) CreateReply(ctx context.Context, data *entity.Comment) (uint64, error) {
	var id uint64
	err := r.WithTransaction(ctx, func(ctx context.Context) error {
		if err := r.ensureArticleExists(ctx, data.ArticleID); err != nil {
			return err
		}

		// 根评论必须存在且是一级评论
		root, err := r.DB(ctx).Comment.Query().
			Where(
				comment.IDEQ(data.RootID),
				comment.ArticleIDEQ(data.ArticleID),
				comment.RootIDEQ(0),
				comment.DeletedAtIsNil(),
			).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return ErrorCommentNotFound
			}
			return err
		}

		// 被回复的评论必须存在，且属于同一篇文章、同一个根评论
		replyTo, err := r.DB(ctx).Comment.Query().
			Where(
				comment.IDEQ(data.ReplyToID),
				comment.ArticleIDEQ(data.ArticleID),
				comment.DeletedAtIsNil(),
			).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return ErrorCommentNotFound
			}
			return err
		}
		if replyTo.ID != root.ID && replyTo.RootID != root.ID {
			return ErrorInvalidParam
		}

		created, err := r.DB(ctx).Comment.Create().
			SetArticleID(data.ArticleID).
			SetPublishUserID(data.PublishUserID).
			SetRootID(root.ID).
			SetReplyID(replyTo.ID).
			SetReplyUserID(replyTo.PublishUserID).
			SetContent(data.Comment).
			Save(ctx)
		if err != nil {
			return err
		}

		// 文章评论数 +1
		if err = r.DB(ctx).Article.UpdateOneID(data.ArticleID).
			AddCommentCount(1).
			Exec(ctx); err != nil {
			return err
		}
		// 根评论的子评论数 +1
		if err = r.DB(ctx).Comment.UpdateOneID(root.ID).
			AddChildCount(1).
			Exec(ctx); err != nil {
			return err
		}

		id = created.ID
		*data = *toEntityComment(created)
		return nil
	})
	return id, err
}

// SoftDelete 软删除评论（评论作者或文章作者可删；一级评论连同回复一起删）
func (r *CommentRepo) SoftDelete(ctx context.Context, id, userID uint64) error {
	return r.WithTransaction(ctx, func(ctx context.Context) error {
		target, err := r.DB(ctx).Comment.Query().
			Where(comment.IDEQ(id), comment.DeletedAtIsNil()).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return ErrorCommentNotFound
			}
			return err
		}

		// 评论作者本人 或 文章作者 才有权限删除
		articleRow, err := r.DB(ctx).Article.Query().
			Where(article.IDEQ(target.ArticleID), article.DeletedAtIsNil()).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return ErrorArticleNotFound
			}
			return err
		}
		if userID != target.PublishUserID && userID != articleRow.AuthorID {
			return ErrorNoPermission
		}

		now := time.Now()
		deletedCount := int64(1)

		if target.RootID == 0 {
			// 一级评论：把它的回复一起软删除
			replies, err := r.DB(ctx).Comment.Update().
				Where(comment.RootIDEQ(target.ID), comment.DeletedAtIsNil()).
				SetDeletedAt(now).
				Save(ctx)
			if err != nil {
				return err
			}
			deletedCount += int64(replies)
		} else {
			// 回复：根评论的子评论数 -1
			if err = r.decChildCount(ctx, target.RootID); err != nil {
				return err
			}
		}

		if err = r.DB(ctx).Comment.UpdateOneID(target.ID).
			SetDeletedAt(now).
			Exec(ctx); err != nil {
			return err
		}

		// 文章评论数减去被删除的数量
		return r.DB(ctx).Article.UpdateOneID(target.ArticleID).
			AddCommentCount(-deletedCount).
			Exec(ctx)
	})
}

// SoftDeleteRepliesByRoot 软删除某个根评论下的全部回复，返回删除数量
func (r *CommentRepo) SoftDeleteRepliesByRoot(ctx context.Context, rootID uint64) (int64, error) {
	var deleted int64
	err := r.WithTransaction(ctx, func(ctx context.Context) error {
		root, err := r.DB(ctx).Comment.Query().
			Where(comment.IDEQ(rootID), comment.DeletedAtIsNil()).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return ErrorCommentNotFound
			}
			return err
		}

		n, err := r.DB(ctx).Comment.Update().
			Where(comment.RootIDEQ(rootID), comment.DeletedAtIsNil()).
			SetDeletedAt(time.Now()).
			Save(ctx)
		if err != nil {
			return err
		}
		deleted = int64(n)

		if deleted > 0 {
			if err = r.DB(ctx).Article.UpdateOneID(root.ArticleID).
				AddCommentCount(-deleted).
				Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	return deleted, err
}

// GetByID 查询单条评论（带评论用户信息）
func (r *CommentRepo) GetByID(ctx context.Context, id uint64) (*aggregate.CommentAggregate, error) {
	data, err := r.DB(ctx).Comment.Query().
		Where(comment.IDEQ(id), comment.DeletedAtIsNil()).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrorCommentNotFound
		}
		return nil, err
	}
	return toCommentAggregate(data), nil
}

// ListByIDs 按 ID 批量查询评论（带评论用户信息）
func (r *CommentRepo) ListByIDs(ctx context.Context, ids []uint64) ([]*aggregate.CommentAggregate, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	comments, err := r.DB(ctx).Comment.Query().
		Where(comment.IDIn(ids...), comment.DeletedAtIsNil()).
		Order(comment.ByID(sql.OrderAsc())).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toCommentAggregates(comments), nil
}

// IncrementChildCount 根评论的子评论数 +1
func (r *CommentRepo) IncrementChildCount(ctx context.Context, rootID uint64) error {
	return r.DB(ctx).Comment.UpdateOneID(rootID).AddChildCount(1).Exec(ctx)
}

// DecrementChildCount 根评论的子评论数 -1（已经为 0 时不处理）
func (r *CommentRepo) DecrementChildCount(ctx context.Context, rootID uint64) error {
	return r.decChildCount(ctx, rootID)
}

// decChildCount 子评论数 -1，用条件更新避免无符号字段下溢
func (r *CommentRepo) decChildCount(ctx context.Context, rootID uint64) error {
	_, err := r.DB(ctx).Comment.Update().
		Where(comment.IDEQ(rootID), comment.ChildCountGT(0)).
		AddChildCount(-1).
		Save(ctx)
	return err
}

// ListTopByArticle 分页查询文章下的一级评论（按时间倒序）
func (r *CommentRepo) ListTopByArticle(ctx context.Context, articleID uint64, offset, limit int) ([]*aggregate.CommentAggregate, error) {
	offset, limit = normalizeOffsetLimit(offset, limit)
	comments, err := r.DB(ctx).Comment.Query().
		Where(
			comment.ArticleIDEQ(articleID),
			comment.RootIDEQ(0),
			comment.DeletedAtIsNil(),
		).
		Order(comment.ByCreatedAt(sql.OrderDesc())).
		Offset(offset).
		Limit(limit).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toCommentAggregates(comments), nil
}

// ListRepliesByRoot 分页查询根评论下的回复（按时间正序）
func (r *CommentRepo) ListRepliesByRoot(ctx context.Context, rootID uint64, offset, limit int) ([]*aggregate.CommentAggregate, error) {
	offset, limit = normalizeOffsetLimit(offset, limit)
	comments, err := r.DB(ctx).Comment.Query().
		Where(
			comment.RootIDEQ(rootID),
			comment.DeletedAtIsNil(),
		).
		Order(comment.ByCreatedAt(sql.OrderAsc())).
		Offset(offset).
		Limit(limit).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toCommentAggregates(comments), nil
}

// ensureArticleExists 校验文章存在且已发布
func (r *CommentRepo) ensureArticleExists(ctx context.Context, articleID uint64) error {
	exists, err := r.DB(ctx).Article.Query().
		Where(
			article.IDEQ(articleID),
			article.DeletedAtIsNil(),
			article.PublishedAtNotNil(),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return ErrorArticleNotFound
	}
	return nil
}
