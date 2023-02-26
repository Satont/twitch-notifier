package db

import (
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
)

type ChatInterface interface {
	GetByID(chatId string, service chat.Service) (*ent.Chat, error)
	GetFollowsByID(chatId string, service chat.Service) ([]*ent.Follow, error)
	Create(chatId string, service chat.Service) (*ent.Chat, error)
	Update(chatId string, service chat.Service, settings *ent.ChatSettings) (*ent.Chat, error)
}

type ChannelUpdateQuery struct {
	IsLive   *bool
	Category *string
	Title    *string
}

type ChannelInterface interface {
	GetByID(channelID string, service channel.Service) (*ent.Channel, error)
	GetFollowsByID(channelID string, service channel.Service) ([]*ent.Follow, error)
	Create(channelID string, service channel.Service) (*ent.Channel, error)
	Update(channelID string, service channel.Service, updateQuery *ChannelUpdateQuery) (*ent.Channel, error)
}

type FollowInterface interface {
	Create(channelID string, channelService channel.Service, chatID string, chatService chat.Service) (*ent.Follow, error)
	Delete(channelID string, channelService channel.Service, chatID string, chatService chat.Service) error
}

type StreamUpdateQuery struct {
	IsLive   *bool
	Category *string
	Title    *string
}

type StreamInterface interface {
	GetByID(streamId string) (*ent.Stream, error)

	GetLatestByChannelID(channelEntityID uuid.UUID) (*ent.Stream, error)
	GetManyByChannelID(channelEntityID uuid.UUID, limit int) ([]*ent.Stream, error)

	UpdateOneByStreamID(streamID string, updateQuery *StreamUpdateQuery) (*ent.Stream, error)
	CreateOneByChannelID(channelEntityID uuid.UUID, updateQuery *StreamUpdateQuery) (*ent.Stream, error)
}
