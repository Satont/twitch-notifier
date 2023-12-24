package domain

import (
	"github.com/google/uuid"
)

type Channel struct {
	ID        uuid.UUID
	ChannelID string
	Service   StreamingService
}

type PlatformChannelInformation struct {
	BroadcasterID   string
	BroadcasterName string
	GameName        string
	Title           string
	ChannelLink     string
}
