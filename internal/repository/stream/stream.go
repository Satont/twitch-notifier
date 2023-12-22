package stream

import (
	"context"

	"github.com/google/uuid"
)

type Stream struct {
	ID uuid.UUID
}

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (Stream, error)
	GetAll(ctx context.Context) ([]Stream, error)
	Create(ctx context.Context, stream Stream) error
	Update(ctx context.Context, stream Stream) error
	Delete(ctx context.Context, id uuid.UUID) error
}
