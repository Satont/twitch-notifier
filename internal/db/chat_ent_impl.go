package db

import (
	"context"

	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/satont/twitch-notifier/ent/chatsettings"
	"github.com/satont/twitch-notifier/internal/db/db_models"
)

type chatService struct {
	entClient *ent.Client
}

func (c *chatService) convertEntity(entity *ent.Chat) *db_models.Chat {
	settings := &db_models.ChatSettings{
		ID:                             entity.Edges.Settings.ID,
		GameChangeNotification:         entity.Edges.Settings.GameChangeNotification,
		OfflineNotification:            entity.Edges.Settings.OfflineNotification,
		TitleChangeNotification:        entity.Edges.Settings.TitleChangeNotification,
		GameAndTitleChangeNotification: entity.Edges.Settings.GameAndTitleChangeNotification,
		ImageInNotification:            entity.Edges.Settings.ImageInNotification,
		ChatLanguage:                   db_models.ChatLanguage(entity.Edges.Settings.ChatLanguage),
		ChatID:                         entity.Edges.Settings.ChatID,
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
	settings *ChatUpdateQuery,
) (*db_models.Chat, error) {
	ch, err := c.entClient.Chat.
		Query().
		Where(chat.ChatID(chatId), chat.ServiceEQ(chat.Service(service))).
		WithSettings().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if settings.Settings != nil {
		updater := ch.Edges.Settings.Update()

		if settings.Settings.ChatLanguage != nil {
			updater.SetChatLanguage(chatsettings.ChatLanguage(*settings.Settings.ChatLanguage))
		}

		if settings.Settings.GameChangeNotification != nil {
			updater.SetGameChangeNotification(*settings.Settings.GameChangeNotification)
		}

		if settings.Settings.OfflineNotification != nil {
			updater.SetOfflineNotification(*settings.Settings.OfflineNotification)
		}

		if settings.Settings.TitleChangeNotification != nil {
			updater.SetTitleChangeNotification(*settings.Settings.TitleChangeNotification)
		}

		if settings.Settings.ImageInNotification != nil {
			updater.SetImageInNotification(*settings.Settings.ImageInNotification)
		}

		if settings.Settings.GameAndTitleChangeNotification != nil {
			updater.SetGameAndTitleChangeNotification(*settings.Settings.GameAndTitleChangeNotification)
		}

		_, err = updater.Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return c.GetByID(ctx, chatId, service)
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

func (c *chatService) GetAllByService(
	ctx context.Context,
	service db_models.ChatService,
) ([]*db_models.Chat, error) {
	chats, err := c.entClient.Chat.
		Query().
		Where(chat.ServiceEQ(chat.Service(service))).
		Order(ent.Desc(chat.FieldChatID)).
		WithSettings().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var result []*db_models.Chat
	for _, ch := range chats {
		result = append(result, c.convertEntity(ch))
	}

	return result, nil
}

func NewChatEntRepository(entClient *ent.Client) ChatInterface {
	return &chatService{
		entClient: entClient,
	}
}
