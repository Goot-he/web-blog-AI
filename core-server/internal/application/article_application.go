package application

import (
	"context"
	"core-server/internal/domain"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/dto"
)

type ArticleService struct {
	articleDomain *domain.ArticleDomainRepo
}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (r *ArticleService) CreateArticle(ctx context.Context, req dto.CreateArticleRequest) error {
	// 调用domain定义的repo层接口

	return nil
}

func (r *ArticleService) EditorArticle(ctx context.Context, req dto.EditorArticleRequest) error {
	return nil
}

func (r *ArticleService) DeleteArticle(ctx context.Context, req dto.DeleteArticleRequest) error {
	return nil
}

func (r *ArticleService) PublishDraft(ctx context.Context, req dto.PublishDraftRequest) error {
	return nil
}

func (r *ArticleService) DeleteDraft(ctx context.Context, req dto.DeleteDraftRequest) error {
	return nil
}

func (r *ArticleService) GetArticle(ctx context.Context, req dto.GetArticleRequest) (*dto.GetArticleResponse, error) {
	var agg *aggregate.ArticleAggregate

	return dto.ToGetArticleResponse(agg), nil
}

func (r *ArticleService) ListArticles(ctx context.Context, req dto.ListArticlesRequest) (*dto.ListArticlesResponse, error) {
	var agg []*aggregate.ArticleAggregate

	return dto.ToListArticlesResponse(agg), nil
}

func (r *ArticleService) ListMyArticles(ctx context.Context, req dto.ListMyArticlesRequest) (*dto.ListMyArticlesResponse, error) {
	var agg []*aggregate.ArticleAggregate

	return dto.ToListMyArticlesResponse(agg), nil
}

func (r *ArticleService) ListUserArticles(ctx context.Context, req dto.ListUserArticlesRequest) (*dto.ListUserArticlesResponse, error) {
	var agg []*aggregate.ArticleAggregate

	return dto.ToListUserArticlesResponse(agg), nil
}

func (r *ArticleService) ListByCategory(ctx context.Context, req dto.ListByCategoryRequest) (*dto.ListByCategoryResponse, error) {
	var agg []*aggregate.ArticleAggregate

	return dto.ToListByCategoryResponse(agg), nil
}

func (r *ArticleService) SearchArticles(ctx context.Context, req dto.SearchArticlesRequest) (*dto.SearchArticlesResponse, error) {
	var agg []*aggregate.ArticleAggregate

	return dto.ToSearchArticlesResponse(agg), nil
}
