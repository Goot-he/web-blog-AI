package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Interaction_like holds the schema definition for the Interaction_like entity.
type Interaction_like struct {
	ent.Schema
}

// Fields of the Interaction_like.
func (Interaction_like) Fields() []ent.Field {
	return []ent.Field{
		// id 创建、更新、删除时间 点赞对象类型 点赞对象id 用户ID 点赞状态
		field.Uint64("id").Comment("编号"),
		field.Time("created_at").Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.Time("deleted_at").Optional().Nillable().Default(time.Now).Comment("删除时间"),

		field.Enum("object_type").Values("article", "comment", "unknown").Default("unknown").Comment("点赞对象"),
		field.Uint64("object_id").Comment("点赞对象ID"),
		field.Uint64("user_id").Comment("点赞用户ID"),
		field.Enum("status").Values("unknow", "thump_up", "nothing").Default("nothing").Comment("点赞状态"),
	}
}

// Edges of the Interaction_like.
func (Interaction_like) Edges() []ent.Edge {
	return nil
} // 固定的静态关系才需要创建边

func (Interaction_like) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("object_type", "object_id", "user_id").Unique().StorageKey("user_object_type_id"), // 创建联合唯一索引
	}
}
