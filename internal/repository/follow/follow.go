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
	UpdatedAt time.Time
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Follow, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID) ([]Follow, error)
	GetByChannelID(ctx context.Context, channelID uuid.UUID) ([]Follow, error)
	GetAll(ctx context.Context) ([]Follow, error)
	Create(ctx context.Context, follow Follow) error
	Update(ctx context.Context, follow Follow) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByChannelID(ctx context.Context, channelID uuid.UUID) ([]Follow, error)
	GetAllByChatID(ctx context.Context, chatID uuid.UUID) ([]Follow, error)
}
