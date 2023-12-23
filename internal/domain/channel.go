package domain

import (
	"github.com/google/uuid"
)

type Channel struct {
	ID        uuid.UUID
	ChannelID string
	Service   StreamingService
}
