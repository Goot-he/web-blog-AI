package dto

import (
	"core-server/internal/model/aggregate"
)

const dateTimeLayout = "2006-01-02 15:04:05"

type CommentInfo struct {
	ID              uint64
	ArticleID       uint64
	UserID          uint64
	RootID          uint64
	ReplyToID       uint64
	Content         string
	LikeCount       uint32
	ChildCount      uint32
	CreatedAt       string
	UserName        string
	UserAvatar      string
	IsLiked         bool
	ReplyToUserName string
}

type CreateCommentRequest struct {
	ArticleID uint64
	UserID    uint64
	Content   string
}

type CreateReplyRequest struct {
	RootID    uint64
	UserID    uint64
	ReplyToID uint64
	Content   string
	ArticleID uint64
}

type DeleteCommentRequest struct {
	ID        uint64
	UserID    uint64
	ArticleID uint64
}

type GetArticleCommentsRequest struct {
	ArticleID uint64
	Page      uint32
	Size      uint32
}
type GetArticleCommentsResponse struct {
	Comments []*CommentInfo
	Page     uint32
	Size     uint32
}

type GetCommentRepliesRequest struct {
	RootID uint64
	Page   uint32
	Size   uint32
}
type GetCommentRepliesResponse struct {
	Replies []*CommentInfo
	Page    uint32
	Size    uint32
}

// ToCommentInfo 评论聚合转评论信息
func ToCommentInfo(agg *aggregate.CommentAggregate) *CommentInfo {
	if agg == nil || agg.Comment == nil {
		return nil
	}
	info := &CommentInfo{
		ID:         agg.Comment.ID,
		ArticleID:  agg.Comment.ArticleID,
		UserID:     agg.Comment.PublishUserID,
		RootID:     agg.Comment.RootID,
		ReplyToID:  agg.Comment.ReplyToID,
		Content:    agg.Comment.Comment,
		LikeCount:  uint32(agg.Comment.LikeCount),
		ChildCount: uint32(agg.Comment.ChildCount),
		CreatedAt:  agg.Comment.CreatedAt.Format(dateTimeLayout),
	}
	if agg.User != nil {
		info.UserName = agg.User.Username
		info.UserAvatar = agg.User.Avatar
	}
	return info
}

// ToCommentInfos 批量转换评论信息
func ToCommentInfos(aggs []*aggregate.CommentAggregate) []*CommentInfo {
	if len(aggs) == 0 {
		return nil
	}
	out := make([]*CommentInfo, 0, len(aggs))
	for _, agg := range aggs {
		out = append(out, ToCommentInfo(agg))
	}
	return out
}
