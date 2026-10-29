package domain

import (
	"context"
	"core-server/internal/model/entity"
)

type CountRepoDomain interface {
	Upsert(ctx context.Context, count *entity.InteractionCount, delta int64) error
	GetByObject(ctx context.Context, objectType string, objectID uint64) ([]*entity.InteractionCount, error)
	GetByObjects(ctx context.Context, objectType string, objectIDs []uint64) ([]*entity.InteractionCount, error)
}
