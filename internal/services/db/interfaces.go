package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
)

type ChatInterface interface {
	GetByID(_ context.Context, chatId string, service chat.Service) (*ent.Chat, error)
	GetFollowsByID(_ context.Context, chatId string, service chat.Service) ([]*ent.Follow, error)
	Create(_ context.Context, chatId string, service chat.Service) (*ent.Chat, error)
	Update(_ context.Context, chatId string, service chat.Service, settings *ent.ChatSettings) (*ent.Chat, error)
}

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
}

type FollowInterface interface {
	Create(_ context.Context, channelID string, channelService channel.Service, chatID string, chatService chat.Service) (*ent.Follow, error)
	Delete(_ context.Context, channelID string, channelService channel.Service, chatID string, chatService chat.Service) error
}

type StreamUpdateQuery struct {
	IsLive   *bool
	Category *string
	Title    *string
}

type StreamInterface interface {
	GetByID(_ context.Context, streamId string) (*ent.Stream, error)

	GetLatestByChannelID(_ context.Context, channelEntityID uuid.UUID) (*ent.Stream, error)
	GetManyByChannelID(_ context.Context, channelEntityID uuid.UUID, limit int) ([]*ent.Stream, error)

	UpdateOneByStreamID(_ context.Context, streamID string, updateQuery *StreamUpdateQuery) (*ent.Stream, error)
	CreateOneByChannelID(_ context.Context, channelEntityID uuid.UUID, updateQuery *StreamUpdateQuery) (*ent.Stream, error)
}
