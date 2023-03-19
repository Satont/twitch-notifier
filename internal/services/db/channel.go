package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
)

type ChannelUpdateQuery struct {
	IsLive   *bool
	Category *string
	Title    *string
}

type ChannelInterface interface {
	GetByID(_ context.Context, channelID string, service channel.Service) (*ent.Channel, error)
	GetFollowsByID(_ context.Context, channelID string, service channel.Service) ([]*ent.Follow, error)
	Create(_ context.Context, channelID string, service channel.Service) (*ent.Channel, error)
	Update(_ context.Context, channelID string, service channel.Service, updateQuery *ChannelUpdateQuery) (*ent.Channel, error)
	GetByIdOrCreate(_ context.Context, channelID string, service channel.Service) (*ent.Channel, error)
}
