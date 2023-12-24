package twitchclient

import (
	"time"

	"github.com/satont/twitch-notifier/internal/domain"
)

type TwitchClient interface {
	GetChannelInformation(channelID string) (*domain.PlatformChannelInformation, error)
	GetLiveStream(channelID string) (*Stream, error)
}

type Stream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	TagIDs       []string  `json:"tag_ids"` //nolint:tagliatelle
	Tags         []string  `json:"tags"`
	IsMature     bool      `json:"is_mature"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
}
