package db_models

import (
	"errors"
	"github.com/google/uuid"
)

var (
	FollowAlreadyExistsError = errors.New("follow already exists")
	FollowNotFoundError      = errors.New("follow not found")
)

type Follow struct {
	ID        uuid.UUID `json:"id,omitempty"`
	ChannelID uuid.UUID `json:"channel_id,omitempty"`
	ChatID    uuid.UUID `json:"chat_id,omitempty"`

	Channel *Channel `json:"channel,omitempty"`
	Chat    *Chat    `json:"chat,omitempty"`
}
