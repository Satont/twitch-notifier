package stream

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type Stream struct {
	ID         uuid.UUID
	ChannelID  uuid.UUID
	Titles     []string
	Categories []string
	StartedAt  time.Time
	UpdatedAt  time.Time
	EndedAt    null.Time
}

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (Stream, error)
	GetLatestByChannelId(ctx context.Context, channelId uuid.UUID) (Stream, error)
	GetByChannelId(ctx context.Context, channelId uuid.UUID) ([]Stream, error)
	Create(ctx context.Context, stream Stream) error
	Update(ctx context.Context, stream Stream) error
	Delete(ctx context.Context, id uuid.UUID) error
}
