package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Post struct {
	ent.Schema
}

func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.Text("body").NotEmpty(),
		field.String("author").NotEmpty(),
		field.Time("created_at"),
	}
}

func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("comments", Comment.Type),
		edge.To("likes", PostLike.Type),
		edge.To("views", PostView.Type),
	}
}
