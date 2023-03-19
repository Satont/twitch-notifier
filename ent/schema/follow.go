package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Follow struct {
	ent.Schema
}

func (Follow) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("channel_id", uuid.UUID{}),
		field.UUID("chat_id", uuid.UUID{}),
	}
}

func (Follow) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("channel", Channel.Type).
			Required().
			Ref("follows").
			Unique().
			Field("channel_id"),
		edge.From("chat", Chat.Type).
			Required().
			Ref("follows").
			Unique().
			Field("chat_id"),
	}
}

func (Follow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("channel_id", "chat_id").
			Unique(),
	}
}
