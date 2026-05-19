package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type PostView struct {
	ent.Schema
}

func (PostView) Fields() []ent.Field {
	return []ent.Field{
		field.Int("post_id"),
		field.String("user").Optional(),
		field.Time("created_at"),
	}
}

func (PostView) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("post", Post.Type).Ref("views").Field("post_id").Required().Unique(),
	}
}
