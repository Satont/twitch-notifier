package db

import (
	"context"
	"github.com/satont/twitch-notifier/internal/db/db_models"
)

type ChannelUpdateQuery struct {
	IsLive   *bool
	Category *string
	Title    *string

	DangerNewChannelId *string
}

type ChannelInterface interface {
	GetByID(
		_ context.Context,
		channelID string,
		service db_models.ChannelService,
	) (*db_models.Channel, error)
	Create(_ context.Context, channelID string, service db_models.ChannelService) (*db_models.Channel, error)
	Update(
		_ context.Context,
		channelID string,
		service db_models.ChannelService,
		updateQuery *ChannelUpdateQuery,
	) (*db_models.Channel, error)
	GetByIdOrCreate(
		_ context.Context,
		channelID string,
		service db_models.ChannelService,
	) (*db_models.Channel, error)
	GetAll(_ context.Context) ([]*db_models.Channel, error)
}
