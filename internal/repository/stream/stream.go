package stream

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
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

var ErrNotFound = errors.New("stream not found")
var ErrCannotCreate = errors.New("cannot create stream")
var ErrCannotDelete = errors.New("cannot delete stream")

//go:generate go run go.uber.org/mock/mockgen -source=stream.go -destination=mocks/mock.go

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.Stream, error)
	GetLatestByChannelId(ctx context.Context, channelId uuid.UUID) (*domain.Stream, error)
	GetByChannelId(ctx context.Context, channelId uuid.UUID) ([]domain.Stream, error)
	Create(ctx context.Context, stream domain.Stream) error
	Update(ctx context.Context, stream domain.Stream) error
	Delete(ctx context.Context, id uuid.UUID) error
}
