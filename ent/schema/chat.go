package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Chat struct {
	ent.Schema
}

type ChatService string

func (c ChatService) String() string {
	return string(c)
}

func (ChatService) Values() []string {
	return []string{Telegram.String()}
}

const (
	Telegram ChatService = "telegram"
)

func (Chat) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("chat_id"),
		field.Enum("service").Values(Telegram.String()),
	}
}

func (Chat) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("chat_id", "service").
			Unique(),
	}
}

func (Chat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("settings", ChatSettings.Type).Unique(),
		//edge.To("id", ChatSettings.Type).Unique().Required(),
		edge.To("follows", Follow.Type),
	}
}
