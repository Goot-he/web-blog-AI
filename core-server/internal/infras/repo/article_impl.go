package repo

import (
	"context"
	"core-server/internal/infras/repo/ent/userstatu"
	"core-server/internal/model/enum"
	"time"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/article"
	"core-server/internal/infras/repo/ent/category"
	"core-server/internal/infras/repo/ent/predicate"
	"core-server/internal/infras/repo/ent/user"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"

	"entgo.io/ent/dialect/sql"
)

type ArticleRepo struct {
	*EntClient
}

// NewArticleRepo 创建文章仓库实例
func NewArticleRepo(client *EntClient) *ArticleRepo {
	return &ArticleRepo{
		client,
	}
}

// Create 新建文章（校验作者和分类，已发布时作者文章数加一，整体在事务中执行）
func (a *ArticleRepo) Create(ctx context.Context, data *entity.Article) error {
	if err := a.WithTransaction(ctx, func(ctx context.Context) error {
		// 先校验作者是否在用户表里面
		if err := a.ensureAuthorExists(ctx, data.AuthorID); err != nil {
			return ErrorUserNotFound
		}
		// 分类是必填的 且必须真实存在 否则插入时会撞上外键约束
		if err := a.ensureCategoryExists(ctx, data.CategoryID); err != nil {
			return ErrorCategoryNotFound
		}

		// 将数据存入到sql表
		create := a.DB(ctx).Article.Create().
			SetTitle(data.Title).
			SetSummary(data.Summary).
			SetContent(data.Content).
			SetCoverage(data.CoverImage).
			SetAuthorID(data.AuthorID).
			SetCategoryID(data.CategoryID).
			SetIsTop(data.IsTop).
			SetIsPublished(data.IsPublished)
		if data.PublishedAt.Valid {
			create.SetPublishedAt(data.PublishedAt.Time)
		}

		// 同时用户状态文章计数加一
		if data.IsPublished {
			err := a.DB(ctx).UserStatu.Update().Where(
				userstatu.UserIDEQ(data.AuthorID), userstatu.DeletedAtIsNil()).
				AddArticleCount(1).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		_, err := create.Save(ctx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Editor 编辑自己名下的文章（可改分类、置顶和发布状态，并按发布状态变化同步作者文章计数 注意：已经发布的文章不可以在撤销为草稿）
func (a *ArticleRepo) Editor(ctx context.Context, data *entity.Article) error {
	if err := a.WithTransaction(ctx, func(ctx context.Context) error {
		// 只能编辑自己名下未被软删除的文章
		exists, err := a.DB(ctx).Article.Query().Where(
			article.IDEQ(data.ID),
			article.AuthorIDEQ(data.AuthorID),
			article.DeletedAtIsNil(),
		).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return ErrorArticleNotFound
			}
			return err
		}
		// 允许修改分类 但新分类必须存在
		if err := a.ensureCategoryExists(ctx, data.CategoryID); err != nil {
			return err
		}

		// 作者不能被编辑修改 所以这里不 SetAuthorID
		update := a.DB(ctx).Article.UpdateOne(exists).
			SetTitle(data.Title).
			SetSummary(data.Summary).
			SetContent(data.Content).
			SetCoverage(data.CoverImage).
			SetCategoryID(data.CategoryID).
			SetIsTop(data.IsTop).
			SetIsPublished(data.IsPublished)

		if data.PublishedAt.Valid {
			update.SetPublishedAt(data.PublishedAt.Time)
		}

		if _, err = update.Save(ctx); err != nil {
			return err
		}

		if exists.IsPublished == false && data.IsPublished == true {
			err := a.DB(ctx).UserStatu.Update().
				Where(userstatu.UserIDEQ(data.AuthorID), userstatu.DeletedAtIsNil()).
				AddArticleCount(1).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Delete 软删除自己名下的文章
func (a *ArticleRepo) Delete(ctx context.Context, articleID, authorID uint64, roleID enum.Role) error {
	if err := a.WithTransaction(ctx, func(ctx context.Context) error {
		// 软删除 条件里带上作者 保证只能删自己的文章
		update := a.DB(ctx).Article.Update().
			Where(article.IDEQ(articleID), article.DeletedAtIsNil()).
			SetDeletedAt(time.Now())

		if roleID != enum.RoleAdmin {
			// 如果不是管理员 加上是否是作者本人校验
			update.Where(article.AuthorIDEQ(authorID))
		}

		affected, err := update.Save(ctx)
		if err != nil {
			return err
		}
		if affected == 0 {
			// 文章不存在、不是自己的或者已经被删除
			return ErrorArticleNotFound
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// PublishDraft 发布草稿（校验文章是草稿后写入发布时间）
func (a *ArticleRepo) PublishDraft(ctx context.Context, articleID, authorID uint64) error {
	if err := a.WithTransaction(ctx, func(ctx context.Context) error {
		exists, err := a.queryOwnArticle(ctx, articleID, authorID)
		if err != nil {
			return ErrorDraftNotFound
		}
		if exists.PublishedAt != nil {
			// 已经发布过的文章不能再发布一次
			return ErrorArticleIsPublished
		}

		err = a.DB(ctx).Article.UpdateOne(exists).
			SetIsPublished(true).
			SetPublishedAt(time.Now()).
			Exec(ctx)

		return err
	}); err != nil {
		return err
	}
	return nil
}

// DeleteDraft 删除草稿（只允许删除自己名下未发布的文章）
func (a *ArticleRepo) DeleteDraft(ctx context.Context, articleID, authorID uint64) error {
	if err := a.WithTransaction(ctx, func(ctx context.Context) error {
		exists, err := a.queryOwnArticle(ctx, articleID, authorID)
		if err != nil {
			return ErrorDraftNotFound
		}
		if exists.PublishedAt != nil {
			// 已发布的文章要走删除文章 不能当成草稿删
			return ErrorArticleIsPublished
		}

		err = a.DB(ctx).Article.UpdateOne(exists).SetDeletedAt(time.Now()).Exec(ctx)
		return err
	}); err != nil {
		return err
	}
	return nil
}

// GetArticle 查询单篇已发布文章详情（带作者和分类）
func (a *ArticleRepo) GetArticle(ctx context.Context, articleID uint64) (*aggregate.ArticleAggregate, error) {
	// 对外只暴露未删除的已发布文章 草稿属于作者私有数据
	data, err := a.DB(ctx).Article.Query().
		Where(
			article.IDEQ(articleID),
			article.DeletedAtIsNil(),
			article.PublishedAtNotNil(),
		).
		WithUser().
		WithCategory().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrorArticleNotFound
		}
		return nil, err
	}
	return toArticleAggregate(data), nil
}

// ListArticles 分页查询全部已发布文章
func (a *ArticleRepo) ListArticles(ctx context.Context, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
	return a.listPublished(ctx, page, size)
}

// ListUserArticles 分页查询指定用户已发布的文章
func (a *ArticleRepo) ListUserArticles(ctx context.Context, authorID uint64, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
	return a.listPublished(ctx, page, size, article.AuthorIDEQ(authorID))
}

// ListMyArticles 分页查询我的文章（isAll为true返回全部，为false只返回草稿箱）
func (a *ArticleRepo) ListMyArticles(ctx context.Context, authorID uint64, isAll bool, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
	preds := []predicate.Article{
		article.AuthorIDEQ(authorID),
		article.DeletedAtIsNil(),
	}
	if !isAll {
		preds = append(preds, article.PublishedAtIsNil())
	}
	return a.list(ctx, page, size, preds, article.ByCreatedAt(sql.OrderDesc()))
}

// ListByCategory 分页查询指定分类下的已发布文章
func (a *ArticleRepo) ListByCategory(ctx context.Context, categoryID uint64, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
	return a.listPublished(ctx, page, size, article.CategoryIDEQ(categoryID))
}

// SearchArticles 按关键词搜索已发布文章（标题、概述、正文任一命中即可）
func (a *ArticleRepo) SearchArticles(ctx context.Context, keyword string, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
	// 正文是longtext 这里用的是 LIKE %keyword% 数据量大了要换成全文索引
	return a.listPublished(ctx, page, size,
		article.Or(
			article.TitleContainsFold(keyword),
			article.SummaryContainsFold(keyword),
			article.ContentContainsFold(keyword),
			// TitleContainsFold 是 ent 为 title 字段生成的谓词函数（命名规则：字段名 + 谓词名）；
			// Contains = SQL 的 LIKE '%值%'；
			// Fold = 忽略大小写（case-folding）；
			// article.Or(...) 不是查询函数，是逻辑组合器：把三个谓词用 OR 连起来。
		),
	)
}

// listPublished 分页查询已发布文章（置顶优先，其余按发布时间倒序）
func (a *ArticleRepo) listPublished(ctx context.Context, page, size uint32, preds ...predicate.Article) ([]*aggregate.ArticleAggregate, error) {
	preds = append(preds, article.DeletedAtIsNil(), article.PublishedAtNotNil())
	return a.list(ctx, page, size, preds,
		article.ByIsTop(sql.OrderDesc()),
		article.ByPublishedAt(sql.OrderDesc()),
	)
}

// list 统一处理分页、排序、预加载作者和分类，并转成聚合返回
func (a *ArticleRepo) list(ctx context.Context, page, size uint32, preds []predicate.Article, orders ...article.OrderOption) ([]*aggregate.ArticleAggregate, error) {
	limit, offset := pageLimitOffset(page, size)
	articles, err := a.DB(ctx).Article.Query().
		Where(preds...).
		Order(orders...).
		Offset(offset).
		Limit(limit).
		WithUser().
		WithCategory().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toArticleAggregates(articles), nil
}

// ensureAuthorExists 校验作者存在且未被删除
func (a *ArticleRepo) ensureAuthorExists(ctx context.Context, authorID uint64) error {
	exists, err := a.DB(ctx).User.Query().
		Where(user.IDEQ(authorID), user.DeletedAtIsNil()).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return ErrorUserNotFound
	}
	return nil
}

// ensureCategoryExists 校验分类存在且未被删除
func (a *ArticleRepo) ensureCategoryExists(ctx context.Context, categoryID uint64) error {
	exists, err := a.DB(ctx).Category.Query().
		Where(category.IDEQ(categoryID), category.DeletedAtIsNil()).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return ErrorCategoryNotFound
	}
	return nil
}

// queryOwnArticle 查询作者自己名下未被软删除的文章
func (a *ArticleRepo) queryOwnArticle(ctx context.Context, articleID, authorID uint64) (*ent.Article, error) {
	data, err := a.DB(ctx).Article.Query().Where(
		article.IDEQ(articleID),
		article.AuthorIDEQ(authorID),
		article.DeletedAtIsNil(),
	).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return data, nil
}
