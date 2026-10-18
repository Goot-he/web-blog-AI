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
	wire.Bind(new(domain.ArticleDomainRepo), new(*ArticleRepo)),
	wire.Bind(new(domain.UserDomainRepo), new(*UserRepo)),
	wire.Bind(new(domain.CategoryDomainRepo), new(*CategoryRepo)),
	wire.Bind(new(domain.CommentDomainRepo), new(*CommentRepo)),
)
