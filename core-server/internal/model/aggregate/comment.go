package aggregate

import "core-server/internal/model/entity"

type CommentAggregate struct {
	Comment *entity.Comment
	User    *entity.User
}
