package dto

import (
	"core-server/internal/model/aggregate"
	"time"
)

type CreateArticleRequest struct {
	Title       string
	Summary     string
	Content     string
	CoverImage  string
	AuthorID    uint64
	CategoryID  uint64
	IsTop       bool
	IsPublished bool
}
type EditorArticleRequest struct {
	ID          uint64
	Title       string
	Summary     string
	Content     string
	CoverImage  string
	AuthorID    uint64
	CategoryID  uint64
	IsTop       bool
	IsPublished bool
}
type DeleteArticleRequest struct {
	ID     uint64
	UserID uint64
}
type PublishDraftRequest struct {
	ID     uint64
	UserID uint64
}
type DeleteDraftRequest struct {
	ID     uint64
	UserID uint64
}
type GetArticleRequest struct {
	ID uint64
}
type ListArticlesRequest struct { // 展示全部文章
	Page uint32
	Size uint32
}
type ListMyArticlesRequest struct {
	UserID uint64
	Page   uint32
	Size   uint32
	IsAll  bool // 为true表示查看自己全部文章 为false表示只看草稿
}
type ListUserArticlesRequest struct {
	UserID uint64
	Page   uint32
	Size   uint32
}
type ListByCategoryRequest struct {
	CategoryID uint64
	Page       uint32
	Size       uint32
}
type SearchArticlesRequest struct {
	Keywords string
	Page     uint32
	Size     uint32
}

type ArticleMsg struct {
	ID          uint64
	Title       string
	Summary     string
	Content     string
	CoverImage  string
	IsTop       bool
	IsPublished bool

	AuthorID     uint64
	AuthorName   string
	AuthorAvatar string

	CategoryID   uint64
	CategoryName string

	ViewCount    uint64
	LikeCount    uint64
	FavorCount   uint64
	CommentCount uint64

	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt time.Time
}

type GetArticleResponse struct {
	*ArticleMsg
}

type ListArticlesResponse struct {
	Articles []*ArticleMsg
}
type ListMyArticlesResponse struct {
	Articles []*ArticleMsg
}
type ListUserArticlesResponse struct {
	Articles []*ArticleMsg
}
type ListByCategoryResponse struct {
	Articles []*ArticleMsg
}
type SearchArticlesResponse struct {
	Articles []*ArticleMsg
}

func ToArticleResponse(a *aggregate.ArticleAggregate) *ArticleMsg {
	article := &ArticleMsg{
		ID:           a.Article.ID,
		Title:        a.Article.Title,
		Summary:      a.Article.Summary,
		Content:      a.Article.Content,
		CoverImage:   a.Article.CoverImage,
		IsTop:        a.Article.IsTop,
		IsPublished:  a.Article.IsPublished,
		AuthorID:     a.Article.AuthorID,
		CategoryID:   a.Article.CategoryID,
		ViewCount:    a.Article.ViewCount,
		LikeCount:    a.Article.LikeCount,
		FavorCount:   a.Article.FavorCount,
		CommentCount: a.Article.CommentCount,
		CreatedAt:    a.Article.CreatedAt,
		UpdatedAt:    a.Article.UpdatedAt,
	}
	if a.Article.PublishedAt.Valid {
		article.PublishedAt = a.Article.PublishedAt.Time
	}

	if a.Author != nil {
	}
	if a.Category != nil {
	}

	return article
}

func ToArticlesResponse(articles []*aggregate.ArticleAggregate) []*ArticleMsg {
	var articleMsgs []*ArticleMsg
	for _, article := range articles {
		articleMsgs = append(articleMsgs, ToArticleResponse(article))
	}
	return articleMsgs
}

func ToGetArticleResponse(article *aggregate.ArticleAggregate) *GetArticleResponse {
	return &GetArticleResponse{
		ArticleMsg: ToArticleResponse(article),
	}
}

func ToListArticlesResponse(articles []*aggregate.ArticleAggregate) *ListArticlesResponse {
	return &ListArticlesResponse{
		Articles: ToArticlesResponse(articles),
	}
}

func ToListMyArticlesResponse(articles []*aggregate.ArticleAggregate) *ListMyArticlesResponse {
	return &ListMyArticlesResponse{
		Articles: ToArticlesResponse(articles),
	}
}

func ToListByCategoryResponse(articles []*aggregate.ArticleAggregate) *ListByCategoryResponse {
	return &ListByCategoryResponse{
		Articles: ToArticlesResponse(articles),
	}
}

func ToListUserArticlesResponse(articles []*aggregate.ArticleAggregate) *ListUserArticlesResponse {
	return &ListUserArticlesResponse{
		Articles: ToArticlesResponse(articles),
	}
}

func ToSearchArticlesResponse(articles []*aggregate.ArticleAggregate) *SearchArticlesResponse {
	return &SearchArticlesResponse{
		Articles: ToArticlesResponse(articles),
	}
}
