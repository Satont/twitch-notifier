package db

import (
	"context"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
)

type ChatInterface interface {
	GetByID(
		_ context.Context,
		chatId string,
		service db_models.ChatService,
	) (*db_models.Chat, error)
	Create(
		_ context.Context,
		chatId string,
		service db_models.ChatService,
	) (*db_models.Chat, error)
	Update(
		_ context.Context,
		chatId string,
		service db_models.ChatService,
		settings *db_models.ChatSettings,
	) (*db_models.Chat, error)
}
