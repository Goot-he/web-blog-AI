package entity

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID            uint64       `db:"id"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at"`
	ArticleID     uint64       `db:"article_id"`
	RootID        uint64       `db:"root_id"`
	ReplyToID     uint64       `db:"reply_to_id"`
	PublishUserID uint64       `db:"publish_user_id"`
	ReplyToUserID uint64       `db:"reply_to_user_id"`
	Comment       string       `db:"comment"`
	LikeCount     uint64       `db:"like_count"`
	ChildCount    uint64       `db:"child_count"`
}
