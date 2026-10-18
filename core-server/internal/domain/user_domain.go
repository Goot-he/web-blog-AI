package domain

import (
	"context"

	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type UserRepo interface {
	// List 分页查询用户列表（按注册时间倒序）
	List(ctx context.Context, page, size uint32) ([]*entity.User, uint64, error)
	// Search 按用户名/邮箱模糊搜索用户
	Search(ctx context.Context, keyword string, page, size uint32) ([]*entity.User, uint64, error)
	// GetByID 查询单个用户
	GetByID(ctx context.Context, userID uint64) (*entity.User, error)
	// GetStatu 查询用户统计信息（没有记录时返回全 0 的统计）
	GetStatu(ctx context.Context, userID uint64) (*entity.Status, error)
	// CountComments 统计用户发表的评论数
	CountComments(ctx context.Context, userID uint64) (uint64, error)
	// UpdateProfile 更新用户资料
	UpdateProfile(ctx context.Context, user *entity.User) error
	// UpdateAvatar 更新用户头像
	UpdateAvatar(ctx context.Context, userID uint64, avatar string) error
	// UpdateStatus 更新用户封禁状态
	UpdateStatus(ctx context.Context, userID uint64, status enum.Status) error
}

type UserDomainRepo interface {
	UserRepo
	ITransaction
}
