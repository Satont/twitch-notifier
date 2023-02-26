package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/stretchr/testify/mock"
)

type ChannelMock struct {
	mock.Mock
}

func (c ChannelMock) GetByID(ctx context.Context, channelID string, service channel.Service) (*ent.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*ent.Channel), args.Error(1)
}

func (c ChannelMock) GetFollowsByID(ctx context.Context, channelID string, service channel.Service) ([]*ent.Follow, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).([]*ent.Follow), args.Error(1)
}

func (c ChannelMock) Create(ctx context.Context, channelID string, service channel.Service) (*ent.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*ent.Channel), args.Error(1)
}

func (c ChannelMock) Update(ctx context.Context, channelID string, service channel.Service, updateQuery *ChannelUpdateQuery) (*ent.Channel, error) {
	args := c.Called(ctx, channelID, service, updateQuery)

	return args.Get(0).(*ent.Channel), args.Error(1)
}

func (c ChannelMock) GetByIdOrCreate(ctx context.Context, channelID string, service channel.Service) (*ent.Channel, error) {
	args := c.Called(ctx, channelID, service)

	return args.Get(0).(*ent.Channel), args.Error(1)
}
