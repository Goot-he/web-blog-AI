package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (uint64, bool) {
	value, exits := c.Get(AuthUserIDContextKey)
	if !exits {
		return 0, false
	}
	id, exits := value.(uint64)
	return id, exits && id > 0
}

func GetRole(c *gin.Context) string {
	role := c.GetString(AuthRoleContextKey)
	return role
}
