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
	Count           int                  `db:"count"`
}

func (InteractionCount) TableName() string {
	return "interaction_count"
}

type InteractionStatus struct {
	LikeCount    int `db:"like_count"`
	CommentCount int `db:"comment_count"`
	ViewCount    int `db:"view_count"`
	FavorCount   int `db:"favor_count"`
	FollowCount  int `db:"follow_count"`
}
