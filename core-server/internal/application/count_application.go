package application

import (
	"context"

	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type CountService struct {
	repo       domain.CountRepoDomain
	userDomain domain.UserDomainRepo
	log        *clog.Log
}

func NewCountService(countRepo domain.CountRepoDomain, userDomain domain.UserDomainRepo, log *clog.Log) *CountService {
	return &CountService{
		repo:       countRepo,
		userDomain: userDomain,
		log:        log,
	}
}

// GetArticleInteractionCount 查询单篇文章的全部互动计数
func (s *CountService) GetArticleInteractionCount(ctx context.Context, articleID uint64) (*entity.InteractionStatus, error) {
	res, err := s.repo.GetByObject(ctx, enum.ObjectTypeArticle.String(), articleID)
	if err != nil {
		s.log.Warn("get article interaction count failed")
		return nil, err
	}

	counts := &entity.InteractionStatus{}
	for _, v := range res {
		switch v.InteractionType {
		case enum.InteractionTypeLike:
			counts.LikeCount = v.Count
		case enum.InteractionTypeComment:
			counts.CommentCount = v.Count
		case enum.InteractionTypeFavor:
			counts.FavorCount = v.Count
		case enum.InteractionTypeFollow:
			counts.FollowCount = v.Count
		case enum.InteractionTypeView:
			counts.ViewCount = v.Count
		}
	}

	return counts, nil
}

// BatchGetArticleInteractionCounts 批量查询多篇文章的互动计数
func (s *CountService) BatchGetArticleInteractionCounts(ctx context.Context, articleIDs []uint64) (map[uint64]*entity.InteractionStatus, error) {
	countsByIds := make(map[uint64]*entity.InteractionStatus, len(articleIDs))
	for _, v := range articleIDs {
		countsByIds[v] = &entity.InteractionStatus{}
	}

	res, err := s.repo.GetByObjects(ctx, enum.ObjectTypeArticle.String(), articleIDs)
	if err != nil {
		s.log.Warn("batch get article interaction counts failed")
		return nil, err
	}
	for _, count := range res {
		status, ok := countsByIds[count.ID]
		if !ok {
			continue
		}
		switch count.InteractionType {
		case enum.InteractionTypeLike:
			status.LikeCount = count.Count
		case enum.InteractionTypeComment:
			status.CommentCount = count.Count
		case enum.InteractionTypeFollow:
			status.FollowCount = count.Count
		case enum.InteractionTypeFavor:
			status.FavorCount = count.Count
		case enum.InteractionTypeView:
			status.ViewCount = count.Count
		}
	}
	return countsByIds, nil
}

// ------------------------------------------------ 计数修改 -----------------------------------------

// UpdateObjectTypeWithInteractionType 按对象类型 + 互动类型增减计数
func (s *CountService) UpdateObjectTypeWithInteractionType(ctx context.Context, object, interaction string, objectID uint64, delta int64) error {
	if err := s.repo.Upsert(ctx, &entity.InteractionCount{
		ObjectID:        objectID,
		ObjectType:      enum.ParseObjectType(object),
		InteractionType: enum.ParseInteractionType(interaction),
	}, delta); err != nil {
		s.log.Error("update interaction count failed")
		return err
	}
	s.log.Debug("interaction count updated")
	return nil
}

// UpdateLikeCount 更新用户点赞数
func (s *CountService) UpdateLikeCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateLikeCount(ctx, userID, num); err != nil {
		s.log.Error("update user like count failed")
		return err
	}
	s.log.Debug("user like count updated")
	return nil
}

// UpdateReceiveLikeCount 更新用户收到的点赞数
func (s *CountService) UpdateReceiveLikeCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateReceiveLikeCount(ctx, userID, num); err != nil {
		s.log.Error("update user receive like count failed")
		return err
	}
	s.log.Debug("user receive like count updated")
	return nil
}

// UpdateCommentCount 更新用户评论数
func (s *CountService) UpdateCommentCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateCommentCount(ctx, userID, num); err != nil {
		s.log.Error("update user comment count failed")
		return err
	}
	s.log.Debug("user comment count updated")
	return nil
}

// UpdateReviewCount 更新用户文章浏览量
func (s *CountService) UpdateReviewCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateViewCount(ctx, userID, num); err != nil {
		s.log.Error("update user view count failed")
		return err
	}
	s.log.Debug("user view count updated")
	return nil
}

// UpdateFavorCount 更新用户收藏数
func (s *CountService) UpdateFavorCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateFavorCount(ctx, userID, num); err != nil {
		s.log.Error("update user favor count failed")
		return err
	}
	s.log.Debug("user favor count updated")
	return nil
}

// UpdateReceiveFavorCount 更新用户收到的收藏数
func (s *CountService) UpdateReceiveFavorCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateReceiveFavorCount(ctx, userID, num); err != nil {
		s.log.Error("update user receive favor count failed")
		return err
	}
	s.log.Debug("user receive favor count updated")
	return nil
}
