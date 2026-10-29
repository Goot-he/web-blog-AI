package repo

import (
	"context"
	"core-server/internal/infras/repo/ent/interactionlike"
	"core-server/internal/model/entity"
	stdsql "database/sql"
	"fmt"

	"entgo.io/ent/dialect/sql"
)

type LikeRepo struct {
	*EntClient
}

func NewLikeRepo(client *EntClient) *LikeRepo {
	return &LikeRepo{client}
}

// Upsert 点赞的repo层函数 返回0表示影响的行数为0，即失败了 返回1表示成功了
func (r *LikeRepo) Upsert(ctx context.Context, like *entity.InteractionLike) (int, error) {
	const query = `
		INSERT INTO interaction_like (user_id, object_type, object_id, status, version)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		  status = IF(status = 'thumb_up' AND version >= VALUES(version), status, VALUES(status)),
		  version = IF(status = 'thumb_up' AND version >= VALUES(version), version, VALUES(version)),
		  updated_at = IF(status = 'thumb_up' AND version >= VALUES(version), updated_at, NOW(3))
	`
	driver, ok := r.driver.(interface {
		ExecContext(context.Context, string, ...any) (stdsql.Result, error)
	})
	if !ok {
		return 0, fmt.Errorf("ent driver does not support ExecContext")
	}
	result, err := driver.ExecContext(ctx, query, like.UserID, like.ObjectType.String(), like.ObjectID, like.LikeStatus.String(), like.Version, entity.LikeStatusTypeThumbUp.String())
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return 0, err
	}
	return 1, nil
}

// UpdateWithCondition 更新点赞状态的函数
func (r *LikeRepo) UpdateWithCondition(ctx context.Context, condition string, like entity.InteractionLike) (int, error) {
	// condition 是修改前的记录状态
	affected, err := r.DB(ctx).InteractionLike.Update().
		Where(
			interactionlike.UserIDEQ(like.UserID),
			interactionlike.ObjectIDEQ(like.ObjectID),
			interactionlike.ObjectTypeEQ(interactionlike.ObjectType(like.ObjectType.String())),
			interactionlike.StatusEQ(interactionlike.Status(condition)),
			interactionlike.VersionLTE(like.Version),
		).
		SetStatus(interactionlike.Status(like.LikeStatus.String())).
		SetVersion(like.Version).
		Save(ctx)
	if err != nil || affected == 0 {
		return 0, err
	}
	return affected, nil
}

// QueryWithCondition 返回一个点赞实例 用于判断指定的点赞状态是否存在
func (r *LikeRepo) QueryWithCondition(ctx context.Context, userID, objectID uint64, objectType string, status string) (*entity.InteractionLike, error) {
	node, err := r.DB(ctx).InteractionLike.Query().Where(
		interactionlike.UserIDEQ(userID),
		interactionlike.ObjectIDEQ(objectID),
		interactionlike.ObjectTypeEQ(interactionlike.ObjectType(objectType)),
		interactionlike.StatusEQ(interactionlike.Status(status)),
	).Only(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityLike(node), nil
}

// CountUserLiked 统计指定用户的指定类型的点赞数量
func (r *LikeRepo) CountUserLiked(ctx context.Context, userID, objectID uint64, objectType string) (int, error) {
	count, err := r.DB(ctx).InteractionLike.Query().Where(
		interactionlike.UserIDEQ(userID),
		interactionlike.ObjectTypeEQ(interactionlike.ObjectType(objectType)),
		interactionlike.ObjectIDEQ(objectID),
	).Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// PageQueryLikeObjects 分页返回批量数据
func (r *LikeRepo) PageQueryLikeObjects(ctx context.Context, userID uint64, objectType string, offset, limit int) ([]*entity.InteractionLike, error) {
	nodes, err := r.DB(ctx).InteractionLike.Query().Where(
		interactionlike.UserIDEQ(userID),
		interactionlike.ObjectTypeEQ(interactionlike.ObjectType(objectType))).
		Order(
			interactionlike.ByVersion(sql.OrderDesc()),
			interactionlike.ByUpdatedAt(sql.OrderDesc()),
			interactionlike.ByID(sql.OrderDesc()),
		).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return toEntityLikes(nodes), nil
}
