package domain

import (
	"context"

	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
)

type CommentRepo interface {
	// Create 新建一级评论（同时更新文章评论数）
	Create(ctx context.Context, comment *entity.Comment) error
	// CreateReply 新建回复（同时更新文章评论数和根评论的子评论数）
	CreateReply(ctx context.Context, comment *entity.Comment) error
	// Delete 软删除评论（评论作者或文章作者可删）
	Delete(ctx context.Context, commentID, userID, articleID uint64) error
	// ListByArticle 分页查询文章下的一级评论
	ListByArticle(ctx context.Context, articleID uint64, page, size uint32) ([]*aggregate.CommentAggregate, error)
	// ListByRoot 分页查询根评论下的回复
	ListByRoot(ctx context.Context, rootID uint64, page, size uint32) ([]*aggregate.CommentAggregate, error)
}

type CommentDomainRepo interface {
	CommentRepo
	ITransaction
}
