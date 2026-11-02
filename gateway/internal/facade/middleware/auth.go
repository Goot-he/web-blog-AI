package middleware

const (
	AuthUserIDContextKey       = "user_id"
	AuthRoleContextKey         = "role"
	refreshedAccessTokenHeader = "Authorization"
)

type Middleware struct {
}
