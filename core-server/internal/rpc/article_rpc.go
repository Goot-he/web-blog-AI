package rpc

import (
	"context"
	"core-server/internal/application"
	"core-server/internal/model/dto"
	"core-server/internal/model/enum"
	"core-server/internal/rpc/articlepb"
)

type ArticleRpc struct {
	articlepb.UnimplementedArticleServiceServer
	articleService application.ArticleService
}

func NewArticleRpc(newArticleService application.ArticleService) *ArticleRpc {
	return &ArticleRpc{
		articleService: newArticleService,
	}
}

func (r *ArticleRpc) CreateArticle(ctx context.Context, req *articlepb.CreateArticleRequest) (*articlepb.CreateArticleResponse, error) {
	//调用application层函数
	err := r.articleService.CreateArticle(ctx, dto.CreateArticleRequest{
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
		return nil, err
	}
	return &articlepb.CreateArticleResponse{
		Success: true,
	}, nil
}

func (r *ArticleRpc) EditorArticle(ctx context.Context, req *articlepb.EditorArticleRequest) (*articlepb.EditorArticleResponse, error) {
	err := r.articleService.EditorArticle(ctx, dto.EditorArticleRequest{
		ID:          req.Id,
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
		return nil, err
	}
	return &articlepb.EditorArticleResponse{
		Success: true,
	}, nil
}

func (r *ArticleRpc) DeleteArticle(ctx context.Context, req *articlepb.DeleteArticleRequest) (*articlepb.DeleteArticleResponse, error) {
	err := r.articleService.DeleteArticle(ctx, dto.DeleteArticleRequest{
		ID:     req.Id,
		UserID: req.UserID,
		Role:   enum.ParseRoleString(req.GetRole()),
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.DeleteArticleResponse{
		Success: true,
	}, nil
}

func (r *ArticleRpc) PublishDraft(ctx context.Context, req *articlepb.PublishDraftRequest) (*articlepb.PublishDraftResponse, error) {
	err := r.articleService.PublishDraft(ctx, dto.PublishDraftRequest{
		ID:     req.Id,
		UserID: req.AuthorID,
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.PublishDraftResponse{
		Success: true,
	}, nil
}

func (r *ArticleRpc) DeleteDraft(ctx context.Context, req *articlepb.DeleteDraftRequest) (*articlepb.DeleteDraftResponse, error) {
	err := r.articleService.DeleteDraft(ctx, dto.DeleteDraftRequest{
		ID:     req.Id,
		UserID: req.AuthorID,
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.DeleteDraftResponse{
		Success: true,
	}, nil
}

func (r *ArticleRpc) GetArticle(ctx context.Context, req *articlepb.GetArticleRequest) (*articlepb.GetArticleResponse, error) {
	resp, err := r.articleService.GetArticle(ctx, dto.GetArticleRequest{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &articlepb.GetArticleResponse{
		Id:           resp.ID,
		Title:        resp.Title,
		Summary:      resp.Summary,
		Content:      resp.Content,
		CoverImage:   resp.CoverImage,
		AuthorID:     resp.AuthorID,
		CategoryID:   resp.CategoryID,
		IsTop:        resp.IsTop,
		IsPublished:  resp.IsPublished,
		ViewCount:    resp.ViewCount,
		LikeCount:    resp.LikeCount,
		CommentCount: resp.CommentCount,
		CreatedAt:    uint64(resp.CreatedAt.Unix()),
		UpdatedAt:    uint64(resp.UpdatedAt.Unix()),
		PublishedAt:  uint64(resp.PublishedAt.Unix()),
		AuthorName:   resp.AuthorName,
		AuthorAvatar: resp.AuthorAvatar,
		FavorCount:   resp.FavorCount,
	}, nil
}

func (r *ArticleRpc) ListArticles(ctx context.Context, req *articlepb.ListArticlesRequest) (*articlepb.ListArticlesResponse, error) {
	resp, err := r.articleService.ListArticles(ctx, dto.ListArticlesRequest{
		Page: req.Page,
		Size: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.ListArticlesResponse{
		Articles: toArticlePBList(resp.Articles),
	}, nil
}

func (r *ArticleRpc) ListMyArticles(ctx context.Context, req *articlepb.ListMyArticlesRequest) (*articlepb.ListMyArticlesResponse, error) {
	resp, err := r.articleService.ListMyArticles(ctx, dto.ListMyArticlesRequest{
		Page: req.Page,
		Size: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.ListMyArticlesResponse{
		Articles: toArticlePBList(resp.Articles),
	}, nil
}

func (r *ArticleRpc) ListUserArticles(ctx context.Context, req *articlepb.ListUserArticlesRequest) (*articlepb.ListUserArticlesResponse, error) {
	resp, err := r.articleService.ListUserArticles(ctx, dto.ListUserArticlesRequest{
		Page:   req.Page,
		Size:   req.PageSize,
		UserID: req.AuthorID,
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.ListUserArticlesResponse{Articles: toArticlePBList(resp.Articles)}, nil
}

func (r *ArticleRpc) ListByCategory(ctx context.Context, req *articlepb.ListByCategoryRequest) (*articlepb.ListByCategoryResponse, error) {
	resp, err := r.articleService.ListByCategory(ctx, dto.ListByCategoryRequest{
		CategoryID: req.CategoryID,
		Page:       req.Page,
		Size:       req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.ListByCategoryResponse{Articles: toArticlePBList(resp.Articles)}, nil
}

func (r *ArticleRpc) SearchArticles(ctx context.Context, req *articlepb.SearchArticlesRequest) (*articlepb.SearchArticlesResponse, error) {
	resp, err := r.articleService.SearchArticles(ctx, dto.SearchArticlesRequest{
		Keywords: req.Q,
		Page:     req.Page,
		Size:     req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.SearchArticlesResponse{Articles: toArticlePBList(resp.Articles)}, nil
}

func toArticlePBList(aggregates []*dto.ArticleMsg) []*articlepb.Article {
	if len(aggregates) == 0 {
		return nil
	}
	items := make([]*articlepb.Article, 0, len(aggregates))
	for _, item := range aggregates {
		items = append(items, toArticlePB(item))
	}
	return items
}

func toArticlePB(a *dto.ArticleMsg) *articlepb.Article {
	if a == nil {
		return nil
	}
	pb := &articlepb.Article{
		Id:           a.ID,
		AuthorID:     a.AuthorID,
		Title:        a.Title,
		Content:      a.Content,
		Summary:      a.Summary,
		CategoryID:   a.CategoryID,
		IsTop:        a.IsTop,
		CoverImage:   a.CoverImage,
		IsPublished:  a.IsPublished,
		ViewCount:    a.ViewCount,
		LikeCount:    a.LikeCount,
		FavorCount:   a.FavorCount,
		CommentCount: a.CommentCount,
		CreatedAt:    uint64(a.CreatedAt.Unix()),
		UpdatedAt:    uint64(a.UpdatedAt.Unix()),
		PublishedAt:  uint64(a.PublishedAt.Unix()),
	}
	pb.AuthorName = a.AuthorName
	pb.AuthorAvatar = a.AuthorAvatar
	return pb
}
