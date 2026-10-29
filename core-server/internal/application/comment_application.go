package application

import (
	"context"
	"errors"
	"strings"

	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"

	"go.uber.org/zap"
)

type CommentService struct {
	CommentDomain domain.CommentRepoDomain
	log           *clog.Log
}

func NewCommentService(commentRepo domain.CommentRepoDomain, log *clog.Log) *CommentService {
	return &CommentService{
		CommentDomain: commentRepo,
		log:           log,
	}
}

// CreateComment 发表一级评论
func (s *CommentService) CreateComment(ctx context.Context, req dto.CreateCommentRequest) error {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		s.log.Debug("empty comment content", zap.Uint64("article_id", req.ArticleID))
		return errors.New("comment content is required")
	}
	if req.ArticleID == 0 || req.UserID == 0 {
		s.log.Debug("invalid comment params",
			zap.Uint64("article_id", req.ArticleID),
			zap.Uint64("user_id", req.UserID),
		)
		return errors.New("articleID and userID are required")
	}
	commentID, err := s.CommentDomain.CreateComment(ctx, &entity.Comment{
		ArticleID:     req.ArticleID,
		PublishUserID: req.UserID,
		Comment:       content,
	})
	if err != nil {
		s.log.Error("create comment failed",
			zap.Error(err),
			zap.Uint64("article_id", req.ArticleID),
			zap.Uint64("user_id", req.UserID),
		)
		return err
	}
	s.log.Info("comment created",
		zap.Uint64("comment_id", commentID),
		zap.Uint64("article_id", req.ArticleID),
		zap.Uint64("user_id", req.UserID),
	)
	return nil
}

// CreateReply 回复某条评论
func (s *CommentService) CreateReply(ctx context.Context, req dto.CreateReplyRequest) error {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		s.log.Debug("empty reply content",
			zap.Uint64("article_id", req.ArticleID),
			zap.Uint64("root_id", req.RootID),
		)
		return errors.New("reply content is required")
	}
	if req.RootID == 0 || req.ReplyToID == 0 || req.ArticleID == 0 || req.UserID == 0 {
		s.log.Debug("invalid reply params",
			zap.Uint64("article_id", req.ArticleID),
			zap.Uint64("root_id", req.RootID),
			zap.Uint64("reply_to_id", req.ReplyToID),
			zap.Uint64("user_id", req.UserID),
		)
		return errors.New("rootID, replyToID, articleID and userID are required")
	}
	commentID, err := s.CommentDomain.CreateReply(ctx, &entity.Comment{
		ArticleID:     req.ArticleID,
		PublishUserID: req.UserID,
		RootID:        req.RootID,
		ReplyToID:     req.ReplyToID,
		Comment:       content,
	})
	if err != nil {
		s.log.Error("create reply failed",
			zap.Error(err),
			zap.Uint64("article_id", req.ArticleID),
			zap.Uint64("root_id", req.RootID),
			zap.Uint64("reply_to_id", req.ReplyToID),
			zap.Uint64("user_id", req.UserID),
		)
		return err
	}
	s.log.Info("reply created",
		zap.Uint64("comment_id", commentID),
		zap.Uint64("article_id", req.ArticleID),
		zap.Uint64("root_id", req.RootID),
		zap.Uint64("user_id", req.UserID),
	)
	return nil
}

// DeleteComment 删除评论（评论作者或文章作者可删）
func (s *CommentService) DeleteComment(ctx context.Context, req dto.DeleteCommentRequest) error {
	if req.ID == 0 || req.UserID == 0 {
		s.log.Debug("invalid delete comment params",
			zap.Uint64("comment_id", req.ID),
			zap.Uint64("user_id", req.UserID),
		)
		return errors.New("id and userID are required")
	}
	if err := s.CommentDomain.SoftDelete(ctx, req.ID, req.UserID); err != nil {
		s.log.Error("delete comment failed",
			zap.Error(err),
			zap.Uint64("comment_id", req.ID),
			zap.Uint64("user_id", req.UserID),
			zap.Uint64("article_id", req.ArticleID),
		)
		return err
	}
	s.log.Info("comment deleted",
		zap.Uint64("comment_id", req.ID),
		zap.Uint64("user_id", req.UserID),
	)
	return nil
}

// GetArticleComments 分页查询文章下的一级评论
func (s *CommentService) GetArticleComments(ctx context.Context, req dto.GetArticleCommentsRequest) (*dto.GetArticleCommentsResponse, error) {
	if req.ArticleID == 0 {
		s.log.Debug("empty article id for comment list")
		return nil, errors.New("articleID is required")
	}
	offset, limit := pagination(req.Page, req.Size)
	comments, err := s.CommentDomain.ListTopByArticle(ctx, req.ArticleID, offset, limit)
	if err != nil {
		s.log.Warn("list article comments failed",
			zap.Error(err),
			zap.Uint64("article_id", req.ArticleID),
		)
		return nil, err
	}
	s.log.Debug("list article comments",
		zap.Uint64("article_id", req.ArticleID),
		zap.Int("count", len(comments)),
	)
	return &dto.GetArticleCommentsResponse{
		Comments: dto.ToCommentInfos(comments),
		Page:     req.Page,
		Size:     req.Size,
	}, nil
}

// GetCommentReplies 分页查询根评论下的回复
func (s *CommentService) GetCommentReplies(ctx context.Context, req dto.GetCommentRepliesRequest) (*dto.GetCommentRepliesResponse, error) {
	if req.RootID == 0 {
		s.log.Debug("empty root id for reply list")
		return nil, errors.New("rootID is required")
	}
	offset, limit := pagination(req.Page, req.Size)
	replies, err := s.CommentDomain.ListRepliesByRoot(ctx, req.RootID, offset, limit)
	if err != nil {
		s.log.Warn("list comment replies failed",
			zap.Error(err),
			zap.Uint64("root_id", req.RootID),
		)
		return nil, err
	}
	s.log.Debug("list comment replies",
		zap.Uint64("root_id", req.RootID),
		zap.Int("count", len(replies)),
	)
	return &dto.GetCommentRepliesResponse{
		Replies: dto.ToCommentInfos(replies),
		Page:    req.Page,
		Size:    req.Size,
	}, nil
}
