package repo

import (
	"context"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/comment"
	"core-server/internal/infras/repo/ent/predicate"
	"core-server/internal/infras/repo/ent/user"
	"core-server/internal/infras/repo/ent/userstatu"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"

	"entgo.io/ent/dialect/sql"
)

type UserRepo struct {
	*EntClient
}

func NewUserRepo(client *EntClient) *UserRepo {
	return &UserRepo{client}
}

// List 分页查询用户列表（按注册时间倒序）
func (r *UserRepo) List(ctx context.Context, page, size uint32) ([]*entity.User, uint64, error) {
	return r.list(ctx, page, size, user.DeletedAtIsNil())
}

// Search 按用户名/邮箱模糊搜索用户
func (r *UserRepo) Search(ctx context.Context, keyword string, page, size uint32) ([]*entity.User, uint64, error) {
	return r.list(ctx, page, size,
		user.DeletedAtIsNil(),
		user.Or(
			user.UsernameContainsFold(keyword),
			user.EmailContainsFold(keyword),
		),
	)
}

// list 分页查询的公共实现，返回当前页数据和总数
func (r *UserRepo) list(ctx context.Context, page, size uint32, preds ...predicate.User) ([]*entity.User, uint64, error) {
	query := r.DB(ctx).User.Query().Where(preds...)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	limit, offset := pageLimitOffset(page, size)
	users, err := query.
		Order(user.ByCreatedAt(sql.OrderDesc())).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return toEntityUsers(users), uint64(total), nil
}

// GetByID 查询单个用户
func (r *UserRepo) GetByID(ctx context.Context, userID uint64) (*entity.User, error) {
	data, err := r.DB(ctx).User.Query().
		Where(user.IDEQ(userID), user.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrorUserNotFound
		}
		return nil, err
	}
	return toEntityUser(data), nil
}

// GetStatu 查询用户统计信息（没有统计记录时返回全 0，保证个人信息页能正常展示）
func (r *UserRepo) GetStatu(ctx context.Context, userID uint64) (*entity.Status, error) {
	data, err := r.DB(ctx).UserStatu.Query().
		Where(userstatu.UserIDEQ(userID), userstatu.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return &entity.Status{UserID: userID}, nil
		}
		return nil, err
	}
	return toEntityUserStatu(data), nil
}

// CountComments 统计用户发表的评论数
func (r *UserRepo) CountComments(ctx context.Context, userID uint64) (uint64, error) {
	count, err := r.DB(ctx).Comment.Query().
		Where(comment.PublishUserIDEQ(userID), comment.DeletedAtIsNil()).
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return uint64(count), nil
}

// UpdateProfile 更新用户资料（用户名为空则不修改，生日为空则清空）
func (r *UserRepo) UpdateProfile(ctx context.Context, data *entity.User) error {
	update := r.DB(ctx).User.UpdateOneID(data.ID).
		SetPhone(data.Phone).
		SetSex(toEntSex(data.Sex)).
		SetSignature(data.Signature)

	if data.Username != "" {
		update.SetUsername(data.Username)
	}
	if data.BirthDate.Valid {
		update.SetBirthDate(data.BirthDate.Time)
	} else {
		update.ClearBirthDate()
	}

	if _, err := update.Save(ctx); err != nil {
		if ent.IsNotFound(err) {
			return ErrorUserNotFound
		}
		return err
	}
	return nil
}

// UpdateAvatar 更新用户头像
func (r *UserRepo) UpdateAvatar(ctx context.Context, userID uint64, avatar string) error {
	if _, err := r.DB(ctx).User.UpdateOneID(userID).SetAvatar(avatar).Save(ctx); err != nil {
		if ent.IsNotFound(err) {
			return ErrorUserNotFound
		}
		return err
	}
	return nil
}

// UpdateStatus 更新用户封禁状态
func (r *UserRepo) UpdateStatus(ctx context.Context, userID uint64, status enum.Status) error {
	if _, err := r.DB(ctx).User.UpdateOneID(userID).SetStatus(toEntStatus(status)).Save(ctx); err != nil {
		if ent.IsNotFound(err) {
			return ErrorUserNotFound
		}
		return err
	}
	return nil
}
