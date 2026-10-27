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
	PageQueryLikeObjects(ctx context.Context, userID uint64, objectType enum.ObjectType, page, size int) (uint64, error)
}

type LikeDomainRepo interface {
	LikeRepo
	ITransaction
}
