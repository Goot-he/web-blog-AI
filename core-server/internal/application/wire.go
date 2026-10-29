package application

import (
	"core-server/internal/infras/clog"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	clog.NewLog,
	NewArticleService,
	NewUserService,
	NewCategoryService,
	NewCommentService,
	NewCountService,
	NewLikeService,
)
