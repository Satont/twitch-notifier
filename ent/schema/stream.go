package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

type Stream struct {
	ent.Schema
}

func (Stream) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.UUID("channel_id", uuid.UUID{}),

		field.Strings("titles").
			Optional().
			Default([]string{}),
		//SchemaType(map[string]string{
		//	"postgres": "text[]",
		//	"sqlite":   "text[]",
		//}),
		field.Strings("categories").
			Optional().
			Default([]string{}),
		//SchemaType(map[string]string{
		//	"postgres": "text[]",
		//	"sqlite":   "text[]",
		//}),

		field.Time("started_at").Optional().Default(time.Now().UTC),
		field.Time("updated_at").Nillable().Optional().Default(nil).UpdateDefault(time.Now().UTC),
		field.Time("ended_at").Nillable().Optional().Default(nil),
	}
}

func (Stream) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("channel", Channel.Type).
			Ref("streams").
			Required().
			Unique().
			Field("channel_id"),
	}
}
