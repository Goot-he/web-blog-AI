package application

import (
	"context"
	"core-server/internal/rpc/articlepb"
)

type ArticleService struct {
}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (s *ArticleService) CreateArticle(ctx *context.Context, req articlepb.CreateArticleRequest) (*articlepb.CreateArticleResponse, error) {
}
