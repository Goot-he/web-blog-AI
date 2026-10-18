package application

import (
	"github.com/google/wire"
)

var SetService = wire.NewSet(
	NewArticleService,
	NewUserService,
	NewCategoryService,
	NewCommentService,
)
