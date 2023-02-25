package db

import (
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

type ChannelInterface interface {
	GetByID(channelID string, service channel.Service) (*ent.Channel, error)
	GetFollowsByID(channelID string, service channel.Service) ([]*ent.Follow, error)
	Create(channelID string, service channel.Service) (*ent.Channel, error)
	Update(channelID string, service channel.Service) (*ent.Channel, error)
}
