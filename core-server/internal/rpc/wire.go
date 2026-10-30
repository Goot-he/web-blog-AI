package rpc

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewArticleRpc,
	NewUserRpc,
	NewCategoryRpc,
	NewCommentRpc,
	NewAuthRPC,
	NewLikeRpc,

	NewServer,
)
