package repo

import (
	"core-server/internal/domain"

	"github.com/google/wire"
)

var SetRepo = wire.NewSet(
	NewEntClient,
	NewArticleRepo,
	NewUserRepo,
	NewCategoryRepo,
	NewCommentRepo,
	NewCountRepo,
	NewLikeRepo,
	wire.Bind(new(domain.ArticleDomainRepo), new(*ArticleRepo)),
	wire.Bind(new(domain.UserDomainRepo), new(*UserRepo)),
	wire.Bind(new(domain.CategoryRepoDomain), new(*CategoryRepo)),
	wire.Bind(new(domain.CommentRepoDomain), new(*CommentRepo)),
	wire.Bind(new(domain.CountRepoDomain), new(*CountRepo)),
	wire.Bind(new(domain.LikeDomainRepo), new(*LikeRepo)),
)
