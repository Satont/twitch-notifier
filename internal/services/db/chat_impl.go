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

func (c *chatService) Update(chatId string, service chat.Service, settings *ent.ChatSettings) (*ent.Chat, error) {
	ch, err := c.entClient.Chat.
		Query().
		Where(chat.ChatID(chatId), chat.ServiceEQ(service)).
		Only(context.Background())
	if err != nil {
		return nil, err
	}

	newChat, err := c.entClient.Chat.
		UpdateOne(ch).
		SetSettings(settings).
		Save(context.Background())
	if err != nil {
		return nil, err
	}

	return newChat, nil
}

func (c *chatService) Create(chatId string, service chat.Service) (*ent.Chat, error) {
	ch, err := c.entClient.Chat.
		Create().
		SetChatID(chatId).
		SetService(service).
		Save(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = c.entClient.ChatSettings.Create().SetChatID(ch.ID).Save(context.Background())
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *chatService) GetByID(chatId string, service chat.Service) (*ent.Chat, error) {
	ch, err := c.entClient.Chat.
		Query().
		Where(chat.ChatID(chatId), chat.ServiceEQ(service)).
		Only(context.Background())

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

func (c *chatService) GetFollowsByID(chatId string, service chat.Service) ([]*ent.Follow, error) {
	follows, err := c.entClient.Follow.
		Query().
		Where(follow.HasChatWith(chat.ChatID(chatId), chat.ServiceEQ(service))).
		WithChat().
		WithChannel().
		All(context.Background())
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
