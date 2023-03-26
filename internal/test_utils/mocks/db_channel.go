package mocks

import (
	"context"

	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/mock"
)

type DbChannelMock struct {
	mock.Mock
}

func (c *DbChannelMock) GetByID(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
) (*db_models.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*db_models.Channel), args.Error(1)
}

func (c *DbChannelMock) GetFollowsByID(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
) ([]*db_models.Follow, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).([]*db_models.Follow), args.Error(1)
}

func (c *DbChannelMock) Create(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
) (*db_models.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*db_models.Channel), args.Error(1)
}

func (c *DbChannelMock) Update(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
	updateQuery *db.ChannelUpdateQuery,
) (*db_models.Channel, error) {
	args := c.Called(ctx, channelID, service, updateQuery)

	return args.Get(0).(*db_models.Channel), args.Error(1)
}

func (c *DbChannelMock) GetByIdOrCreate(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
) (*db_models.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*db_models.Channel), args.Error(1)
}

func (c *DbChannelMock) GetAll(ctx context.Context) ([]*db_models.Channel, error) {
	args := c.Called(ctx)

	return args.Get(0).([]*db_models.Channel), args.Error(1)
}
