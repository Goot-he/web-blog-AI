package repo

import (
	"database/sql"
	"time"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/user"
	"core-server/internal/model/aggregate"
	"core-server/internal/model/entity"
	"core-server/internal/model/enum"
)

const (
	defaultPageSize uint32 = 10  // 不传分页大小时使用
	maxPageSize     uint32 = 100 // 单页最多返回的数量
)

// pageLimitOffset 把从1开始的页码换算成数据库的 limit/offset，并兜底非法入参
func pageLimitOffset(page, size uint32) (limit, offset int) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return int(size), int(uint64(page-1) * uint64(size))
}

// limitOffset 兜底 domain 层传入的 limit/offset（limit 默认 10、上限 100，offset 不能为负）
func limitOffset(limit, offset int32) (int, int) {
	if limit <= 0 {
		limit = int32(defaultPageSize)
	}
	if limit > int32(maxPageSize) {
		limit = int32(maxPageSize)
	}
	if offset < 0 {
		offset = 0
	}
	return int(limit), int(offset)
}

// normalizeOffsetLimit 兜底 offset/limit（offset 不能为负，limit 默认 10、上限 100）
func normalizeOffsetLimit(offset, limit int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = int(defaultPageSize)
	}
	if limit > int(maxPageSize) {
		limit = int(maxPageSize)
	}
	return offset, limit
}

// nullTime 把 ent 的可空时间（*time.Time）转成 entity 的 sql.NullTime
func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// toEntityArticle 把 ent.Article 转成 entity.Article
func toEntityArticle(a *ent.Article) *entity.Article {
	if a == nil {
		return nil
	}
	return &entity.Article{
		ID:           a.ID,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
		DeletedAt:    nullTime(a.DeletedAt),
		PublishedAt:  nullTime(a.PublishedAt),
		Title:        a.Title,
		Summary:      a.Summary,
		Content:      a.Content,
		CoverImage:   a.Coverage,
		AuthorID:     a.AuthorID,
		CategoryID:   a.CategoryID,
		IsTop:        a.IsTop,
		IsPublished:  a.IsPublished,
		ViewCount:    a.ViewCount,
		LikeCount:    a.LikeCount,
		FavorCount:   a.FavorCount,
		CommentCount: a.CommentCount,
	}
}

// toEntityUser 把 ent.User 转成 entity.User
func toEntityUser(u *ent.User) *entity.User {
	if u == nil {
		return nil
	}
	return &entity.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: nullTime(u.DeletedAt),
		BirthDate: nullTime(u.BirthDate),
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Sex:       enum.ParseSex(string(u.Sex)),
		Role:      toEntityRole(u.Role),
		Signature: u.Signature,
		Status:    enum.ParseStatus(string(u.Status)),
	}
}

// toEntityUsers 批量把 ent.User 转成 entity.User
func toEntityUsers(users []*ent.User) []*entity.User {
	if len(users) == 0 {
		return nil
	}
	out := make([]*entity.User, 0, len(users))
	for _, u := range users {
		out = append(out, toEntityUser(u))
	}
	return out
}

// toEntityRole 把数据库的字符串角色枚举映射成领域模型的数字枚举
func toEntityRole(r user.Role) enum.UserRole {
	switch r {
	case user.RoleAdmin:
		return enum.UserRoleAdmin
	case user.RoleUser:
		return enum.UserRoleUser
	default:
		return enum.UserRoleUnknown
	}
}

// toEntityCategory 把 ent.Category 转成 entity.Category
func toEntityCategory(c *ent.Category) *entity.Category {
	if c == nil {
		return nil
	}
	return &entity.Category{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: nullTime(c.DeletedAt),
		ParentID:  c.TypeID,
		Category:  c.Category,
	}
}

// toEntityCategories 批量把 ent.Category 转成 entity.Category
func toEntityCategories(categories []*ent.Category) []*entity.Category {
	if len(categories) == 0 {
		return nil
	}
	out := make([]*entity.Category, 0, len(categories))
	for _, c := range categories {
		out = append(out, toEntityCategory(c))
	}
	return out
}

// toArticleAggregate 把 ent.Article（含预加载的作者和分类）转成文章聚合
func toArticleAggregate(a *ent.Article) *aggregate.ArticleAggregate {
	if a == nil {
		return nil
	}
	agg := &aggregate.ArticleAggregate{
		Article: toEntityArticle(a),
	}
	if a.Edges.User != nil {
		agg.Author = toEntityUser(a.Edges.User)
	}
	if a.Edges.Category != nil {
		agg.Category = toEntityCategory(a.Edges.Category)
	}
	return agg
}

// toArticleAggregates 批量把 ent.Article 列表转成文章聚合列表
func toArticleAggregates(articles []*ent.Article) []*aggregate.ArticleAggregate {
	if len(articles) == 0 {
		return nil
	}
	aggs := make([]*aggregate.ArticleAggregate, 0, len(articles))
	for _, a := range articles {
		aggs = append(aggs, toArticleAggregate(a))
	}
	return aggs
}

// toEntityUserStatu 把 ent.UserStatu 转成 entity.UserStatus
func toEntityUserStatu(s *ent.UserStatu) *entity.UserStatus {
	if s == nil {
		return nil
	}
	return &entity.UserStatus{
		ID:                 uint64(s.ID), // user_statu 的主键是 int
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
		DeletedAt:          nullTime(s.DeletedAt),
		UserID:             s.UserID,
		ArticleCount:       s.ArticleCount,
		LikeCount:          s.LikeCount,
		ReceiveLikeCount:   s.ReceiveLikeCount,
		FavorCount:         s.FavorCount,
		ReceiveFavorCount:  s.ReceiveFavorCount,
		FollowCount:        s.FollowCount,
		ReceiveFollowCount: s.ReceiveFollowCount,
		ViewCount:          s.ViewCount,
		CommentCount:       s.CommentCount,
	}
}

// toEntityComment 把 ent.Comment 转成 entity.Comment
func toEntityComment(c *ent.Comment) *entity.Comment {
	if c == nil {
		return nil
	}
	return &entity.Comment{
		ID:            c.ID,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		DeletedAt:     nullTime(c.DeletedAt),
		ArticleID:     c.ArticleID,
		RootID:        c.RootID,
		ReplyToID:     c.ReplyID,
		PublishUserID: c.PublishUserID,
		ReplyToUserID: c.ReplyUserID,
		Comment:       c.Content,
		LikeCount:     c.LikeCount,
		ChildCount:    c.ChildCount,
	}
}

// toCommentAggregate 把 ent.Comment（含预加载的用户）转成评论聚合
func toCommentAggregate(c *ent.Comment) *aggregate.CommentAggregate {
	if c == nil {
		return nil
	}
	agg := &aggregate.CommentAggregate{
		Comment: toEntityComment(c),
	}
	if c.Edges.User != nil {
		agg.User = toEntityUser(c.Edges.User)
	}
	return agg
}

// toCommentAggregates 批量把 ent.Comment 列表转成评论聚合列表
func toCommentAggregates(comments []*ent.Comment) []*aggregate.CommentAggregate {
	if len(comments) == 0 {
		return nil
	}
	aggs := make([]*aggregate.CommentAggregate, 0, len(comments))
	for _, c := range comments {
		aggs = append(aggs, toCommentAggregate(c))
	}
	return aggs
}

// toEntSex 领域模型的性别枚举转成 ent 的枚举
func toEntSex(s enum.Sex) user.Sex {
	return user.Sex(string(s))
}

// toEntStatus 领域模型的用户状态枚举转成 ent 的枚举
func toEntStatus(s enum.UserStatus) user.Status {
	return user.Status(string(s))
}

// toEntRole 数字角色枚举转成 ent 的字符串枚举（数据库只支持 admin/user）
func toEntRole(r enum.UserRole) user.Role {
	if r == enum.UserRoleAdmin {
		return user.RoleAdmin
	}
	return user.RoleUser
}

func toEntityLike(node *ent.InteractionLike) *entity.InteractionLike {
	if node == nil {
		return nil
	}
	return &entity.InteractionLike{
		ID:         node.ID,
		ObjectType: enum.ParseObjectType(node.ObjectType.String()),
		ObjectID:   node.ObjectID,
		LikeStatus: entity.ParseInteractionLikeType(node.Status.String()),
		UserID:     node.UserID,
		Version:    node.Version,
	}
}

func toEntityLikes(likes []*ent.InteractionLike) []*entity.InteractionLike {
	items := make([]*entity.InteractionLike, 0, len(likes))
	for _, l := range likes {
		items = append(items, toEntityLike(l))
	}
	return items
}

func toEntityInteractionCount(interactionCount *ent.InteractionCount) *entity.InteractionCount {
	if interactionCount == nil {
		return nil
	}
	return &entity.InteractionCount{
		ID:              interactionCount.ID,
		ObjectID:        interactionCount.ObjectID,
		ObjectType:      enum.ParseObjectType(interactionCount.ObjectType.String()),
		InteractionType: enum.ParseInteractionType(interactionCount.InteractionType.String()),
		Count:           interactionCount.Count,
	}
}

func toEntityInteractionCounts(interactionCount []*ent.InteractionCount) []*entity.InteractionCount {
	items := make([]*entity.InteractionCount, 0, len(interactionCount))
	for _, i := range interactionCount {
		items = append(items, toEntityInteractionCount(i))
	}
	return items
}
