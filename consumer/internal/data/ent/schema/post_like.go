package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type PostLike struct {
	ent.Schema
}

func (PostLike) Fields() []ent.Field {
	return []ent.Field{
		field.Int("post_id"),
		field.String("user").NotEmpty(),
		field.Time("created_at"),
	}
}

func (PostLike) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("post", Post.Type).Ref("likes").Field("post_id").Required().Unique(),
	}
}

func (PostLike) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("post_id", "user").Unique(),
	}
}
