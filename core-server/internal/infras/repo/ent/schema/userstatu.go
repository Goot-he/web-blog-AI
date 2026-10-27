package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	entschema "entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserStatu holds the schema definition for the UserStatu entity.
type UserStatu struct {
	ent.Schema
}

// Annotations of the UserStatu.
func (UserStatu) Annotations() []entschema.Annotation {
	return []entschema.Annotation{
		entsql.Annotation{Charset: "utf8mb4", Collation: "utf8mb4_general_ci"},
	}
}

// Fields of the UserStatu.
func (UserStatu) Fields() []ent.Field {
	return []ent.Field{
		//用户ID 文章数量 点赞量 收到点赞量 收藏量 收到收藏量 关注量 收到关注量 文章被浏览量
		field.Uint64("user_id").Comment("用户ID"),
		field.Time("created_at").Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.Time("deleted_at").Optional().Nillable().Comment("删除时间 假删除"),

		field.Uint64("article_count").Default(0).Comment("文章数量"),
		field.Uint64("like_count").Default(0).Comment("点赞量"),
		field.Uint64("receive_like_count").Default(0).Comment("收到点赞量"),
		field.Uint64("favor_count").Default(0).Comment("收藏量"),
		field.Uint64("receive_favor_count").Default(0).Comment("收到收藏量"),
		field.Uint64("follow_count").Default(0).Comment("关注量"),
		field.Uint64("receive_follow_count").Default(0).Comment("粉丝量"),
		field.Uint64("view_count").Default(0).Comment("被浏览量"),
		field.Uint64("comment_count").Default(0).Comment("评论数量"),
	}
}

// Edges of the UserStatu.
func (UserStatu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("statu").Field("user_id").Unique().Required(),
	}
}

// Indexes of the UserStatu.
func (UserStatu) Indexes() []ent.Index {
	return []ent.Index{
		// 每个用户只能有一条统计记录
		index.Fields("user_id").Unique().StorageKey("uk_user_status_user"),
	}
}
