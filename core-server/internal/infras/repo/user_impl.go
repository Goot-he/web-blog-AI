package repo

import (
	"context"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/predicate"
	"core-server/internal/infras/repo/ent/user"
	"core-server/internal/infras/repo/ent/userstatu"
	"core-server/internal/model/aggregate"
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

// GetByID 按 ID 查询用户
func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	return r.getOne(ctx, user.IDEQ(id))
}

// GetByName 按用户名查询用户
func (r *UserRepo) GetByName(ctx context.Context, name string) (*entity.User, error) {
	return r.getOne(ctx, user.UsernameEQ(name))
}

// GetByEmail 按邮箱查询用户
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.getOne(ctx, user.EmailEQ(email))
}

// getOne 按条件查询单个未删除用户
func (r *UserRepo) getOne(ctx context.Context, preds ...predicate.User) (*entity.User, error) {
	data, err := r.DB(ctx).User.Query().
		Where(append(preds, user.DeletedAtIsNil())...).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrorUserNotFound
		}
		return nil, err
	}
	return toEntityUser(data), nil
}

// GetUserMsgByID 查询用户资料和统计信息
func (r *UserRepo) GetUserMsgByID(ctx context.Context, id uint64) (*aggregate.UserAggregate, error) {
	u, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	statu, err := r.GetUserStat(ctx, id)
	if err != nil {
		return nil, err
	}
	return &aggregate.UserAggregate{
		User:  u,
		Statu: statu,
	}, nil
}

// GetUserStat 查询用户统计信息（没有记录时返回全 0）
func (r *UserRepo) GetUserStat(ctx context.Context, userID uint64) (*entity.UserStatus, error) {
	data, err := r.DB(ctx).UserStatu.Query().
		Where(userstatu.UserIDEQ(userID), userstatu.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return &entity.UserStatus{UserID: userID}, nil
		}
		return nil, err
	}
	return toEntityUserStatu(data), nil
}

// Search 按用户名/邮箱模糊搜索用户，返回当前页数据和总数
func (r *UserRepo) Search(ctx context.Context, keyword string, limit, offset int32) ([]*entity.User, uint64, error) {
	return r.list(ctx, limit, offset,
		user.Or(
			user.UsernameContainsFold(keyword),
			user.EmailContainsFold(keyword),
		),
	)
}

// List 分页查询用户列表，返回当前页数据和总数
func (r *UserRepo) List(ctx context.Context, limit, offset int32) ([]*entity.User, uint64, error) {
	return r.list(ctx, limit, offset)
}

// list 用户分页查询的公共实现
func (r *UserRepo) list(ctx context.Context, limit, offset int32, preds ...predicate.User) ([]*entity.User, uint64, error) {
	query := r.DB(ctx).User.Query().
		Where(append(preds, user.DeletedAtIsNil())...)

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	l, o := limitOffset(limit, offset)
	users, err := query.
		Order(user.ByCreatedAt(sql.OrderDesc())).
		Offset(o).
		Limit(l).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return toEntityUsers(users), uint64(total), nil
}

// ListByIDs 按 ID 批量查询用户
func (r *UserRepo) ListByIDs(ctx context.Context, ids []uint64) ([]*entity.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	users, err := r.DB(ctx).User.Query().
		Where(user.IDIn(ids...), user.DeletedAtIsNil()).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityUsers(users), nil
}

// CreateUser 新建用户
func (r *UserRepo) CreateUser(ctx context.Context, data *entity.User) error {
	create := r.DB(ctx).User.Create().
		SetUsername(data.Username).
		SetPassword(data.Password).
		SetEmail(data.Email).
		SetPhone(data.Phone).
		SetAvatar(data.Avatar).
		SetSex(toEntSex(data.Sex)).
		SetRole(toEntRole(data.Role)).
		SetSignature(data.Signature).
		SetStatus(toEntStatus(data.Status))

	if data.BirthDate.Valid {
		create.SetBirthDate(data.BirthDate.Time)
	}

	created, err := create.Save(ctx)
	if err != nil {
		return err
	}

	data.ID = created.ID
	data.CreatedAt = created.CreatedAt
	data.UpdatedAt = created.UpdatedAt
	return nil
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

// UpdatePassword 更新用户密码（调用方负责加密）
func (r *UserRepo) UpdatePassword(ctx context.Context, userID uint64, password string) error {
	if _, err := r.DB(ctx).User.UpdateOneID(userID).SetPassword(password).Save(ctx); err != nil {
		if ent.IsNotFound(err) {
			return ErrorUserNotFound
		}
		return err
	}
	return nil
}

// UpdateStatus 更新用户封禁状态
func (r *UserRepo) UpdateStatus(ctx context.Context, userID uint64, status enum.UserStatus) error {
	if _, err := r.DB(ctx).User.UpdateOneID(userID).SetStatus(toEntStatus(status)).Save(ctx); err != nil {
		if ent.IsNotFound(err) {
			return ErrorUserNotFound
		}
		return err
	}
	return nil
}

// UpdateLikeCount 更新用户点赞数
func (r *UserRepo) UpdateLikeCount(ctx context.Context, userID uint64, delta int64) error {
	return r.updateStat(ctx, userID, delta,
		func(u *ent.UserStatuUpdate) *ent.UserStatuUpdate { return u.AddLikeCount(delta) },
		func(c *ent.UserStatuCreate) *ent.UserStatuCreate { return c.SetLikeCount(uint64(delta)) },
	)
}

// UpdateReceiveLikeCount 更新用户收到的点赞数
func (r *UserRepo) UpdateReceiveLikeCount(ctx context.Context, userID uint64, delta int64) error {
	return r.updateStat(ctx, userID, delta,
		func(u *ent.UserStatuUpdate) *ent.UserStatuUpdate { return u.AddReceiveLikeCount(delta) },
		func(c *ent.UserStatuCreate) *ent.UserStatuCreate { return c.SetReceiveLikeCount(uint64(delta)) },
	)
}

// UpdateFavorCount 更新用户收藏数
func (r *UserRepo) UpdateFavorCount(ctx context.Context, userID uint64, delta int64) error {
	return r.updateStat(ctx, userID, delta,
		func(u *ent.UserStatuUpdate) *ent.UserStatuUpdate { return u.AddFavorCount(delta) },
		func(c *ent.UserStatuCreate) *ent.UserStatuCreate { return c.SetFavorCount(uint64(delta)) },
	)
}

// UpdateReceiveFavorCount 更新用户收到的收藏数
func (r *UserRepo) UpdateReceiveFavorCount(ctx context.Context, userID uint64, delta int64) error {
	return r.updateStat(ctx, userID, delta,
		func(u *ent.UserStatuUpdate) *ent.UserStatuUpdate { return u.AddFavorCount(delta) },
		func(c *ent.UserStatuCreate) *ent.UserStatuCreate { return c.SetFavorCount(uint64(delta)) },
	)
}

// UpdateArticleCount 更新用户文章数
func (r *UserRepo) UpdateArticleCount(ctx context.Context, userID uint64, delta int64) error {
	return r.updateStat(ctx, userID, delta,
		func(u *ent.UserStatuUpdate) *ent.UserStatuUpdate { return u.AddArticleCount(delta) },
		func(c *ent.UserStatuCreate) *ent.UserStatuCreate { return c.SetArticleCount(uint64(delta)) },
	)
}

// UpdateViewCount 更新用户文章被浏览量
func (r *UserRepo) UpdateViewCount(ctx context.Context, userID uint64, delta int64) error {
	return r.updateStat(ctx, userID, delta,
		func(u *ent.UserStatuUpdate) *ent.UserStatuUpdate { return u.AddViewCount(delta) },
		func(c *ent.UserStatuCreate) *ent.UserStatuCreate { return c.SetViewCount(uint64(delta)) },
	)
}

// UpdateCommentCount 更新用户评论数
func (r *UserRepo) UpdateCommentCount(ctx context.Context, userID uint64, delta int64) error {
	return r.updateStat(ctx, userID, delta,
		func(u *ent.UserStatuUpdate) *ent.UserStatuUpdate { return u.AddCommentCount(delta) },
		func(c *ent.UserStatuCreate) *ent.UserStatuCreate { return c.SetCommentCount(uint64(delta)) },
	)
}

// updateStat 更新统计字段；统计行不存在且是正向增量时补建一行
func (r *UserRepo) updateStat(
	ctx context.Context,
	userID uint64,
	delta int64,
	add func(*ent.UserStatuUpdate) *ent.UserStatuUpdate,
	create func(*ent.UserStatuCreate) *ent.UserStatuCreate,
) error {
	affected, err := add(r.DB(ctx).UserStatu.Update().
		Where(userstatu.UserIDEQ(userID), userstatu.DeletedAtIsNil())).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 && delta > 0 {
		_, err = create(r.DB(ctx).UserStatu.Create().SetUserID(userID)).Save(ctx)
		return err
	}
	return nil
}
