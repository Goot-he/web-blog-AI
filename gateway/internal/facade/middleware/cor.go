package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
// 允许所有来源的跨域请求 支持GET POST DELETE PUT OPTIONS方法
// 允许携带 Authorization 等自定义请求头
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源的访问
		c.Header("Access-Control-Allow-Origin", "*")
		// 允许的请求头
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// 如果是预检请求（options） 直接返回，不继续处理
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
