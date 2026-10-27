package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	entschema "entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []entschema.Annotation {
	return []entschema.Annotation{
		entsql.Annotation{Charset: "utf8mb4", Collation: "utf8mb4_general_ci"},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		// id 用户名 密码 邮箱 手机号 头像 创建、更新、删除时间 性别 签名 生日 role 用户状态
		field.Uint64("id").Comment("用户编号"),
		field.Time("created_at").Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.Time("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.Time("birth_date").Optional().Nillable().Comment("用户生日"),
		field.String("username").Unique().Comment("用户账号"),
		field.String("password").Comment("用户密码"),
		field.String("email").Optional().Default("").Comment("用户邮箱"),
		field.String("avatar").Optional().Default("").Comment("用户头像"),
		field.String("phone").Optional().Default("").Comment("用户手机号"),
		field.Enum("sex").Values("male", "female", "unknown").Default("unknown").Comment("性别"),
		field.Enum("role").Values("admin", "user").Default("user").Comment("账号权限 默认用户"),
		field.String("signature").Optional().Default("").Comment("用户签名"),
		field.Enum("status").Values("approved", "blocked").Default("approved").Comment("用户状态 默认为未封禁"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// 一个用户对应多个文章
		edge.To("articles", Article.Type),
		// 一个用户对应多个评论
		edge.To("comments", Comment.Type),
		// 一个用户对应一个用户状态
		edge.To("statu", UserStatu.Type),
	}
}
