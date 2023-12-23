package follow

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
)

type Follow struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	ChannelID uuid.UUID
	CreatedAt time.Time
}

var ErrNotFound = errors.New("follow not found")
var ErrCannotCreate = errors.New("cannot create follow")
var ErrCannotDelete = errors.New("cannot delete follow")

//go:generate go run go.uber.org/mock/mockgen -source=follow.go -destination=mocks/mock.go

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Follow, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID) ([]domain.Follow, error)
	GetByChannelID(ctx context.Context, channelID uuid.UUID) ([]domain.Follow, error)
	Create(ctx context.Context, follow domain.Follow) error
	Delete(ctx context.Context, id uuid.UUID) error
}
