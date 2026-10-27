package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	entschema "entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Comment holds the schema definition for the Comment entity.
type Comment struct {
	ent.Schema
}

// Annotations of the Comment.
func (Comment) Annotations() []entschema.Annotation {
	return []entschema.Annotation{
		entsql.Annotation{Charset: "utf8mb4", Collation: "utf8mb4_general_ci"},
	}
}

// Fields of the Comment.
func (Comment) Fields() []ent.Field {
	return []ent.Field{
		// 评论所属文章ID 评论ID 根评论id 回复评论ID 评论内容 发布、删除时间 发布用户ID 回复用户ID 子评论数量 点赞量
		field.Uint64("id").Comment("评论编号"),
		field.Time("created_at").Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.Time("deleted_at").Optional().Nillable().Comment("删除时间"),

		field.Uint64("article_id").Comment("评论所属文章ID"),
		field.Uint64("publish_user_id").Comment("发布评论的用户ID"),
		field.Uint64("root_id").Optional().Default(0).Comment("根评论ID"),
		field.Uint64("reply_id").Optional().Default(0).Comment("回复的评论ID"),
		field.Uint64("reply_user_id").Optional().Default(0).Comment("回复的用户的ID"),
		field.String("content").
			SchemaType(map[string]string{dialect.MySQL: "text"}).
			Comment("评论内容"),
		field.Uint64("like_count").Default(0).Comment("评论点赞量"),
		field.Uint64("child_count").Default(0).Comment("评论下的子评论数量"),
	}
}

// Edges of the Comment.
func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		// 多个评论来自一篇文章
		edge.From("article", Article.Type).Ref("comments").Field("article_id").Unique().Required(),
		// 多个评论来自一个用户
		edge.From("user", User.Type).Ref("comments").Field("publish_user_id").Unique().Required(),
	}
}
