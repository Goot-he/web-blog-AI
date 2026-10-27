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

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Annotations of the Article.
func (Article) Annotations() []entschema.Annotation {
	return []entschema.Annotation{
		entsql.Annotation{Charset: "utf8mb4", Collation: "utf8mb4_general_ci"},
	}
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		// 编号 创建、更新、删除、发布时间 标题 概要 正文 作者编号 分类编号 是否发布、置顶 评论、点赞、收藏、浏览数量 封面图片
		field.Uint64("id").Comment("文章编号"),
		field.Time("created_at").Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		// 软删除字段不能给默认值 否则新建的文章一插入就带上了删除时间
		field.Time("deleted_at").Optional().Nillable().Comment("删除时间 假删除"),
		// 同理 草稿的发布时间必须为nil 由业务在发布时写入 更新时也不能自动刷新
		field.Time("published_at").Optional().Nillable().Comment("发布时间 为nil表示未发布"),

		field.String("title").Comment("文章标题"),
		field.String("summary").Comment("文章概述"),
		field.String("content").
			SchemaType(map[string]string{dialect.MySQL: "longtext"}).
			Comment("文章正文"),
		field.String("coverage").Optional().Default("").Comment("文章封面"),
		field.Uint64("author_id").Comment("作者ID"),
		field.Uint64("category_id").Comment("分类ID"),

		field.Bool("is_top").Optional().Default(false).Comment("是否置顶 默认不置顶"),
		field.Bool("is_published").Optional().Default(true).Comment("是否发布 默认发布"),

		field.Uint64("comment_count").Optional().Default(0).Comment("评论数量"),
		field.Uint64("like_count").Optional().Default(0).Comment("点赞数量"),
		field.Uint64("favor_count").Optional().Default(0).Comment("收藏数量"),
		field.Uint64("view_count").Optional().Default(0).Comment("浏览数量"),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		// 一篇文章对应多个评论
		edge.To("comments", Comment.Type),

		// 多个文章属于一个分类
		edge.From("category", Category.Type).
			Ref("articles").      // Category里面的边叫做articles 这里表示的是反向边 反向边必须和正向边的名字一致
			Field("category_id"). // 外键字段
			Unique().             // 表示一篇文章只属于一个分类
			Required(),           // 表示创建文章时必须指定分类

		// 多篇文章属于一个用户
		edge.From("user", User.Type).
			Ref("articles").    // User里面的边叫作articles 这里表示的是反向边
			Field("author_id"). // 外键字段
			Unique().           // 表示一篇文章只属于一个用户
			Required(),         // 表示创建文章时必须指定作者
	}
}
