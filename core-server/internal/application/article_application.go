package application

import (
	"context"

	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"

	"go.uber.org/zap"
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
		r.log.Error("create article failed",
			zap.Error(err),
			zap.Uint64("author_id", req.AuthorID),
			zap.Uint64("category_id", req.CategoryID),
			zap.String("title", req.Title),
		)
		return err
	}
	r.log.Info("article created",
		zap.Uint64("author_id", req.AuthorID),
		zap.Uint64("category_id", req.CategoryID),
		zap.Bool("published", req.IsPublished),
	)
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
		r.log.Error("edit article failed",
			zap.Error(err),
			zap.Uint64("article_id", req.ID),
			zap.Uint64("author_id", req.AuthorID),
			zap.Uint64("category_id", req.CategoryID),
		)
		return err
	}
	r.log.Info("article updated",
		zap.Uint64("article_id", req.ID),
		zap.Uint64("author_id", req.AuthorID),
		zap.Bool("published", req.IsPublished),
	)
	return nil
}

// DeleteArticle 删除文章
func (r *ArticleService) DeleteArticle(ctx context.Context, req dto.DeleteArticleRequest) error {
	err := r.ArtDomain.Delete(ctx, req.ID, req.UserID, req.Role)
	if err != nil {
		r.log.Error("delete article failed",
			zap.Error(err),
			zap.Uint64("article_id", req.ID),
			zap.Uint64("user_id", req.UserID),
			zap.Int("role", int(req.Role)),
		)
		return err
	}
	r.log.Info("article deleted",
		zap.Uint64("article_id", req.ID),
		zap.Uint64("user_id", req.UserID),
		zap.Int("role", int(req.Role)),
	)
	return nil
}

// PublishDraft 发布草稿
func (r *ArticleService) PublishDraft(ctx context.Context, req dto.PublishDraftRequest) error {
	err := r.ArtDomain.PublishDraft(ctx, req.ID, req.UserID)
	if err != nil {
		r.log.Error("publish draft failed",
			zap.Error(err),
			zap.Uint64("article_id", req.ID),
			zap.Uint64("user_id", req.UserID),
		)
		return err
	}
	r.log.Info("draft published",
		zap.Uint64("article_id", req.ID),
		zap.Uint64("user_id", req.UserID),
	)
	return nil
}

// DeleteDraft 删除草稿
func (r *ArticleService) DeleteDraft(ctx context.Context, req dto.DeleteDraftRequest) error {
	err := r.ArtDomain.DeleteDraft(ctx, req.ID, req.UserID)
	if err != nil {
		r.log.Error("delete draft failed",
			zap.Error(err),
			zap.Uint64("article_id", req.ID),
			zap.Uint64("user_id", req.UserID),
		)
		return err
	}
	r.log.Info("draft deleted",
		zap.Uint64("article_id", req.ID),
		zap.Uint64("user_id", req.UserID),
	)
	return nil
}

// GetArticle 查询文章详情
func (r *ArticleService) GetArticle(ctx context.Context, req dto.GetArticleRequest) (*dto.GetArticleResponse, error) {
	agg, err := r.ArtDomain.GetArticle(ctx, req.ID)
	if err != nil {
		r.log.Warn("get article failed",
			zap.Error(err),
			zap.Uint64("article_id", req.ID),
		)
		return nil, err
	}
	return dto.ToGetArticleResponse(agg), nil
}

// ListArticles 分页查询全部文章
func (r *ArticleService) ListArticles(ctx context.Context, req dto.ListArticlesRequest) (*dto.ListArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		r.log.Warn("list articles failed",
			zap.Error(err),
			zap.Uint32("page", req.Page),
			zap.Uint32("size", req.Size),
		)
		return nil, err
	}
	r.log.Debug("list articles", zap.Int("count", len(aggs)))
	return dto.ToListArticlesResponse(aggs), nil
}

// ListMyArticles 分页查询我的文章
func (r *ArticleService) ListMyArticles(ctx context.Context, req dto.ListMyArticlesRequest) (*dto.ListMyArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListMyArticles(ctx, req.UserID, req.IsAll, req.Page, req.Size)
	if err != nil {
		r.log.Warn("list my articles failed",
			zap.Error(err),
			zap.Uint64("user_id", req.UserID),
			zap.Bool("is_all", req.IsAll),
		)
		return nil, err
	}
	r.log.Debug("list my articles", zap.Uint64("user_id", req.UserID), zap.Int("count", len(aggs)))
	return dto.ToListMyArticlesResponse(aggs), nil
}

// ListUserArticles 分页查询指定用户的文章
func (r *ArticleService) ListUserArticles(ctx context.Context, req dto.ListUserArticlesRequest) (*dto.ListUserArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		r.log.Warn("list user articles failed",
			zap.Error(err),
			zap.Uint64("user_id", req.UserID),
		)
		return nil, err
	}
	r.log.Debug("list user articles", zap.Uint64("user_id", req.UserID), zap.Int("count", len(aggs)))
	return dto.ToListUserArticlesResponse(aggs), nil
}

// ListByCategory 分页查询分类下的文章
func (r *ArticleService) ListByCategory(ctx context.Context, req dto.ListByCategoryRequest) (*dto.ListByCategoryResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		r.log.Warn("list articles by category failed",
			zap.Error(err),
			zap.Uint64("category_id", req.CategoryID),
		)
		return nil, err
	}
	r.log.Debug("list articles by category", zap.Uint64("category_id", req.CategoryID), zap.Int("count", len(aggs)))
	return dto.ToListByCategoryResponse(aggs), nil
}

// SearchArticles 搜索文章
func (r *ArticleService) SearchArticles(ctx context.Context, req dto.SearchArticlesRequest) (*dto.SearchArticlesResponse, error) {
	aggs, err := r.ArtDomain.ListArticles(ctx, req.Page, req.Size)
	if err != nil {
		r.log.Warn("search articles failed",
			zap.Error(err),
			zap.String("keyword", req.Keywords),
			zap.Uint32("page", req.Page),
			zap.Uint32("size", req.Size),
		)
		return nil, err
	}
	r.log.Debug("search articles",
		zap.String("keyword", req.Keywords),
		zap.Int("count", len(aggs)),
	)
	return dto.ToSearchArticlesResponse(aggs), nil
}
