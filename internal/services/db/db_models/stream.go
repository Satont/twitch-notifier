package db_models

import (
	"github.com/google/uuid"
	"time"
)

type Stream struct {
	ID         string     `json:"id,omitempty"`
	ChannelID  uuid.UUID  `json:"channel_id,omitempty"`
	Titles     []string   `json:"titles,omitempty"`
	Categories []string   `json:"categories,omitempty"`
	StartedAt  time.Time  `json:"started_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}
