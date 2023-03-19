package db_models

import "github.com/google/uuid"

type Follow struct {
	ID        uuid.UUID `json:"id,omitempty"`
	ChannelID uuid.UUID `json:"channel_id,omitempty"`
	ChatID    uuid.UUID `json:"chat_id,omitempty"`
}
