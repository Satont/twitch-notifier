package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/follow"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
)

type channelService struct {
	entClient *ent.Client
}

func (c *channelService) convertEntity(ch *ent.Channel) *db_models.Channel {
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

func (c *channelService) GetByIdOrCreate(
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

func (c *channelService) GetByID(
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
			return nil, nil
		}

		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return c.convertEntity(ch), nil
}

func (c *channelService) GetFollowsByID(
	ctx context.Context,
	channelID string,
	service db_models.ChannelService,
) ([]*db_models.Follow, error) {
	channelService := channel.Service(service.String())

	follows, err := c.entClient.Follow.
		Query().
		Where(follow.HasChannelWith(
			channel.ChannelID(channelID),
			channel.ServiceEQ(channelService)),
		).
		WithChat().
		WithChannel().
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*db_models.Follow, 0, len(follows))

	for _, f := range follows {
		result = append(result, &db_models.Follow{
			ID:        f.ID,
			ChatID:    f.Edges.Chat.ID,
			ChannelID: f.Edges.Channel.ID,
		})
	}

	return result, nil
}

func (c *channelService) Create(
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

func (c *channelService) Update(
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

func NewChannelEntRepository(entClient *ent.Client) ChannelInterface {
	return &channelService{
		entClient: entClient,
	}
}
