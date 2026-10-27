package application

import (
	"context"
	"errors"
	"strings"

	"core-server/internal/domain"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
)

type CommentService struct {
	CommentDomain domain.CommentRepoDomain
}

func NewCommentService(commentRepo domain.CommentRepoDomain) *CommentService {
	return &CommentService{
		CommentDomain: commentRepo,
	}
}

// CreateComment 发表一级评论
func (s *CommentService) CreateComment(ctx context.Context, req dto.CreateCommentRequest) error {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return errors.New("comment content is required")
	}
	if req.ArticleID == 0 || req.UserID == 0 {
		return errors.New("articleID and userID are required")
	}
	_, err := s.CommentDomain.CreateComment(ctx, &entity.Comment{
		ArticleID:     req.ArticleID,
		PublishUserID: req.UserID,
		Comment:       content,
	})
	return err
}

// CreateReply 回复某条评论
func (s *CommentService) CreateReply(ctx context.Context, req dto.CreateReplyRequest) error {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return errors.New("reply content is required")
	}
	if req.RootID == 0 || req.ReplyToID == 0 || req.ArticleID == 0 || req.UserID == 0 {
		return errors.New("rootID, replyToID, articleID and userID are required")
	}
	_, err := s.CommentDomain.CreateReply(ctx, &entity.Comment{
		ArticleID:     req.ArticleID,
		PublishUserID: req.UserID,
		RootID:        req.RootID,
		ReplyToID:     req.ReplyToID,
		Comment:       content,
	})
	return err
}

// DeleteComment 删除评论（评论作者或文章作者可删）
func (s *CommentService) DeleteComment(ctx context.Context, req dto.DeleteCommentRequest) error {
	if req.ID == 0 || req.UserID == 0 {
		return errors.New("id and userID are required")
	}
	return s.CommentDomain.SoftDelete(ctx, req.ID, req.UserID)
}

// GetArticleComments 分页查询文章下的一级评论
func (s *CommentService) GetArticleComments(ctx context.Context, req dto.GetArticleCommentsRequest) (*dto.GetArticleCommentsResponse, error) {
	if req.ArticleID == 0 {
		return nil, errors.New("articleID is required")
	}
	offset, limit := pagination(req.Page, req.Size)
	comments, err := s.CommentDomain.ListTopByArticle(ctx, req.ArticleID, offset, limit)
	if err != nil {
		return nil, err
	}
	return &dto.GetArticleCommentsResponse{
		Comments: dto.ToCommentInfos(comments),
		Page:     req.Page,
		Size:     req.Size,
	}, nil
}

// GetCommentReplies 分页查询根评论下的回复
func (s *CommentService) GetCommentReplies(ctx context.Context, req dto.GetCommentRepliesRequest) (*dto.GetCommentRepliesResponse, error) {
	if req.RootID == 0 {
		return nil, errors.New("rootID is required")
	}
	offset, limit := pagination(req.Page, req.Size)
	replies, err := s.CommentDomain.ListRepliesByRoot(ctx, req.RootID, offset, limit)
	if err != nil {
		return nil, err
	}
	return &dto.GetCommentRepliesResponse{
		Replies: dto.ToCommentInfos(replies),
		Page:    req.Page,
		Size:    req.Size,
	}, nil
}
