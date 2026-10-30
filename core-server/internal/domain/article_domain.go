package domain

import (
	"context"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type ArticleRepo interface {
	Create(ctx context.Context, data *entity.Article) error
	Editor(ctx context.Context, data *entity.Article) error
	Delete(ctx context.Context, articleID, authorID uint64, roleID enum.UserRole) error

	PublishDraft(ctx context.Context, articleID, authorID uint64) error
	DeleteDraft(ctx context.Context, articleID, authorID uint64) error

	GetArticle(ctx context.Context, articleID uint64) (*aggregate.ArticleAggregate, error)
	ListArticles(ctx context.Context, page, size uint32) ([]*aggregate.ArticleAggregate, error)
	ListByIds(ctx context.Context, articleIDs []uint64) ([]*aggregate.ArticleAggregate, error)
	ListUserArticles(ctx context.Context, authorID uint64, page, size uint32) ([]*aggregate.ArticleAggregate, error)
	// ListMyArticles isPublished 为 nil 查全部，true 只查已发布，false 只查草稿
	ListMyArticles(ctx context.Context, authorID uint64, isPublished *bool, page, size uint32) ([]*aggregate.ArticleAggregate, error)
	ListByCategory(ctx context.Context, categoryID uint64, page, size uint32) ([]*aggregate.ArticleAggregate, error)

	SearchArticles(ctx context.Context, keyword string, page, size uint32) ([]*aggregate.ArticleAggregate, error)
}

type ArticleDomainRepo interface {
	ArticleRepo  // 文章接口
	ITransaction // 事务通用接口
}
