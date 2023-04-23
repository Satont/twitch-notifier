package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

type StreamCategory struct {
	ent.Schema
}

func (StreamCategory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("stream_id"),
		field.String("name"),

		field.Time("setted_at").Optional().Default(time.Now().UTC),
	}
}

func (StreamCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stream", Stream.Type).
			Ref("stream_categories").
			Required().
			Unique().
			Field("stream_id"),
	}
}
