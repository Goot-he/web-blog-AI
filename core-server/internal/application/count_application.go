package application

import (
	"context"
	"core-server/internal/domain"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type CountService struct {
	repo       domain.CountRepoDomain
	userDomain domain.UserDomainRepo
}

func NewCountService(countRepo domain.CountRepoDomain) *CountService {
	return &CountService{
		repo: countRepo,
	}
}

func (s *CountService) GetArticleInteractionCount(ctx context.Context, articleID uint64) (*entity.InteractionStatus, error) {
	res, err := s.repo.GetByObject(ctx, enum.ObjectTypeArticle, articleID)
	if err != nil {
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

func (s *CountService) BatchGetArticleInteractionCounts(ctx context.Context, articleIDs []uint64) (map[uint64]*entity.InteractionStatus, error) {
	countsByIds := make(map[uint64]*entity.InteractionStatus, len(articleIDs))
	for _, v := range articleIDs {
		countsByIds[v] = &entity.InteractionStatus{}
	}

	res, err := s.repo.GetByObjects(ctx, enum.ObjectTypeArticle, articleIDs)
	if err != nil {
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
		case enum.InteractionTypeFavor:
			status.FavorCount = count.Count
		case enum.InteractionTypeFollow:
			status.FollowCount = count.Count
		case enum.InteractionTypeView:
			status.ViewCount = count.Count
		}
	}
	return countsByIds, nil
}

// ------------------------------------------------ 计数修改 -----------------------------------------

func (s *CountService) UpdateObjectTypeWithInteractionType(ctx context.Context, object, interaction string, objectID uint64, delta int64) error {
	if err := s.repo.Upsert(ctx, &entity.InteractionCount{
		ObjectID:        objectID,
		ObjectType:      enum.ParseObjectType(object),
		InteractionType: enum.ParseInteractionType(interaction),
	}, delta); err != nil {
		return err
	}

	return nil
}

func (s *CountService) UpdateLikeCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateLikeCount(ctx, userID, num); err != nil {
		return err
	}
	return nil
}

func (s *CountService) UpdateReceiveLikeCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateReceiveLikeCount(ctx, userID, num); err != nil {
		return err
	}
	return nil
}

func (s *CountService) UpdateCommentCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateCommentCount(ctx, userID, num); err != nil {
		return err
	}
	return nil
}

func (s *CountService) UpdateReviewCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateViewCount(ctx, userID, num); err != nil {
		return err
	}
	return nil
}

func (s *CountService) UpdateFavorCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateFavorCount(ctx, userID, num); err != nil {
		return err
	}
	return nil
}

func (s *CountService) UpdateReceiveFavorCount(ctx context.Context, userID uint64, num int64) error {
	if err := s.userDomain.UpdateReceiveFavorCount(ctx, userID, num); err != nil {
		return err
	}
	return nil
}
