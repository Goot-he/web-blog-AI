package domain

import (
	"context"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
)

type ArticleRepo interface {
	Create(ctx context.Context, article *entity.Article) error
	Editor(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, articleID, authorID uint64) error

	PublishDraft(ctx context.Context, articleID uint64) error
	DeleteDraft(ctx context.Context, articleID uint64) error

	GetArticle(ctx context.Context, articleID uint64) (*aggregate.ArticleAggregate, error)
	ListArticles(ctx context.Context, page, size uint32) ([]*aggregate.ArticleAggregate, error)
	ListUserArticles(ctx context.Context, AuthorID uint64, page, size uint32) ([]*aggregate.ArticleAggregate, error)
	ListMyArticles(ctx context.Context, page, size uint32) ([]*aggregate.ArticleAggregate, error)
	ListByCategory(ctx context.Context, categoryID uint64, page, size uint32) ([]*aggregate.ArticleAggregate, error)

	SearchArticles(ctx context.Context, keyword string, page, size uint32) ([]*dto.SearchArticlesResponse, error)
}

type ArticleDomainRepo interface {
	ArticleRepo  // 文章接口
	ITransaction // 事务通用接口
}
