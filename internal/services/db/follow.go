package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
)

type FollowInterface interface {
	Create(_ context.Context, channelID uuid.UUID, chatID uuid.UUID) (*ent.Follow, error)
	Delete(_ context.Context, id uuid.UUID) error
	GetByChatAndChannel(_ context.Context, channelId uuid.UUID, chatId uuid.UUID) (*ent.Follow, error)
}
