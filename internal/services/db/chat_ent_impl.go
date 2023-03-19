package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/satont/twitch-notifier/ent/chatsettings"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
)

type chatService struct {
	entClient *ent.Client
}

func (c *chatService) convertEntity(entity *ent.Chat) *db_models.Chat {
	settings := &db_models.ChatSettings{
		ID:                     entity.Edges.Settings.ID,
		GameChangeNotification: entity.Edges.Settings.GameChangeNotification,
		OfflineNotification:    entity.Edges.Settings.OfflineNotification,
		ChatLanguage:           db_models.ChatLanguage(entity.Edges.Settings.ChatLanguage),
		ChatID:                 entity.Edges.Settings.ChatID,
	}

	return &db_models.Chat{
		ID:       entity.ID,
		ChatID:   entity.ChatID,
		Service:  db_models.ChatService(entity.Service),
		Settings: settings,
	}
}

func (c *chatService) Update(
	ctx context.Context,
	chatId string,
	service db_models.ChatService,
	settings *db_models.ChatSettings,
) (*db_models.Chat, error) {
	ch, err := c.entClient.Chat.
		Query().
		Where(chat.ChatID(chatId), chat.ServiceEQ(chat.Service(service))).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	newSettings := &ent.ChatSettings{
		ID:                     settings.ID,
		GameChangeNotification: settings.GameChangeNotification,
		OfflineNotification:    settings.OfflineNotification,
		ChatLanguage:           chatsettings.ChatLanguage(db_models.ChatLanguageEn),
	}

	newChat, err := c.entClient.Chat.
		UpdateOne(ch).
		SetSettings(newSettings).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return c.convertEntity(newChat), nil
}

func (c *chatService) Create(
	ctx context.Context,
	chatId string,
	service db_models.ChatService,
) (*db_models.Chat, error) {
	ch, err := c.entClient.Chat.
		Create().
		SetChatID(chatId).
		SetService(chat.Service(service.String())).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	settings, err := c.entClient.ChatSettings.Create().SetChatID(ch.ID).Save(ctx)
	if err != nil {
		return nil, err
	}

	ch.Edges.Settings = settings

	return c.convertEntity(ch), nil
}

func (c *chatService) GetByID(
	ctx context.Context,
	chatId string,
	service db_models.ChatService,
) (*db_models.Chat, error) {
	ch, err := c.entClient.Chat.
		Query().
		Where(chat.ChatID(chatId), chat.ServiceEQ(chat.Service(service))).
		WithSettings().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return c.convertEntity(ch), nil
}

func NewChatEntRepository(entClient *ent.Client) ChatInterface {
	return &chatService{
		entClient: entClient,
	}
}
