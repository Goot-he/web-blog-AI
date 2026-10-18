package entity

import (
	"database/sql"
	"time"
)

type Status struct {
	ID        uint64       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`

	UserID             uint64 `db:"user_id"`
	ArticleCount       uint64 `db:"article_count"`
	LikeCount          uint64 `db:"like_count"`
	ReceiveLikeCount   uint64 `db:"receive_like_count"`
	FavorCount         uint64 `db:"favor_count"`
	ReceiveFavorCount  uint64 `db:"receive_favor_count"`
	FollowCount        uint64 `db:"follow_count"`
	ReceiveFollowCount uint64 `db:"receive_follow_count"`
	ViewCount          uint64 `db:"view_count"`
}
