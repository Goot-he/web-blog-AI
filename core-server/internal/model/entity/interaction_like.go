package entity

import (
	"core-server/internal/model/enum"
	"database/sql"
	"time"
)

type InteractionLike struct {
	ID        uint64       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`

	ObjectType enum.ObjectType `db:"object_type"`
	ObjectID   uint64          `db:"object_id"`
	UserID     uint64          `db:"user_id"`
	LikeStatus enum.Status     `db:"like_status"`
}

type InteractionLikeType string

// 点赞状态机
const (
	LikeStatusTypeUnknown InteractionLikeType = "unknown"  // 这个状态用来规避就是删除实体之后，所有的点赞关系置为此，表示失效
	LikeStatusTypeThumbUp InteractionLikeType = "thumb_up" // 已点赞
	LikeStatusTypeNothing InteractionLikeType = "nothing"  // 取消点赞
)

func (t InteractionLikeType) String() string {
	return string(t)
}

func ParseInteractionLikeType(s string) InteractionLikeType {
	switch s {
	case "thumb_up":
		return LikeStatusTypeThumbUp
	case "nothing":
		return LikeStatusTypeNothing
	default:
		return LikeStatusTypeUnknown
	}
}
