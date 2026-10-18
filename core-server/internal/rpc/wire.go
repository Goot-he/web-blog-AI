package rpc

import "github.com/google/wire"

var SetRpc = wire.NewSet(
	NewArticleRpc,
	NewUserRpc,
	NewCategoryRpc,
	NewCommentRpc,
)
