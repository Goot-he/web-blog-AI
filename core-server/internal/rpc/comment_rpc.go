package rpc

import (
	"context"

	"core-server/internal/application"
	"core-server/internal/model/dto"
	"core-server/internal/rpc/commentpb"
)

type CommentRpc struct {
	commentpb.UnimplementedCommentServiceServer
	commentService *application.CommentService
}

func NewCommentRpc(commentService *application.CommentService) *CommentRpc {
	return &CommentRpc{
		commentService: commentService,
	}
}

// CreateComment 发表一级评论
func (r *CommentRpc) CreateComment(ctx context.Context, req *commentpb.CreateCommentReq) (*commentpb.CreateCommentResp, error) {
	err := r.commentService.CreateComment(ctx, dto.CreateCommentRequest{
		ArticleID: req.GetArticleId(),
		UserID:    req.GetUserId(),
		Content:   req.GetContent(),
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.CreateCommentResp{Success: true}, nil
}

// CreateReply 回复某条评论
func (r *CommentRpc) CreateReply(ctx context.Context, req *commentpb.CreateReplyReq) (*commentpb.CreateReplyResp, error) {
	err := r.commentService.CreateReply(ctx, dto.CreateReplyRequest{
		RootID:    req.GetRootId(),
		UserID:    req.GetUserId(),
		ReplyToID: req.GetReplyToId(),
		Content:   req.GetContent(),
		ArticleID: req.GetArticleId(),
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.CreateReplyResp{Success: true}, nil
}

// DeleteComment 删除评论
func (r *CommentRpc) DeleteComment(ctx context.Context, req *commentpb.DeleteCommentReq) (*commentpb.DeleteCommentResp, error) {
	err := r.commentService.DeleteComment(ctx, dto.DeleteCommentRequest{
		ID:        req.GetId(),
		UserID:    req.GetUserId(),
		ArticleID: req.GetArticleId(),
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.DeleteCommentResp{Success: true}, nil
}

// GetArticleComments 分页查询文章下的一级评论
func (r *CommentRpc) GetArticleComments(ctx context.Context, req *commentpb.GetArticleCommentsReq) (*commentpb.GetArticleCommentsResp, error) {
	resp, err := r.commentService.GetArticleComments(ctx, dto.GetArticleCommentsRequest{
		ArticleID: req.GetArticleId(),
		Page:      req.GetPage(),
		Size:      req.GetSize(),
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.GetArticleCommentsResp{
		Comments: toPbCommentInfos(resp.Comments),
		Page:     resp.Page,
		Size:     resp.Size,
	}, nil
}

// GetCommentReplies 分页查询根评论下的回复
func (r *CommentRpc) GetCommentReplies(ctx context.Context, req *commentpb.GetCommentRepliesReq) (*commentpb.GetCommentRepliesResp, error) {
	resp, err := r.commentService.GetCommentReplies(ctx, dto.GetCommentRepliesRequest{
		RootID: req.GetRootId(),
		Page:   req.GetPage(),
		Size:   req.GetSize(),
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.GetCommentRepliesResp{
		Replies: toPbCommentInfos(resp.Replies),
		Page:    int32(resp.Page),
		Size:    resp.Size,
	}, nil
}

// toPbCommentInfo dto.CommentInfo 转 proto
func toPbCommentInfo(c *dto.CommentInfo) *commentpb.CommentInfo {
	if c == nil {
		return nil
	}
	return &commentpb.CommentInfo{
		Id:              c.ID,
		ArticleId:       c.ArticleID,
		UserId:          c.UserID,
		RootId:          c.RootID,
		ReplyToId:       c.ReplyToID,
		Content:         c.Content,
		LikeCount:       c.LikeCount,
		ChildCount:      c.ChildCount,
		CreatedAt:       c.CreatedAt,
		UserName:        c.UserName,
		UserAvatar:      c.UserAvatar,
		IsLiked:         c.IsLiked,
		ReplyToUserName: c.ReplyToUserName,
	}
}

// toPbCommentInfos 批量转换
func toPbCommentInfos(comments []*dto.CommentInfo) []*commentpb.CommentInfo {
	if len(comments) == 0 {
		return nil
	}
	out := make([]*commentpb.CommentInfo, 0, len(comments))
	for _, c := range comments {
		out = append(out, toPbCommentInfo(c))
	}
	return out
}
