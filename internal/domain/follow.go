package domain

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	ChannelID uuid.UUID
	CreatedAt time.Time
}
