package repo

import (
	"context"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
)

type ArticleRepo struct {
}

func (a *ArticleRepo) Create(ctx context.Context, article *entity.Article) error    {}
func (a *ArticleRepo) Editor(ctx context.Context, article *entity.Article) error    {}
func (a *ArticleRepo) Delete(ctx context.Context, articleID, authorID uint64) error {}

func (a *ArticleRepo) PublishDraft(ctx context.Context, articleID uint64) error {}
func (a *ArticleRepo) DeleteDraft(ctx context.Context, articleID uint64) error  {}

func (a *ArticleRepo) GetArticle(ctx context.Context, articleID uint64) (*aggregate.ArticleAggregate, error) {
}
func (a *ArticleRepo) ListArticles(ctx context.Context, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
}
func (a *ArticleRepo) ListUserArticles(ctx context.Context, AuthorID uint64, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
}
func (a *ArticleRepo) ListMyArticles(ctx context.Context, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
}
func (a *ArticleRepo) ListByCategory(ctx context.Context, categoryID uint64, page, size uint32) ([]*aggregate.ArticleAggregate, error) {
}

func (a *ArticleRepo) SearchArticles(ctx context.Context, keyword string, page, size uint32) ([]*dto.SearchArticlesResponse, error) {
}
