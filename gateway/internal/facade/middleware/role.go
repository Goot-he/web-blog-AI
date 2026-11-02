package middleware

import (
	"gateway/internal/model/response"

	"github.com/gin-gonic/gin"
)

func Role() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetRole(c)
		if role != "admin" {
			// 只有管理员才有权限操作
			response.Forbidden(c)
			return
		}
	}
}
