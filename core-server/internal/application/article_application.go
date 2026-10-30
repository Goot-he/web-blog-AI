package application

import (
	"context"

	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
)

type ArticleService struct {
	ArtDomain domain.ArticleDomainRepo
	log       *clog.Log
}

func NewArticleService(articleRepo domain.ArticleDomainRepo, log *clog.Log) *ArticleService {
	return &ArticleService{
		ArtDomain: articleRepo,
		log:       log,
	}
}

// CreateArticle 新建文章
func (r *ArticleService) CreateArticle(ctx context.Context, req dto.CreateArticleRequest) error {
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
		r.log.Error("create article failed")
		return err
	}
	r.log.Info("article created")
	return nil
}

// EditorArticle 编辑文章
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
		r.log.Error("edit article failed")
		return err
	}
	r.log.Info("article updated")
	return nil
}

// DeleteArticle 删除文章
func (r *ArticleService) DeleteArticle(ctx context.Context, req dto.DeleteArticleRequest) error {
	err := r.ArtDomain.Delete(ctx, req.ID, req.UserID, req.Role)
	if err != nil {
		r.log.Error("delete article failed")
		return err
	}
	r.log.Info("article deleted")
	return nil
}

// PublishDraft 发布草稿
func (r *ArticleService) PublishDraft(ctx context.Context, req dto.PublishDraftRequest) error {
	err := r.ArtDomain.PublishDraft(ctx, req.ID, req.UserID)
	if err != nil {
		r.log.Error("publish draft failed")
		return err
	}
	r.log.Info("draft published")
	return nil
}

// DeleteDraft 删除草稿
func (r *ArticleService) DeleteDraft(ctx context.Context, req dto.DeleteDraftRequest) error {
	err := r.ArtDomain.DeleteDraft(ctx, req.ID, req.UserID)
	if err != nil {
		r.log.Error("delete draft failed")
		return err
	}
	r.log.Info("draft deleted")
	return nil
}

// GetArticle 查询文章详情
func (r *ArticleService) GetArticle(ctx context.Context, req dto.GetArticleRequest) (*dto.GetArticleResponse, error) {
	agg, err := r.ArtDomain.GetArticle(ctx, req.ID)
	if err != nil {
		r.log.Warn("get article failed")
		return nil, err
	}
	return dto.ToGetArticleResponse(agg), nil
}

// ListArticles 分页查询全部文章
func (r *ArticleService) ListArticles(ctx context.Context, req dto.ListArticlesRequest) (*dto.ListArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		r.log.Warn("list articles failed")
		return nil, err
	}
	r.log.Debug("list articles")
	return dto.ToListArticlesResponse(aggs), nil
}

// ListMyArticles 分页查询我的文章
func (r *ArticleService) ListMyArticles(ctx context.Context, req dto.ListMyArticlesRequest) (*dto.ListMyArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListMyArticles(ctx, req.UserID, req.IsPublished, req.Page, req.Size)
	if err != nil {
		r.log.Warn("list my articles failed")
		return nil, err
	}
	r.log.Debug("list my articles")
	return dto.ToListMyArticlesResponse(aggs), nil
}

// ListUserArticles 分页查询指定用户的文章
func (r *ArticleService) ListUserArticles(ctx context.Context, req dto.ListUserArticlesRequest) (*dto.ListUserArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListUserArticles(ctx, req.UserID, req.Page, req.Size)
	if err != nil {
		r.log.Warn("list user articles failed")
		return nil, err
	}
	r.log.Debug("list user articles")
	return dto.ToListUserArticlesResponse(aggs), nil
}

// ListByCategory 分页查询分类下的文章
func (r *ArticleService) ListByCategory(ctx context.Context, req dto.ListByCategoryRequest) (*dto.ListByCategoryResponse, error) {
	aggs, err := r.ArtDomain.ListByCategory(ctx, req.CategoryID, req.Page, req.Size)
	if err != nil {
		r.log.Warn("list articles by category failed")
		return nil, err
	}
	r.log.Debug("list articles by category")
	return dto.ToListByCategoryResponse(aggs), nil
}

// SearchArticles 搜索文章
func (r *ArticleService) SearchArticles(ctx context.Context, req dto.SearchArticlesRequest) (*dto.SearchArticlesResponse, error) {
	aggs, err := r.ArtDomain.SearchArticles(ctx, req.Keywords, req.Page, req.Size)
	if err != nil {
		r.log.Warn("search articles failed")
		return nil, err
	}
	r.log.Debug("search articles")
	return dto.ToSearchArticlesResponse(aggs), nil
}
