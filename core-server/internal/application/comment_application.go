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
	CommentDomain domain.CommentDomainRepo
}

func NewCommentService(commentRepo domain.CommentDomainRepo) *CommentService {
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
	return s.CommentDomain.Create(ctx, &entity.Comment{
		ArticleID:     req.ArticleID,
		PublishUserID: req.UserID,
		Comment:       content,
	})
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
	return s.CommentDomain.CreateReply(ctx, &entity.Comment{
		ArticleID:     req.ArticleID,
		PublishUserID: req.UserID,
		RootID:        req.RootID,
		ReplyToID:     req.ReplyToID,
		Comment:       content,
	})
}

// DeleteComment 删除评论（评论作者或文章作者可删）
func (s *CommentService) DeleteComment(ctx context.Context, req dto.DeleteCommentRequest) error {
	if req.ID == 0 || req.UserID == 0 || req.ArticleID == 0 {
		return errors.New("id, userID and articleID are required")
	}
	return s.CommentDomain.Delete(ctx, req.ID, req.UserID, req.ArticleID)
}

// GetArticleComments 分页查询文章下的一级评论
func (s *CommentService) GetArticleComments(ctx context.Context, req dto.GetArticleCommentsRequest) (*dto.GetArticleCommentsResponse, error) {
	if req.ArticleID == 0 {
		return nil, errors.New("articleID is required")
	}
	comments, err := s.CommentDomain.ListByArticle(ctx, req.ArticleID, req.Page, req.Size)
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
	replies, err := s.CommentDomain.ListByRoot(ctx, req.RootID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &dto.GetCommentRepliesResponse{
		Replies: dto.ToCommentInfos(replies),
		Page:    req.Page,
		Size:    req.Size,
	}, nil
}
