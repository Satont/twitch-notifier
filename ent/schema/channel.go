package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	ent.Schema
}

type ChannelService string

func (c ChannelService) String() string {
	return string(c)
}

const (
	Twitch ChannelService = "twitch"
)

func (Channel) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("channel_id"),
		field.Enum("service").Values(Twitch.String()),
		field.Bool("is_live").Default(false),
		field.String("title").Nillable().Optional(),
		field.String("category").Nillable().Optional(),
		field.Time("updated_at").Nillable().Optional().Default(nil).UpdateDefault(time.Now().UTC),
	}
}

func (Channel) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("channel_id", "service").
			Unique(),
	}
}

func (Channel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("follows", Follow.Type),
		edge.To("streams", Stream.Type),
	}
}
