package db_models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ChannelNotFoundError = errors.New("channel not found")
)

type ChannelService string

const (
	ChannelServiceTwitch ChannelService = "twitch"
)

func (s ChannelService) String() string {
	return string(s)
}

type Channel struct {
	ID        uuid.UUID      `json:"id,omitempty"`
	ChannelID string         `json:"channel_id,omitempty"`
	Service   ChannelService `json:"service,omitempty"`
	IsLive    bool           `json:"is_live,omitempty"`
	Title     *string        `json:"title,omitempty"`
	Category  *string        `json:"category,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`

	Follows []*Follow `json:"follows,omitempty"`
	Streams []*Stream `json:"streams,omitempty"`
}
