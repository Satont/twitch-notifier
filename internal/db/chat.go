package db

import (
	"context"

	"github.com/satont/twitch-notifier/internal/db/db_models"
)

type ChatUpdateSettingsQuery struct {
	GameChangeNotification  *bool
	OfflineNotification     *bool
	TitleChangeNotification *bool
	ImageInNotification     *bool
	ChatLanguage            *db_models.ChatLanguage
}

type ChatUpdateQuery struct {
	Settings *ChatUpdateSettingsQuery
}

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
		query *ChatUpdateQuery,
	) (*db_models.Chat, error)
	GetAllByService(_ context.Context, service db_models.ChatService) ([]*db_models.Chat, error)
}
