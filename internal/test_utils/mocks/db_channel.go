package mocks

import (
	"context"

	"github.com/satont/twitch-notifier/internal/db"
	db_models2 "github.com/satont/twitch-notifier/internal/db/db_models"

	"github.com/stretchr/testify/mock"
)

type DbChannelMock struct {
	mock.Mock
}

func (c *DbChannelMock) GetByID(
	ctx context.Context,
	id string,
	service db_models2.ChannelService,
) (*db_models2.Channel, error) {
	args := c.Called(ctx, id, service)

	return args.Get(0).(*db_models2.Channel), args.Error(1)
}

func (c *DbChannelMock) GetByChannelID(
	ctx context.Context,
	channelID string,
	service db_models2.ChannelService,
) (*db_models2.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*db_models2.Channel), args.Error(1)
}

func (c *DbChannelMock) GetFollowsByID(
	ctx context.Context,
	channelID string,
	service db_models2.ChannelService,
) ([]*db_models2.Follow, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).([]*db_models2.Follow), args.Error(1)
}

func (c *DbChannelMock) Create(
	ctx context.Context,
	channelID string,
	service db_models2.ChannelService,
) (*db_models2.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*db_models2.Channel), args.Error(1)
}

func (c *DbChannelMock) Update(
	ctx context.Context,
	channelID string,
	service db_models2.ChannelService,
	updateQuery *db.ChannelUpdateQuery,
) (*db_models2.Channel, error) {
	args := c.Called(ctx, channelID, service, updateQuery)

	return args.Get(0).(*db_models2.Channel), args.Error(1)
}

func (c *DbChannelMock) GetByIdOrCreate(
	ctx context.Context,
	channelID string,
	service db_models2.ChannelService,
) (*db_models2.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*db_models2.Channel), args.Error(1)
}

func (c *DbChannelMock) GetAll(ctx context.Context) ([]*db_models2.Channel, error) {
	args := c.Called(ctx)

	return args.Get(0).([]*db_models2.Channel), args.Error(1)
}
