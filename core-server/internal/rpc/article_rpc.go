package rpc

import (
	"context"
	"core-server/internal/rpc/articlepb"
)

type ArticleRpc struct {
	articlepb.UnimplementedArticleServiceServer
}

func NewArticleRpc() *ArticleRpc {
	return &ArticleRpc{}
}

func (r *ArticleRpc) CreateArticle(ctx *context.Context, req articlepb.CreateArticleRequest) (*articlepb.CreateArticleResponse, error) {
}

func (r *ArticleRpc) EditorArticle(ctx *context.Context, req articlepb.EditorArticleRequest) (*articlepb.EditorArticleResponse, error) {
}

func (r *ArticleRpc) DeleteArticle(ctx *context.Context, req articlepb.DeleteArticleRequest) (*articlepb.DeleteArticleResponse, error) {
}

func (r *ArticleRpc) PublishDraft(ctx *context.Context, req articlepb.PublishDraftRequest) (*articlepb.PublishDraftResponse, error) {
}

func (r *ArticleRpc) DeleteDraft(ctx *context.Context, req articlepb.DeleteDraftRequest) (*articlepb.DeleteDraftResponse, error) {
}

func (r *ArticleRpc) GetArticle(ctx *context.Context, req articlepb.GetArticleRequest) (*articlepb.GetArticleResponse, error) {
}

func (r *ArticleRpc) ListArticles(ctx *context.Context, req articlepb.ListArticlesRequest) (*articlepb.ListArticlesResponse, error) {
}

func (r *ArticleRpc) ListMyArticles(ctx *context.Context, req articlepb.ListMyArticlesRequest) (*articlepb.ListMyArticlesResponse, error) {
}

func (r *ArticleRpc) ListUserArticles(ctx *context.Context, req articlepb.ListUserArticlesRequest) (*articlepb.ListUserArticlesResponse, error) {
}

func (r *ArticleRpc) ListByCategory(ctx *context.Context, req articlepb.ListByCategoryRequest) (*articlepb.ListByCategoryResponse, error) {
}

func (r *ArticleRpc) SearchArticles(ctx *context.Context, req articlepb.SearchArticlesRequest) (*articlepb.SearchArticlesResponse, error) {
}
