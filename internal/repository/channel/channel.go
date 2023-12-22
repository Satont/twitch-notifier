package channel

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
)

type Channel struct {
	ID        uuid.UUID
	ChannelID string
	Service   domain.StreamingService
}

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (Channel, error)
	GetByStreamServiceAndID(ctx context.Context, service domain.StreamingService, id string) (
		Channel,
		error,
	)
	GetAll(ctx context.Context) ([]Channel, error)
	Create(ctx context.Context, channel Channel) error
	Update(ctx context.Context, channel Channel) error
	Delete(ctx context.Context, id uuid.UUID) error
}
