package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{
		"code": status,
		"msg":  msg,
	})
}

// Unauthorized 未授权或者无权限
func Unauthorized(c *gin.Context) {
	Fail(c, http.StatusUnauthorized, "unauthorized")
}

// Forbidden 无权限访问
func Forbidden(c *gin.Context) {
	Fail(c, http.StatusForbidden, "Access denied for this role")
}

// StatusBadRequest 参数解析错误
func StatusBadRequest(c *gin.Context) {
	Fail(c, http.StatusBadRequest, "invalid request parameters")
}

// InternalServerError 服务器内部错误
func InternalServerError(c *gin.Context, msg string) {
	Fail(c, http.StatusInternalServerError, msg)
}
