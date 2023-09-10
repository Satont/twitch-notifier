package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type ChatSettings struct {
	ent.Schema
}

type ChatLanguage string

const (
	ChatLanguageRu ChatLanguage = "ru"
	ChatLanguageEn ChatLanguage = "en"
	ChatLanguageUk ChatLanguage = "uk"
)

func (c ChatLanguage) String() string {
	return string(c)
}

func (ChatSettings) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Bool("game_change_notification").Default(true),
		field.Bool("title_change_notification").Default(false),
		field.Bool("game_and_title_change_notification").Default(false),
		field.Bool("offline_notification").Default(true),
		field.Bool("image_in_notification").Default(true),
		field.Enum("chat_language").
			Values(ChatLanguageRu.String(), ChatLanguageEn.String(), ChatLanguageUk.String()).
			Default(ChatLanguageEn.String()),
		field.String("test").Default("test"),
		field.UUID("chat_id", uuid.UUID{}),
	}
}

func (ChatSettings) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat", Chat.Type).
			Ref("settings").
			Unique().
			Field("chat_id").
			Required(),
	}
}
