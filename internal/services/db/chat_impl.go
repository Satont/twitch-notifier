package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/satont/twitch-notifier/ent/follow"
)

type chatService struct {
	entClient *ent.Client
}

func (c *chatService) Update(ctx context.Context, chatId string, service chat.Service, settings *ent.ChatSettings) (*ent.Chat, error) {
	ch, err := c.entClient.Chat.
		Query().
		Where(chat.ChatID(chatId), chat.ServiceEQ(service)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	newChat, err := c.entClient.Chat.
		UpdateOne(ch).
		SetSettings(settings).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return newChat, nil
}

func (c *chatService) Create(ctx context.Context, chatId string, service chat.Service) (*ent.Chat, error) {
	ch, err := c.entClient.Chat.
		Create().
		SetChatID(chatId).
		SetService(service).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	settings, err := c.entClient.ChatSettings.Create().SetChatID(ch.ID).Save(ctx)
	if err != nil {
		return nil, err
	}

	ch.Edges.Settings = settings

	return ch, nil
}

func (c *chatService) GetByID(ctx context.Context, chatId string, service chat.Service) (*ent.Chat, error) {
	ch, err := c.entClient.Chat.
		Query().
		Where(chat.ChatID(chatId), chat.ServiceEQ(service)).
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

	return ch, nil
}

func (c *chatService) GetFollowsByID(ctx context.Context, chatId string, service chat.Service) ([]*ent.Follow, error) {
	follows, err := c.entClient.Follow.
		Query().
		Where(follow.HasChatWith(chat.ChatID(chatId), chat.ServiceEQ(service))).
		WithChat().
		WithChannel().
		All(ctx)
	if err != nil {
		return nil, err
	}

	return follows, nil
}

func NewChatService(entClient *ent.Client) ChatInterface {
	return &chatService{
		entClient: entClient,
	}
}
