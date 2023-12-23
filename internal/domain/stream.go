package domain

import (
	"time"

	"github.com/google/uuid"
)

type Stream struct {
	ID         uuid.UUID
	ChannelID  uuid.UUID
	Titles     []string
	Categories []string
	StartedAt  time.Time
	UpdatedAt  time.Time
	EndedAt    *time.Time
}
