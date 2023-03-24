package db

import (
	"context"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/mock"
)

type ChatMock struct {
	mock.Mock
}

func (c *ChatMock) GetByID(ctx context.Context, chatId string, service db_models.ChatService) (*db_models.Chat, error) {
	args := c.Called(ctx, chatId, service)

	return args.Get(0).(*db_models.Chat), args.Error(1)
}

func (c *ChatMock) Create(ctx context.Context, chatId string, service db_models.ChatService) (*db_models.Chat, error) {
	args := c.Called(ctx, chatId, service)

	return args.Get(0).(*db_models.Chat), args.Error(1)
}

func (c *ChatMock) Update(ctx context.Context, chatId string, service db_models.ChatService, query *ChatUpdateQuery) (*db_models.Chat, error) {
	args := c.Called(ctx, chatId, service, query)

	return args.Get(0).(*db_models.Chat), args.Error(1)
}

func (c *ChatMock) GetAllByService(ctx context.Context, service db_models.ChatService) ([]*db_models.Chat, error) {
	args := c.Called(ctx, service)

	return args.Get(0).([]*db_models.Chat), args.Error(1)
}
