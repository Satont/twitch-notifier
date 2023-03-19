package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
)

type FollowInterface interface {
	Create(_ context.Context, channelID uuid.UUID, chatID uuid.UUID) (*db_models.Follow, error)
	Delete(_ context.Context, id uuid.UUID) error
	GetByChatAndChannel(
		_ context.Context,
		channelId uuid.UUID,
		chatId uuid.UUID,
	) (*db_models.Follow, error)
	GetByChannelID(_ context.Context, channelId uuid.UUID) ([]*db_models.Follow, error)
}
