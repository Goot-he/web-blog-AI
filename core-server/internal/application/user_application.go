package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"core-server/internal/domain"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type UserService struct {
	UserDomain domain.UserDomainRepo
}

func NewUserService(userRepo domain.UserDomainRepo) *UserService {
	return &UserService{
		UserDomain: userRepo,
	}
}

// ListUsers 分页查询用户列表
func (s *UserService) ListUsers(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	users, total, err := s.UserDomain.List(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	return &dto.ListUsersResponse{
		Users: dto.ToUserInfos(users),
		Total: total,
	}, nil
}

// SearchUsers 按关键词搜索用户（空关键词直接返回空列表）
func (s *UserService) SearchUsers(ctx context.Context, req dto.SearchUsersRequest) (*dto.SearchUsersResponse, error) {
	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return &dto.SearchUsersResponse{}, nil
	}
	users, total, err := s.UserDomain.Search(ctx, keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	return &dto.SearchUsersResponse{
		Users: dto.ToUserInfos(users),
		Total: total,
	}, nil
}

// GetProfile 查询用户个人主页（资料 + 统计信息）
func (s *UserService) GetProfile(ctx context.Context, req dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	user, err := s.UserDomain.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	statu, err := s.UserDomain.GetStatu(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	commentCount, err := s.UserDomain.CountComments(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return dto.ToGetProfileResponse(user, statu, commentCount), nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(ctx context.Context, req dto.UpdateProfileRequest) error {
	data := &entity.User{
		ID:        req.UserID,
		Username:  strings.TrimSpace(req.Username),
		Phone:     req.Phone,
		Sex:       enum.ParseSex(req.Sex),
		Signature: req.Signature,
	}

	birthday := strings.TrimSpace(req.Birthday)
	if birthday != "" {
		t, err := time.Parse("2006-01-02", birthday)
		if err != nil {
			return fmt.Errorf("invalid birthday %q, expected yyyy-MM-dd", req.Birthday)
		}
		data.BirthDate = sql.NullTime{Time: t, Valid: true}
	}

	return s.UserDomain.UpdateProfile(ctx, data)
}

// UpdateAvatar 更新用户头像，返回新的头像地址
func (s *UserService) UpdateAvatar(ctx context.Context, req dto.UpdateAvatarRequest) (string, error) {
	avatar := strings.TrimSpace(req.Avatar)
	if avatar == "" {
		return "", errors.New("avatar is required")
	}
	if err := s.UserDomain.UpdateAvatar(ctx, req.UserID, avatar); err != nil {
		return "", err
	}
	return avatar, nil
}

// UpdateUserStatus 更新用户封禁状态
func (s *UserService) UpdateUserStatus(ctx context.Context, req dto.UpdateUserStatusRequest) error {
	status := enum.ParseStatus(req.Status)
	if status != enum.StatusApproved && status != enum.StatusBlocked {
		return fmt.Errorf("invalid status %q, expected approved or blocked", req.Status)
	}
	return s.UserDomain.UpdateStatus(ctx, req.UserID, status)
}
