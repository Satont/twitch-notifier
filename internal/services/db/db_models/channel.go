package db_models

import (
	"github.com/google/uuid"
	"time"
)

type ChannelService string

const (
	ServiceTwitch ChannelService = "twitch"
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
