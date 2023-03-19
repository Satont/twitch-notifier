package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
)

type StreamUpdateQuery struct {
	StreamID string
	IsLive   *bool
	Category *string
	Title    *string
}

type StreamInterface interface {
	GetByID(_ context.Context, streamId string) (*db_models.Stream, error)

	GetLatestByChannelID(
		_ context.Context,
		channelEntityID uuid.UUID,
	) (*db_models.Stream, error)
	GetManyByChannelID(
		_ context.Context,
		channelEntityID uuid.UUID,
		limit int,
	) ([]*db_models.Stream, error)

	UpdateOneByStreamID(
		_ context.Context,
		streamID string,
		updateQuery *StreamUpdateQuery,
	) (*db_models.Stream, error)
	CreateOneByChannelID(
		_ context.Context,
		channelEntityID uuid.UUID,
		updateQuery *StreamUpdateQuery,
	) (*db_models.Stream, error)
}
