package aggregate

import "core-server/internal/model/entity"

type UserAggregate struct {
	User  *entity.User
	Statu *entity.UserStatus
}
