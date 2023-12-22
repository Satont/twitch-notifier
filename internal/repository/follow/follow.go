package follow

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	ChannelID uuid.UUID
	CreatedAt time.Time
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Follow, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID) ([]Follow, error)
	GetByChannelID(ctx context.Context, channelID uuid.UUID) ([]Follow, error)
	Create(ctx context.Context, follow Follow) error
	Delete(ctx context.Context, id uuid.UUID) error
}
