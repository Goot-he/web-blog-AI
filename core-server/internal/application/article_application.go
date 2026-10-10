package application

import (
	"context"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/dto"
)

type ArticleService struct {
}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (r *ArticleService) CreateArticle(ctx *context.Context, req dto.CreateArticleRequest) error {
	// 调用domain定义的repo层接口

	return nil
}

func (r *ArticleService) EditorArticle(ctx *context.Context, req dto.EditorArticleRequest) error {
}

func (r *ArticleService) DeleteArticle(ctx *context.Context, req dto.DeleteArticleRequest) error {
}

func (r *ArticleService) PublishDraft(ctx *context.Context, req dto.PublishDraftRequest) error {
}

func (r *ArticleService) DeleteDraft(ctx *context.Context, req dto.DeleteDraftRequest) error {
}

func (r *ArticleService) GetArticle(ctx *context.Context, req dto.GetArticleRequest) (dto.GetArticleResponse, error) {
	var aggregate aggregate.ArticleAggregate

	return dto.ToArticleResponse(aggregate)
}

func (r *ArticleService) ListArticles(ctx *context.Context, req dto.ListArticlesRequest) (dto.ListArticlesResponse, error) {
}

func (r *ArticleService) ListMyArticles(ctx *context.Context, req dto.ListMyArticlesRequest) (dto.ListMyArticlesResponse, error) {
}

func (r *ArticleService) ListUserArticles(ctx *context.Context, req dto.ListUserArticlesRequest) (dto.ListUserArticlesResponse, error) {
}

func (r *ArticleService) ListByCategory(ctx *context.Context, req dto.ListByCategoryRequest) (dto.ListByCategoryResponse, error) {
}

func (r *ArticleService) SearchArticles(ctx *context.Context, req dto.SearchArticlesRequest) (dto.SearchArticlesResponse, error) {
}
