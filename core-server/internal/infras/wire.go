package infras

import (
	"core-server/internal/domain"
	"core-server/internal/infras/cache"
	"core-server/internal/infras/clog"
	"core-server/internal/infras/repo"

	"github.com/google/wire"
)

var RepoProviderSet = wire.NewSet(
	repo.NewEntClient,
	repo.NewArticleRepo,
	repo.NewUserRepo,
	repo.NewCategoryRepo,
	repo.NewCommentRepo,
	repo.NewCountRepo,
	repo.NewLikeRepo,
	wire.Bind(new(domain.ArticleDomainRepo), new(*repo.ArticleRepo)),
	wire.Bind(new(domain.UserDomainRepo), new(*repo.UserRepo)),
	wire.Bind(new(domain.CategoryRepoDomain), new(*repo.CategoryRepo)),
	wire.Bind(new(domain.CommentRepoDomain), new(*repo.CommentRepo)),
	wire.Bind(new(domain.CountRepoDomain), new(*repo.CountRepo)),
	wire.Bind(new(domain.LikeDomainRepo), new(*repo.LikeRepo)),
)

var CacheProviderSet = wire.NewSet(
	cache.NewCacheClient,
	cache.NewLikeCache,

	wire.Bind(new(domain.LikeCacheDomain), new(*cache.LikeCache)),
)

var LogProviderSet = wire.NewSet(
	clog.NewLog,
)

var JobProviderSet = wire.NewSet(
	RepoProviderSet,
	CacheProviderSet,
	LogProviderSet,
)
