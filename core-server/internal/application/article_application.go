package application

import (
	"context"
	"core-server/internal/domain"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
)

type ArticleService struct {
	ArtDomain domain.ArticleDomainRepo
}

func NewArticleService(articleRepo domain.ArticleDomainRepo) *ArticleService {
	return &ArticleService{
		ArtDomain: articleRepo,
	}
}

func (r *ArticleService) CreateArticle(ctx context.Context, req dto.CreateArticleRequest) error {
	// 调用domain定义的repo层接口
	err := r.ArtDomain.Create(ctx, &entity.Article{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		AuthorID:    req.AuthorID,
		CategoryID:  req.CategoryID,
		IsTop:       req.IsTop,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleService) EditorArticle(ctx context.Context, req dto.EditorArticleRequest) error {
	err := r.ArtDomain.Editor(ctx, &entity.Article{
		ID:          req.ID,
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		AuthorID:    req.AuthorID,
		CategoryID:  req.CategoryID,
		IsTop:       req.IsTop,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ArticleService) DeleteArticle(ctx context.Context, req dto.DeleteArticleRequest) error {
	err := r.ArtDomain.Delete(ctx, req.ID, req.UserID, req.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleService) PublishDraft(ctx context.Context, req dto.PublishDraftRequest) error {
	err := r.ArtDomain.PublishDraft(ctx, req.ID, req.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleService) DeleteDraft(ctx context.Context, req dto.DeleteDraftRequest) error {
	err := r.ArtDomain.DeleteDraft(ctx, req.ID, req.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleService) GetArticle(ctx context.Context, req dto.GetArticleRequest) (*dto.GetArticleResponse, error) {
	agg, err := r.ArtDomain.GetArticle(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return dto.ToGetArticleResponse(agg), nil
}

func (r *ArticleService) ListArticles(ctx context.Context, req dto.ListArticlesRequest) (*dto.ListArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return dto.ToListArticlesResponse(aggs), nil
}

func (r *ArticleService) ListMyArticles(ctx context.Context, req dto.ListMyArticlesRequest) (*dto.ListMyArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListMyArticles(ctx, req.UserID, req.IsAll, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return dto.ToListMyArticlesResponse(aggs), nil
}

func (r *ArticleService) ListUserArticles(ctx context.Context, req dto.ListUserArticlesRequest) (*dto.ListUserArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return dto.ToListUserArticlesResponse(aggs), nil
}

func (r *ArticleService) ListByCategory(ctx context.Context, req dto.ListByCategoryRequest) (*dto.ListByCategoryResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return dto.ToListByCategoryResponse(aggs), nil
}

func (r *ArticleService) SearchArticles(ctx context.Context, req dto.SearchArticlesRequest) (*dto.SearchArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return dto.ToSearchArticlesResponse(aggs), nil
}
