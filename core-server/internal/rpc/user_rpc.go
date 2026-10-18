package rpc

import (
	"context"

	"core-server/internal/application"
	"core-server/internal/model/dto"
	"core-server/internal/rpc/userpb"
)

type UserRpc struct {
	userpb.UnimplementedUserServiceServer
	userService *application.UserService
}

func NewUserRpc(userService *application.UserService) *UserRpc {
	return &UserRpc{
		userService: userService,
	}
}

// pageValue 把 proto 里的 int32 页码安全地转成 uint32
func pageValue(v int32) uint32 {
	if v <= 0 {
		return 0
	}
	return uint32(v)
}

// ListUsers 分页查询用户列表
func (r *UserRpc) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	resp, err := r.userService.ListUsers(ctx, dto.ListUsersRequest{
		Page:     pageValue(req.GetPage()),
		PageSize: pageValue(req.GetPageSize()),
	})
	if err != nil {
		return nil, err
	}
	return &userpb.ListUsersResponse{
		Users: toPbUserInfos(resp.Users),
		Total: resp.Total,
	}, nil
}

// SearchUsers 按关键词搜索用户
func (r *UserRpc) SearchUsers(ctx context.Context, req *userpb.SearchUsersRequest) (*userpb.SearchUsersResponse, error) {
	resp, err := r.userService.SearchUsers(ctx, dto.SearchUsersRequest{
		Keyword:  req.GetKeyword(),
		Page:     pageValue(req.GetPage()),
		PageSize: pageValue(req.GetPageSize()),
	})
	if err != nil {
		return nil, err
	}
	return &userpb.SearchUsersResponse{
		Users: toPbUserInfos(resp.Users),
		Total: resp.Total,
	}, nil
}

// GetProfile 查询用户个人主页
func (r *UserRpc) GetProfile(ctx context.Context, req *userpb.GetProfileRequest) (*userpb.GetProfileResponse, error) {
	resp, err := r.userService.GetProfile(ctx, dto.GetProfileRequest{
		UserID: req.GetUserId(),
	})
	if err != nil {
		return nil, err
	}
	return &userpb.GetProfileResponse{
		Id:                resp.ID,
		Username:          resp.Username,
		Email:             resp.Email,
		Phone:             resp.Phone,
		Avatar:            resp.Avatar,
		Sex:               resp.Sex,
		Birthday:          resp.Birthday,
		Role:              resp.Role,
		Status:            resp.Status,
		Signature:         resp.Signature,
		FollowersCount:    resp.FollowersCount,
		FollowingCount:    resp.FollowingCount,
		LikeCount:         resp.LikeCount,
		ReceiveLikeCount:  resp.ReceiveLikeCount,
		FavorCount:        resp.FavorCount,
		ReceiveFavorCount: resp.ReceiveFavorCount,
		ArticleCount:      resp.ArticleCount,
		ViewCount:         resp.ViewCount,
		CommentCount:      resp.CommentCount,
	}, nil
}

// UpdateProfile 更新用户资料
func (r *UserRpc) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*userpb.UpdateProfileResponse, error) {
	err := r.userService.UpdateProfile(ctx, dto.UpdateProfileRequest{
		UserID:    req.GetUserId(),
		Username:  req.GetUsername(),
		Phone:     req.GetPhone(),
		Sex:       req.GetSex(),
		Birthday:  req.GetBirthday(),
		Signature: req.GetSignature(),
	})
	if err != nil {
		return nil, err
	}
	return &userpb.UpdateProfileResponse{Success: true}, nil
}

// UpdateAvatar 更新用户头像
func (r *UserRpc) UpdateAvatar(ctx context.Context, req *userpb.UpdateAvatarRequest) (*userpb.UpdateAvatarResponse, error) {
	url, err := r.userService.UpdateAvatar(ctx, dto.UpdateAvatarRequest{
		UserID: req.GetUserId(),
		Avatar: req.GetAvatar(),
	})
	if err != nil {
		return nil, err
	}
	return &userpb.UpdateAvatarResponse{Url: url}, nil
}

// UpdateUserStatus 更新用户封禁状态
func (r *UserRpc) UpdateUserStatus(ctx context.Context, req *userpb.UpdateUserStatusRequest) (*userpb.UpdateUserStatusResponse, error) {
	err := r.userService.UpdateUserStatus(ctx, dto.UpdateUserStatusRequest{
		UserID: req.GetUserId(),
		Status: req.GetStatus(),
	})
	if err != nil {
		return nil, err
	}
	return &userpb.UpdateUserStatusResponse{Success: true}, nil
}

// toPbUserInfo dto.UserInfo 转 proto
func toPbUserInfo(u *dto.UserInfo) *userpb.UserInfo {
	if u == nil {
		return nil
	}
	return &userpb.UserInfo{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Sex:       u.Sex,
		Birthday:  u.Birthday,
		Signature: u.Signature,
		Role:      u.Role,
		Status:    u.Status,
	}
}

// toPbUserInfos 批量转换
func toPbUserInfos(users []*dto.UserInfo) []*userpb.UserInfo {
	if len(users) == 0 {
		return nil
	}
	out := make([]*userpb.UserInfo, 0, len(users))
	for _, u := range users {
		out = append(out, toPbUserInfo(u))
	}
	return out
}
