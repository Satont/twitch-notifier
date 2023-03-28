package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/db/db_models"
)

type FollowInterface interface {
	Create(_ context.Context, channelID uuid.UUID, chatID uuid.UUID) (*db_models.Follow, error)
	Delete(_ context.Context, id uuid.UUID) error
	GetByChatAndChannel(
		_ context.Context,
		channelID uuid.UUID,
		chatID uuid.UUID,
	) (*db_models.Follow, error)
	GetByChannelID(_ context.Context, channelID uuid.UUID) ([]*db_models.Follow, error)
	GetByChatID(_ context.Context, chatID uuid.UUID, limit, offset int) ([]*db_models.Follow, error)
	CountByChatID(_ context.Context, chatID uuid.UUID) (int, error)
}
