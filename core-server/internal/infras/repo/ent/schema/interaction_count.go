package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Interaction_count holds the schema definition for the Interaction_count entity.
type Interaction_count struct {
	ent.Schema
}

// Fields of the Interaction_count.
func (Interaction_count) Fields() []ent.Field {
	return []ent.Field{
		// id 创建、更新、删除时间 对象类型 对象id 互动类型 数量
		field.Uint64("id").Comment("编号"),
		field.Time("created_at").Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.Time("deleted_at").Optional().Nillable().Comment("删除时间"),

		field.Enum("object_type").Values("article", "comment", "life").Comment("对象类型"),
		field.Uint64("object_id").Comment("对象ID"),
		field.Enum("interaction_type").Values("like", "comment", "favor", "unknown").Optional().Default("unknown").Comment("互动类型"),
		field.Uint64("count").Default(0).Comment("计数"),
	}
}

// Edges of the Interaction_count.
func (Interaction_count) Edges() []ent.Edge {
	return nil
} // 固定的静态关系才需要创建边 像这里的一对多或者多对多的关系不需要创建边

func (Interaction_count) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("object_type", "object_id", "interaction_type").Unique().StorageKey("interaction_object_type_id"), // 创建联合唯一索引
	}
}
