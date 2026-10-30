package application

import (
	"context"
	"core-server/internal/infras/repo"
	"core-server/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type UserService struct {
	repo domain.UserDomainRepo
	log  *clog.Log
}

func NewUserService(userRepo domain.UserDomainRepo, log *clog.Log) *UserService {
	return &UserService{
		repo: userRepo,
		log:  log,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserBlocked        = errors.New("USER_BLOCKED")
	ErrEmailAlreadyInUse  = errors.New("email is already registered")
	ErrPasswordMismatch   = errors.New("passwords do not match")
)

// 用户登入方面

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1.Does the user exist?
	user, err := s.repo.GetByName(ctx, req.Username)
	if err != nil {
		// DB error or Not found
		s.log.Error("UserService/Login error:")
		return nil, err
	}

	// 2. Determine whether the user has been blocked.
	if user.Status != enum.StatusApproved {
		return nil, ErrUserBlocked
	}

	// 3. Check whether the password is correct
	if user.Password != utils.Bcrypt(req.Password) {
		return nil, ErrInvalidCredentials
	}

	return dto.ToLoginResponse(user), nil
}

func (s *UserService) EmailLogin(ctx context.Context, req *dto.EmailLoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.log.Error("UserService/EmailLogin error:")
		return nil, err
	}

	if user.Status != enum.StatusApproved {
		return nil, ErrUserBlocked
	}
	return dto.ToLoginResponse(user), nil
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, repo.ErrorUserNotFound) {
		s.log.Error("UserService/Register error:")
		return err
	}
	if existing != nil {
		return ErrEmailAlreadyInUse
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: utils.Bcrypt(req.Password),
		Role:     enum.UserRoleUser,
		Status:   enum.StatusApproved,
		Sex:      enum.Unknown,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.log.Error("UserService/Register error:")
		return err
	}

	return nil
}

func (s *UserService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	// 1.校验两次密码是否相同
	if req.Password != req.Confirm {
		return ErrPasswordMismatch
	}

	// 2.判断该邮箱是否存在
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.log.Error("UserService/ForgotPassword error:")
		return err
	}

	// 3.更新密码
	err = s.repo.UpdatePassword(ctx, user.ID, utils.Bcrypt(req.Password))
	if err != nil {
		s.log.Error("UserService/ForgotPassword error:")
		return err
	}

	return nil
}

// =====================================================================================================================

// ListUsers 分页查询用户列表
func (s *UserService) ListUsers(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	offset, limit := pagination(req.Page, req.PageSize)
	users, total, err := s.repo.List(ctx, int32(limit), int32(offset))
	if err != nil {
		s.log.Warn("list users failed")
		return nil, err
	}
	s.log.Debug("list users")
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
	offset, limit := pagination(req.Page, req.PageSize)
	users, total, err := s.repo.Search(ctx, keyword, int32(limit), int32(offset))
	if err != nil {
		s.log.Warn("search users failed")
		return nil, err
	}
	s.log.Debug("search users")
	return &dto.SearchUsersResponse{
		Users: dto.ToUserInfos(users),
		Total: total,
	}, nil
}

// GetProfile 查询用户个人主页（资料 + 统计信息）
func (s *UserService) GetProfile(ctx context.Context, req dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	agg, err := s.repo.GetUserMsgByID(ctx, req.UserID)
	if err != nil {
		s.log.Warn("get user profile failed")
		return nil, err
	}

	var commentCount uint64
	if agg.Statu != nil {
		commentCount = agg.Statu.CommentCount
	}
	return dto.ToGetProfileResponse(agg.User, agg.Statu, commentCount), nil
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
			s.log.Debug("invalid birthday")
			return fmt.Errorf("invalid birthday %q, expected yyyy-MM-dd", req.Birthday)
		}
		data.BirthDate = sql.NullTime{Time: t, Valid: true}
	}

	if err := s.repo.UpdateProfile(ctx, data); err != nil {
		s.log.Error("update user profile failed")
		return err
	}
	s.log.Info("user profile updated")
	return nil
}

// UpdateAvatar 更新用户头像，返回新的头像地址
func (s *UserService) UpdateAvatar(ctx context.Context, req dto.UpdateAvatarRequest) (string, error) {
	avatar := strings.TrimSpace(req.Avatar)
	if avatar == "" {
		s.log.Debug("empty avatar")
		return "", errors.New("avatar is required")
	}
	if err := s.repo.UpdateAvatar(ctx, req.UserID, avatar); err != nil {
		s.log.Error("update user avatar failed")
		return "", err
	}
	s.log.Info("user avatar updated")
	return avatar, nil
}

// UpdateUserStatus 更新用户封禁状态
func (s *UserService) UpdateUserStatus(ctx context.Context, req dto.UpdateUserStatusRequest) error {
	status := enum.ParseStatus(req.Status)
	if status != enum.StatusApproved && status != enum.StatusBlocked {
		s.log.Debug("invalid user status")
		return fmt.Errorf("invalid status %q, expected approved or blocked", req.Status)
	}
	if err := s.repo.UpdateStatus(ctx, req.UserID, status); err != nil {
		s.log.Error("update user status failed")
		return err
	}
	s.log.Info("user status updated")
	return nil
}
