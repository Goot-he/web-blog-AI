package middleware

import (
	"context"
	"gateway/internal/config"
	"gateway/internal/infras/cache"
	"gateway/internal/model/response"
	"gateway/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthUserIDContextKey       = "user_id"
	AuthRoleContextKey         = "role"
	refreshedAccessTokenHeader = "Authorization"
)

type AuthMiddleware struct {
	cfg       *config.Config
	blacklist *cache.Blacklist
}

func NewAuthMiddleware(cfg *config.Config, blacklist *cache.Blacklist) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:       cfg,
		blacklist: blacklist,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 校验access token是否在黑名单
		// 如果没有access token 获取refresh token
		// 校验refresh token是否在黑名单
		// 如果不在根据它生成新的access token
		// 设置新的access token到cache
		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			response.Fail(c, http.StatusUnauthorized, "invalid or expired access-token")
			return
		}

		isBlack, err := m.blacklist.IsTokenBlacklist(context.Background(), accessToken)
		if err != nil {
			c.Abort()
			response.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if isBlack {
			c.Abort()
			response.Fail(c, http.StatusUnauthorized, "access-token in blacklist")
		}

		claims, err := utils.GetClaims(accessToken, []byte(m.cfg.Auth.JwtSecret))
		if err != nil {
			// 说明claims有问题 需要重新申请
			refreshToken, err := utils.RefreshTokenFromCookie(c.Request)
			if err != nil {
				c.Abort()
				response.Fail(c, http.StatusUnauthorized, "invalid or expired refresh-token")
				return
			}

			isBlack, err = m.blacklist.IsTokenBlacklist(c.Request.Context(), refreshToken)
			if err != nil {
				c.Abort()
				response.Fail(c, http.StatusUnauthorized, "internal server error")
				return
			}
			if isBlack {
				c.Abort()
				response.Fail(c, http.StatusUnauthorized, "refresh-token in blacklist")
				return
			}

			refreshClaims, err := utils.GetClaims(refreshToken, []byte(m.cfg.Auth.JwtSecret))
			if err != nil {
				c.Abort()
				response.Fail(c, http.StatusUnauthorized, "invalid or expired refresh-token")
				return
			}

			newAccessToken, err := utils.CreateAccessToken([]byte(m.cfg.Auth.JwtSecret), *refreshClaims, m.cfg.Auth.AccessExpire)
			if err != nil {
				c.Abort()
				response.Fail(c, http.StatusUnauthorized, "internal server error")
				return
			}

			c.Header(refreshedAccessTokenHeader, "Bearer "+newAccessToken)
			c.Header("Access-Control-Expose-Headers", refreshedAccessTokenHeader)
			claims = refreshClaims
		}

		UserId, role := claims.UserID, claims.Role
		// 校验用户ID是否在黑名单
		isBlack, err = m.blacklist.IsUserBlacklisted(context.Background(), UserId)
		if err != nil {
			c.Abort()
			response.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if isBlack {
			c.Abort()
			response.Fail(c, http.StatusUnauthorized, "user in blacklist")
			return
		}
		c.Set(AuthUserIDContextKey, UserId)
		c.Set(AuthRoleContextKey, role)
		c.Next()
	}
}

func bearerToken(value string) (string, bool) {
	parts := strings.Fields(value)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
