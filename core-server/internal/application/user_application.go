package application

import (
	"context"
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

	"go.uber.org/zap"
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

// 用户登入方面

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1.Does the user exist?
	user, err := s.repo.GetByName(ctx, req.Username)
	if err != nil {
		// DB error or Not found
		s.log.Error("UserService/Login error:", zap.Error(err))
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
		s.log.Error("UserService/EmailLogin error:", zap.Error(err))
		return nil, err
	}

	if user.Status != enum.StatusApproved {
		return nil, ErrUserBlocked
	}
	return dto.ToLoginResponse(user), nil
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, repo.ErrUserNotFound) {
		s.log.Error("UserService/Register error:", zap.Error(err))
		return err
	}
	if existing != nil {
		return ErrEmailAlreadyInUse
	}

	user := &entity.User{
		Name:     req.Username,
		Email:    req.Email,
		Password: utils.Bcrypt(req.Password),
		Role:     enum.UserRoleUser,
		Status:   enum.StatusApproved,
		Sex:      enum.UserUnknown,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.log.Error("UserService/Register error:", zap.Error(err))
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
		s.log.Error("UserService/ForgotPassword error:", zap.Error(err))
		return err
	}

	// 3.更新密码
	err = s.repo.UpdatePassword(ctx, user.ID, utils.Bcrypt(req.Password))
	if err != nil {
		s.log.Error("UserService/ForgotPassword error:", zap.Error(err))
		return err
	}

	return nil
}

// =====================================================================================================================

// ListUsers 分页查询用户列表
func (s *UserService) ListUsers(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	offset, limit := pagination(req.Page, req.PageSize)
	users, total, err := s.UserDomain.List(ctx, int32(limit), int32(offset))
	if err != nil {
		s.log.Warn("list users failed",
			zap.Error(err),
			zap.Uint32("page", req.Page),
			zap.Uint32("page_size", req.PageSize),
		)
		return nil, err
	}
	s.log.Debug("list users", zap.Int("count", len(users)), zap.Uint64("total", total))
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
	users, total, err := s.UserDomain.Search(ctx, keyword, int32(limit), int32(offset))
	if err != nil {
		s.log.Warn("search users failed",
			zap.Error(err),
			zap.String("keyword", keyword),
			zap.Uint32("page", req.Page),
			zap.Uint32("page_size", req.PageSize),
		)
		return nil, err
	}
	s.log.Debug("search users",
		zap.String("keyword", keyword),
		zap.Int("count", len(users)),
		zap.Uint64("total", total),
	)
	return &dto.SearchUsersResponse{
		Users: dto.ToUserInfos(users),
		Total: total,
	}, nil
}

// GetProfile 查询用户个人主页（资料 + 统计信息）
func (s *UserService) GetProfile(ctx context.Context, req dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	agg, err := s.UserDomain.GetUserMsgByID(ctx, req.UserID)
	if err != nil {
		s.log.Warn("get user profile failed",
			zap.Error(err),
			zap.Uint64("user_id", req.UserID),
		)
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
			s.log.Debug("invalid birthday",
				zap.Uint64("user_id", req.UserID),
				zap.String("birthday", req.Birthday),
			)
			return fmt.Errorf("invalid birthday %q, expected yyyy-MM-dd", req.Birthday)
		}
		data.BirthDate = sql.NullTime{Time: t, Valid: true}
	}

	if err := s.UserDomain.UpdateProfile(ctx, data); err != nil {
		s.log.Error("update user profile failed",
			zap.Error(err),
			zap.Uint64("user_id", req.UserID),
		)
		return err
	}
	s.log.Info("user profile updated", zap.Uint64("user_id", req.UserID))
	return nil
}

// UpdateAvatar 更新用户头像，返回新的头像地址
func (s *UserService) UpdateAvatar(ctx context.Context, req dto.UpdateAvatarRequest) (string, error) {
	avatar := strings.TrimSpace(req.Avatar)
	if avatar == "" {
		s.log.Debug("empty avatar", zap.Uint64("user_id", req.UserID))
		return "", errors.New("avatar is required")
	}
	if err := s.UserDomain.UpdateAvatar(ctx, req.UserID, avatar); err != nil {
		s.log.Error("update user avatar failed",
			zap.Error(err),
			zap.Uint64("user_id", req.UserID),
		)
		return "", err
	}
	s.log.Info("user avatar updated", zap.Uint64("user_id", req.UserID))
	return avatar, nil
}

// UpdateUserStatus 更新用户封禁状态
func (s *UserService) UpdateUserStatus(ctx context.Context, req dto.UpdateUserStatusRequest) error {
	status := enum.ParseStatus(req.Status)
	if status != enum.StatusApproved && status != enum.StatusBlocked {
		s.log.Debug("invalid user status",
			zap.Uint64("user_id", req.UserID),
			zap.String("status", req.Status),
		)
		return fmt.Errorf("invalid status %q, expected approved or blocked", req.Status)
	}
	if err := s.UserDomain.UpdateStatus(ctx, req.UserID, status); err != nil {
		s.log.Error("update user status failed",
			zap.Error(err),
			zap.Uint64("user_id", req.UserID),
			zap.String("status", status.String()),
		)
		return err
	}
	s.log.Info("user status updated",
		zap.Uint64("user_id", req.UserID),
		zap.String("status", status.String()),
	)
	return nil
}
