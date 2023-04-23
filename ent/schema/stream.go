package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type Stream struct {
	ent.Schema
}

func (Stream) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.UUID("channel_id", uuid.UUID{}),

		field.Other("titles", pq.StringArray{}).
			SchemaType(map[string]string{
				dialect.Postgres: "text[]",
				dialect.SQLite:   "JSON",
			}).
			Default(pq.StringArray{}).
			Optional(),
		//SchemaType(map[string]string{
		//	"postgres": "text[]",
		//	"sqlite":   "text[]",
		//}),

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
		edge.To("stream_categories", StreamCategory.Type).Annotations(entsql.Annotation{
			OnDelete: entsql.Cascade,
		}),
	}
}
