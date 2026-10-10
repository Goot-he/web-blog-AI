package aggregate

import "core-server/internal/model/entity"

type ArticleAggregate struct {
	Article  *entity.Article
	Category *entity.Category
	Author   *entity.User
}
