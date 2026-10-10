package rpc

import (
	"context"
	"core-server/internal/application"
	"core-server/internal/model/dto"
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

func (r *ArticleRpc) CreateArticle(ctx *context.Context, req *articlepb.CreateArticleRequest) (*articlepb.CreateArticleResponse, error) {
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

func (r *ArticleRpc) EditorArticle(ctx *context.Context, req *articlepb.EditorArticleRequest) (*articlepb.EditorArticleResponse, error) {
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

func (r *ArticleRpc) DeleteArticle(ctx *context.Context, req *articlepb.DeleteArticleRequest) (*articlepb.DeleteArticleResponse, error) {
	err := r.articleService.DeleteArticle(ctx, dto.DeleteArticleRequest{
		ID:     req.Id,
		UserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &articlepb.DeleteArticleResponse{
		Success: true,
	}, nil
}

func (r *ArticleRpc) PublishDraft(ctx *context.Context, req *articlepb.PublishDraftRequest) (*articlepb.PublishDraftResponse, error) {
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

func (r *ArticleRpc) DeleteDraft(ctx *context.Context, req *articlepb.DeleteDraftRequest) (*articlepb.DeleteDraftResponse, error) {
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

func (r *ArticleRpc) GetArticle(ctx *context.Context, req *articlepb.GetArticleRequest) (*articlepb.GetArticleResponse, error) {
	resp, err := r.articleService.GetArticle(ctx, dto.GetArticleRequest{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}
	article := DtoToPB(dto.ArticleMsg{
		Title:       resp.Title,
		Summary:     resp.Summary,
		Content:     resp.Content,
		CoverImage:  resp.CoverImage,
		AuthorID:    resp.AuthorID,
		CategoryID:  resp.CategoryID,
		IsTop:       resp.IsTop,
		IsPublished: resp.IsPublished,
	})
	return article, nil
}

func (r *ArticleRpc) ListArticles(ctx *context.Context, req *articlepb.ListArticlesRequest) (*articlepb.ListArticlesResponse, error) {
}

func (r *ArticleRpc) ListMyArticles(ctx *context.Context, req *articlepb.ListMyArticlesRequest) (*articlepb.ListMyArticlesResponse, error) {
}

func (r *ArticleRpc) ListUserArticles(ctx *context.Context, req *articlepb.ListUserArticlesRequest) (*articlepb.ListUserArticlesResponse, error) {
}

func (r *ArticleRpc) ListByCategory(ctx *context.Context, req *articlepb.ListByCategoryRequest) (*articlepb.ListByCategoryResponse, error) {
}

func (r *ArticleRpc) SearchArticles(ctx *context.Context, req *articlepb.SearchArticlesRequest) (*articlepb.SearchArticlesResponse, error) {
}

func DtoToPB(msg dto.ArticleMsg) *articlepb.GetArticleResponse {
	return &articlepb.GetArticleResponse{
		Title:       msg.Title,
		Summary:     msg.Summary,
		Content:     msg.Content,
		CoverImage:  msg.CoverImage,
		AuthorID:    msg.AuthorID,
		CategoryID:  msg.CategoryID,
		IsTop:       msg.IsTop,
		IsPublished: msg.IsPublished,
	}
}

func DtoToPBs() {

}
