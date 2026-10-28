package application

import (
	"context"
	"core-server/internal/config"
	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
	"time"

	"github.com/avast/retry-go"
)

type LikeService struct {
	cfg          *config.Config
	log          *clog.Log
	repo         domain.LikeDomainRepo
	articleRepo  domain.ArticleDomainRepo
	userRepo     domain.UserDomainRepo
	countService *CountService
}

func NewLikeService(
	log *clog.Log,
	repo domain.LikeDomainRepo,
	articleRepo domain.ArticleDomainRepo,
	userRepo domain.UserDomainRepo,
	countService *CountService,
	cfg *config.Config,
) (*LikeService, error) {
	return &LikeService{
		cfg:          cfg,
		repo:         repo,
		articleRepo:  articleRepo,
		userRepo:     userRepo,
		countService: countService,
		log:          log,
	}, nil
}

func (s *LikeService) HasThumbUp(ctx context.Context, userID, objectID uint64, objectType enum.ObjectType) (bool, error) {
	like, err := s.repo.QueryWithCondition(ctx, userID, objectID, objectType, entity.LikeStatusTypeThumbUp.String())
	if err != nil {
		return false, err
	}
	return like != nil, nil
}

// ------------------------------------------------------

func (s *LikeService) ThumbUp(ctx context.Context, userID, objectID uint64, objectType enum.ObjectType) error {
	// 先判断是否已经存在点赞状态
	isThumbUp, err := s.HasThumbUp(ctx, userID, objectID, objectType)
	if err != nil {
		return err
	}
	if isThumbUp { // 这里表示已经处于点赞状态了
		return nil
	}

	like := &entity.InteractionLike{
		ObjectID:   objectID,
		UserID:     userID,
		ObjectType: objectType,
		LikeStatus: entity.LikeStatusTypeThumbUp,
		Version:    time.Now().Unix(),
	}

	mark := true

	if err := s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		affected, err := s.repo.Upsert(ctx, like)
		if err != nil {
			return err
		}
		if affected == 0 { // 这里表示数据库未执行任何操作
			mark = false
			return nil
		}

		if err := s.countService.UpdateObjectTypeWithInteractionType(ctx, objectType.String(), enum.InteractionTypeLike.String(), objectID, 1); err != nil {
			return err
		}

		// 更新用户的点赞数
		if err := s.userRepo.UpdateLikeCount(ctx, userID, 1); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	// 如果没有问题 异步增加此作者收到的点赞总数
	if mark {
		s.asyncUpdateAuthorReceiveLikeCount(objectType, objectID, 1)
	}
	return nil
}

func (s *LikeService) CancelThumbUp(ctx context.Context, userID, objectID uint64, objectType enum.ObjectType) error {
	// 先判断是否处于点赞状态
	isThumbUp, err := s.HasThumbUp(ctx, userID, objectID, objectType)
	if err != nil {
		return err
	}
	if !isThumbUp {
		return nil
	}

	mark := true

	like := &entity.InteractionLike{
		ObjectID:   objectID,
		UserID:     userID,
		ObjectType: objectType,
		LikeStatus: entity.LikeStatusTypeNothing,
		Version:    time.Now().Unix(),
	}

	if err := s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		affected, err := s.repo.Upsert(ctx, like)
		if err != nil {
			return err
		}
		if affected == 0 {
			mark = false
			return nil
		}

		if err := s.countService.UpdateObjectTypeWithInteractionType(ctx, objectType.String(), enum.InteractionTypeLike.String(), objectID, -1); err != nil {
			return err
		}

		// 更新用户的点赞数量
		if err := s.userRepo.UpdateLikeCount(ctx, userID, -1); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	// 这里异步更新作者收到的点赞数
	if mark {
		s.asyncUpdateAuthorReceiveLikeCount(objectType, objectID, -1)
	}
	return nil
}

func (s *LikeService) asyncUpdateAuthorReceiveLikeCount(objectType enum.ObjectType, objectID uint64, delta int) {
	go func() {
		syncContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := retry.Do(
			func() error {
				authorID, err := s.getAuthorIdByObject(syncContext, objectType, objectID)
				if err != nil {
					return err
				}
				return s.userRepo.UpdateReceiveLikeCount(syncContext, authorID, int64(delta))
			}, retry.Attempts(3),
			retry.MaxDelay(10*time.Second),
			retry.DelayType(retry.BackOffDelay),
		)
		if err != nil {
			return
		}
	}()
}

func (s *LikeService) getAuthorIdByObject(ctx context.Context, objectType enum.ObjectType, objectID uint64) (uint64, error) {
	switch objectType {
	case enum.ObjectTypeArticle:
		article, err := s.articleRepo.GetArticle(ctx, objectID)
		if err != nil {
			return 0, err
		}
		return article.Article.AuthorID, nil

	default:
		return 0, nil
	}
}

// -------------------------------------------------------

func (s *LikeService) UserLikeList(ctx context.Context, userID uint64, objectType enum.ObjectType, page, size int) ([]*aggregate.ArticleAggregate, error) {
	// 大致流程是线 先拿取点赞列表 再通过点赞列表拿取文章
	limit := (page - 1) * size
	likes, err := s.repo.PageQueryLikeObjects(ctx, userID, objectType, limit, size)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(likes))
	for _, like := range likes {
		ids = append(ids, like.ObjectID)
	}
	articles, err := s.articleRepo.ListByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
