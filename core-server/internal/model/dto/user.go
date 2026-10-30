package dto

import (
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

type LoginRequest struct {
	Username string
	Password string
}
type EmailLoginRequest struct {
	Email string
}
type LoginResponse struct {
	Id     uint64
	Name   string
	Avatar string
	Status enum.UserStatus
	Role   enum.UserRole
}

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}
type ForgotPasswordRequest struct {
	Email    string
	Password string
	Confirm  string
}

type UserInfo struct {
	ID        uint64
	Username  string
	Email     string
	Phone     string
	Avatar    string
	Sex       string
	Birthday  string
	Signature string
	Role      string
	Status    string
}

type ListUsersRequest struct {
	Page     uint32
	PageSize uint32
}
type ListUsersResponse struct {
	Users []*UserInfo
	Total uint64
}

type SearchUsersRequest struct {
	Keyword  string
	Page     uint32
	PageSize uint32
}
type SearchUsersResponse struct {
	Users []*UserInfo
	Total uint64
}

type GetProfileRequest struct {
	UserID uint64
}
type GetProfileResponse struct {
	ID        uint64
	Username  string
	Email     string
	Phone     string
	Avatar    string
	Sex       string
	Birthday  string
	Role      string
	Status    string
	Signature string

	FollowersCount    uint64
	FollowingCount    uint64
	LikeCount         uint64
	ReceiveLikeCount  uint64
	FavorCount        uint64
	ReceiveFavorCount uint64
	ArticleCount      uint64
	ViewCount         uint64
	CommentCount      uint64
}

type UpdateProfileRequest struct {
	UserID    uint64
	Username  string
	Phone     string
	Sex       string
	Birthday  string
	Signature string
}

type UpdateAvatarRequest struct {
	UserID uint64
	Avatar string
}

type UpdateUserStatusRequest struct {
	UserID uint64
	Status string
}

// ToLoginResponse 转登录信息
func ToLoginResponse(user *entity.User) *LoginResponse {
	return &LoginResponse{
		Id:     user.ID,
		Name:   user.Username,
		Avatar: user.Avatar,
		Status: user.Status,
		Role:   user.Role,
	}
}

// ToUserInfo entity.User 转用户信息
func ToUserInfo(u *entity.User) *UserInfo {
	if u == nil {
		return nil
	}
	return &UserInfo{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Sex:       string(u.Sex),
		Birthday:  birthdayString(u),
		Signature: u.Signature,
		Role:      roleName(u.Role),
		Status:    string(u.Status),
	}
}

// ToUserInfos 批量转换用户信息
func ToUserInfos(users []*entity.User) []*UserInfo {
	if len(users) == 0 {
		return nil
	}
	out := make([]*UserInfo, 0, len(users))
	for _, u := range users {
		out = append(out, ToUserInfo(u))
	}
	return out
}

// ToGetProfileResponse 用户 + 统计信息 转个人主页响应
func ToGetProfileResponse(u *entity.User, statu *entity.UserStatus, commentCount uint64) *GetProfileResponse {
	if u == nil {
		return nil
	}
	resp := &GetProfileResponse{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		Phone:        u.Phone,
		Avatar:       u.Avatar,
		Sex:          string(u.Sex),
		Birthday:     birthdayString(u),
		Role:         roleName(u.Role),
		Status:       string(u.Status),
		Signature:    u.Signature,
		CommentCount: commentCount,
	}
	if statu != nil {
		resp.FollowersCount = statu.ReceiveFollowCount
		resp.FollowingCount = statu.FollowCount
		resp.LikeCount = statu.LikeCount
		resp.ReceiveLikeCount = statu.ReceiveLikeCount
		resp.FavorCount = statu.FavorCount
		resp.ReceiveFavorCount = statu.ReceiveFavorCount
		resp.ArticleCount = statu.ArticleCount
		resp.ViewCount = statu.ViewCount
	}
	return resp
}

// birthdayString 把可空生日格式化成 yyyy-MM-dd
func birthdayString(u *entity.User) string {
	if !u.BirthDate.Valid {
		return ""
	}
	return u.BirthDate.Time.Format("2006-01-02")
}

// roleName 角色枚举转接口使用的字符串
func roleName(r enum.UserRole) string {
	switch r {
	case enum.UserRoleAdmin:
		return "admin"
	case enum.UserRoleUser:
		return "user"
	default:
		return "unknown"
	}
}
