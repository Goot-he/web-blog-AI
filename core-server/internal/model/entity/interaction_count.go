package entity

import (
	"core-server/internal/model/enum"
	"database/sql"
	"time"
)

type InteractionCount struct {
	ID        uint64       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`

	ObjectType      enum.ObjectType      `db:"object_type"`
	ObjectID        uint64               `db:"object_id"`
	InteractionType enum.InteractionType `db:"interaction_type"`
	Count           uint64               `db:"count"`
}
