package domain

import (
	"context"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type LikeRepo interface {
	Upsert(ctx context.Context, like *entity.InteractionLike) (int, error)
	UpdateWithCondition(ctx context.Context, condition string, like entity.InteractionLike) (int, error)
	QueryWithCondition(ctx context.Context, userID, objectID uint64, objectType enum.ObjectType, status string) (*entity.InteractionLike, error)
	CountUserLiked(ctx context.Context, userID, objectID uint64, objectType enum.ObjectType) (int, error)
	PageQueryLikeObjects(ctx context.Context, userID uint64, objectType enum.ObjectType, offset, limit int) ([]*entity.InteractionLike, error)
}

type LikeCacheRepo interface {
	ThumbUp(ctx context.Context, userID uint64, objectType string, objectID uint64, score int64) error
	CancelThumbUp(ctx context.Context, userID uint64, objectType string, objectID uint64) (int, int64, error)
	ExistZSetMember(ctx context.Context, userID uint64, objectType string, objectID uint64) (bool, error)
	CompensationCountDecr(ctx context.Context, objectID uint64, objectType string)
	CompensationCountIncr(ctx context.Context, objectID uint64, objectType string)
	GetObjectLikeCount(ctx context.Context, objectID uint64, objectType string) (uint64, error)
	SetLikeList(ctx context.Context, userID uint64, objectType string, objectIDs []uint64, scores []float64)
	PageQueryObjects(ctx context.Context, userID uint64, objectType string, page, size int) ([]uint64, error)
}

type LikeDomainRepo interface {
	LikeRepo
	LikeCacheRepo
	ITransaction
}
