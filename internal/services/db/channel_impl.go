package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/follow"
)

type channelService struct {
	entClient *ent.Client
}

func (c *channelService) GetByIdOrCreate(ctx context.Context, channelID string, service channel.Service) (*ent.Channel, error) {
	ch := c.entClient.Channel.
		Query().
		Where(channel.ChannelID(channelID), channel.ServiceEQ(service)).
		OnlyX(ctx)

	if ch == nil {
		newChannel, err := c.Create(ctx, channelID, service)
		if err != nil {
			return nil, err
		}
		ch = newChannel
	}

	return ch, nil
}

func (c *channelService) GetByID(ctx context.Context, channelID string, service channel.Service) (*ent.Channel, error) {
	ch, err := c.entClient.Channel.
		Query().
		Where(channel.ChannelID(channelID), channel.ServiceEQ(service)).
		Only(ctx)

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

func (c *channelService) GetFollowsByID(ctx context.Context, channelID string, service channel.Service) ([]*ent.Follow, error) {
	follows, err := c.entClient.Follow.
		Query().
		Where(follow.HasChannelWith(channel.ChannelID(channelID), channel.ServiceEQ(service))).
		WithChat().
		WithChannel().
		All(ctx)
	if err != nil {
		return nil, err
	}

	return follows, nil
}

func (c *channelService) Create(ctx context.Context, channelID string, service channel.Service) (*ent.Channel, error) {
	ch, err := c.entClient.Channel.Create().
		SetChannelID(channelID).
		SetService(service).Save(ctx)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *channelService) Update(ctx context.Context, channelID string, service channel.Service, query *ChannelUpdateQuery) (*ent.Channel, error) {
	ch, err := c.entClient.Channel.
		Query().
		Where(channel.ChannelIDIn(channelID), channel.ServiceEQ(service)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	updateQuery := c.entClient.Channel.UpdateOne(ch)

	if query.IsLive != nil {
		updateQuery.SetIsLive(*query.IsLive)
	}

	if query.Category != nil {
		updateQuery.SetCategory(*query.Category)
	}

	if query.Title != nil {
		updateQuery.SetTitle(*query.Title)
	}

	newChannel, err := updateQuery.Save(context.Background())

	if err != nil {
		return nil, err
	}

	return newChannel, nil
}

func NewChannelService(entClient *ent.Client) ChannelInterface {
	return &channelService{
		entClient: entClient,
	}
}
