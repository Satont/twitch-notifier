package channel

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
)

type Channel struct {
	ID        uuid.UUID
	ChannelID string
	Service   StreamingService
}

type StreamingService string

func (s StreamingService) String() string {
	return string(s)
}

const (
	StreamingServiceTwitch StreamingService = "twitch"
)

var ErrNotFound = errors.New("channel not found")
var ErrCannotCreate = errors.New("cannot create channel")
var ErrCannotDelete = errors.New("cannot delete channel")

//go:generate go run go.uber.org/mock/mockgen -source=channel.go -destination=mocks/mock.go

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.Channel, error)
	GetByStreamServiceAndID(ctx context.Context, service StreamingService, id string) (
		*domain.Channel,
		error,
	)
	GetAll(ctx context.Context) ([]domain.Channel, error)
	Create(ctx context.Context, channel domain.Channel) error
	// Update(ctx context.Context, channel Channel) error
	Delete(ctx context.Context, id uuid.UUID) error
}
