package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
)

type StreamUpdateQuery struct {
	IsLive   *bool
	Category *string
	Title    *string
}

type StreamInterface interface {
	GetByID(_ context.Context, streamId string) (*ent.Stream, error)

	GetLatestByChannelID(_ context.Context, channelEntityID uuid.UUID) (*ent.Stream, error)
	GetManyByChannelID(_ context.Context, channelEntityID uuid.UUID, limit int) ([]*ent.Stream, error)

	UpdateOneByStreamID(_ context.Context, streamID string, updateQuery *StreamUpdateQuery) (*ent.Stream, error)
	CreateOneByChannelID(_ context.Context, channelEntityID uuid.UUID, updateQuery *StreamUpdateQuery) (*ent.Stream, error)
}