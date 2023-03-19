package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
)

type channelEntService struct {
	entClient *ent.Client
}

func (c *channelEntService) convertEntity(ch *ent.Channel) *db_models.Channel {
	return &db_models.Channel{
		ID:        ch.ID,
		ChannelID: ch.ChannelID,
		Service:   db_models.ChannelService(ch.Service.String()),
		IsLive:    ch.IsLive,
		Title:     ch.Title,
		Category:  ch.Category,
		UpdatedAt: ch.UpdatedAt,
	}
}

func (c *channelEntService) GetByIdOrCreate(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
) (*db_models.Channel, error) {
	channelService := channel.Service(service.String())

	ch, err := c.entClient.Channel.
		Query().
		Where(channel.ChannelID(channelID), channel.ServiceEQ(channelService)).
		First(ctx)

	if ent.IsNotFound(err) {
		newChannel, err := c.Create(ctx, channelID, service)
		if err != nil {
			return nil, err
		}
		return newChannel, nil
	} else if err != nil {
		return nil, err
	}

	return c.convertEntity(ch), nil
}

func (c *channelEntService) GetByID(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
) (*db_models.Channel, error) {
	channelService := channel.Service(service.String())

	ch, err := c.entClient.Channel.
		Query().
		Where(channel.ChannelID(channelID), channel.ServiceEQ(channelService)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, db_models.ChannelNotFoundError
		}

		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return c.convertEntity(ch), nil
}

func (c *channelEntService) Create(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
) (*db_models.Channel, error) {
	channelService := channel.Service(service.String())

	ch, err := c.entClient.Channel.Create().
		SetChannelID(channelID).
		SetService(channelService).Save(ctx)
	if err != nil {
		return nil, err
	}

	return c.convertEntity(ch), nil
}

func (c *channelEntService) Update(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
	query *ChannelUpdateQuery,
) (*db_models.Channel, error) {
	channelService := channel.Service(service.String())
	ch, err := c.entClient.Channel.
		Query().
		Where(channel.ChannelIDIn(channelID), channel.ServiceEQ(channelService)).
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

	return c.convertEntity(newChannel), nil
}

func NewChannelEntService(entClient *ent.Client) ChannelInterface {
	return &channelEntService{
		entClient: entClient,
	}
}
