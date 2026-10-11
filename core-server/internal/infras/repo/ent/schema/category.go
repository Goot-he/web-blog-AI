package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		// 编号 分类名 上级分类ID 创建、更新、删除时间
		field.Uint64("id").Comment("用户编号"),
		field.Time("created_at").Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.Time("deleted_at").Optional().Nillable().Default(time.Now).Comment("删除时间"),
		field.String("category").Comment("分类名"),
		field.Uint64("type_id").Optional().Default(0).Comment("上级分类ID"),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		// 一个分类对应多个文章
		edge.To("articles", Article.Type),
	}
}
