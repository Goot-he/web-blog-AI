package repo

import (
	"context"
	"time"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/article"
	"core-server/internal/infras/repo/ent/comment"
	"core-server/internal/infras/repo/ent/predicate"
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

// Create 新建一级评论（校验文章，写入评论并把文章评论数 +1）
func (r *CommentRepo) Create(ctx context.Context, data *entity.Comment) error {
	return r.WithTransaction(ctx, func(ctx context.Context) error {
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

		*data = *toEntityComment(created)
		return nil
	})
}

// CreateReply 新建回复（校验根评论和被回复评论，更新文章评论数和根评论子评论数）
func (r *CommentRepo) CreateReply(ctx context.Context, data *entity.Comment) error {
	return r.WithTransaction(ctx, func(ctx context.Context) error {
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

		// 文章评论数 +1，根评论的子评论数 +1
		if err = r.DB(ctx).Article.UpdateOneID(data.ArticleID).
			AddCommentCount(1).
			Exec(ctx); err != nil {
			return err
		}
		if err = r.DB(ctx).Comment.UpdateOneID(root.ID).
			AddChildCount(1).
			Exec(ctx); err != nil {
			return err
		}

		*data = *toEntityComment(created)
		return nil
	})
}

// Delete 软删除评论（评论作者本人或文章作者可删，一级评论会连同回复一起删）
func (r *CommentRepo) Delete(ctx context.Context, commentID, userID, articleID uint64) error {
	return r.WithTransaction(ctx, func(ctx context.Context) error {
		target, err := r.DB(ctx).Comment.Query().
			Where(
				comment.IDEQ(commentID),
				comment.ArticleIDEQ(articleID),
				comment.DeletedAtIsNil(),
			).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return ErrorCommentNotFound
			}
			return err
		}

		// 评论作者本人 或 文章作者 才有权限删除
		articleRow, err := r.DB(ctx).Article.Query().
			Where(article.IDEQ(articleID), article.DeletedAtIsNil()).
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
			if err = r.DB(ctx).Comment.UpdateOneID(target.RootID).
				AddChildCount(-1).
				Exec(ctx); err != nil {
				return err
			}
		}

		if err = r.DB(ctx).Comment.UpdateOneID(target.ID).
			SetDeletedAt(now).
			Exec(ctx); err != nil {
			return err
		}

		// 文章评论数减去被删除的数量
		return r.DB(ctx).Article.UpdateOneID(articleID).
			AddCommentCount(-deletedCount).
			Exec(ctx)
	})
}

// ListByArticle 分页查询文章下的一级评论（按时间倒序，带评论用户信息）
func (r *CommentRepo) ListByArticle(ctx context.Context, articleID uint64, page, size uint32) ([]*aggregate.CommentAggregate, error) {
	return r.list(ctx, page, size,
		[]predicate.Comment{
			comment.ArticleIDEQ(articleID),
			comment.RootIDEQ(0),
			comment.DeletedAtIsNil(),
		},
		comment.ByCreatedAt(sql.OrderDesc()),
	)
}

// ListByRoot 分页查询根评论下的回复（按时间正序，带评论用户信息）
func (r *CommentRepo) ListByRoot(ctx context.Context, rootID uint64, page, size uint32) ([]*aggregate.CommentAggregate, error) {
	return r.list(ctx, page, size,
		[]predicate.Comment{
			comment.RootIDEQ(rootID),
			comment.DeletedAtIsNil(),
		},
		comment.ByCreatedAt(sql.OrderAsc()),
	)
}

// list 评论列表的公共实现（分页、排序、预加载用户、转聚合）
func (r *CommentRepo) list(ctx context.Context, page, size uint32, preds []predicate.Comment, order comment.OrderOption) ([]*aggregate.CommentAggregate, error) {
	limit, offset := pageLimitOffset(page, size)
	comments, err := r.DB(ctx).Comment.Query().
		Where(preds...).
		Order(order).
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
