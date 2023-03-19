package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/chat"
)

type ChatInterface interface {
	GetByID(_ context.Context, chatId string, service chat.Service) (*ent.Chat, error)
	GetFollowsByID(_ context.Context, chatId string, service chat.Service) ([]*ent.Follow, error)
	Create(_ context.Context, chatId string, service chat.Service) (*ent.Chat, error)
	Update(_ context.Context, chatId string, service chat.Service, settings *ent.ChatSettings) (*ent.Chat, error)
}
