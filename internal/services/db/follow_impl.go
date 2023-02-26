package db

import (
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
)

type followService struct {
	entClient *ent.Client
}

func (f followService) Create(channelID string, channelService channel.Service, chatID string, chatService chat.Service) (*ent.Follow, error) {
	//TODO implement me
	panic("implement me")
}

func (f followService) Delete(channelID string, channelService channel.Service, chatID string, chatService chat.Service) error {
	//TODO implement me
	panic("implement me")
}

func NewFollowService(entClient *ent.Client) FollowInterface {
	return &followService{entClient: entClient}
}
