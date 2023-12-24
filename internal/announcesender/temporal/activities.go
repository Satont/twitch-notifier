package temporal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/repository/channel"
	"github.com/satont/twitch-notifier/internal/twitchclient"
	"go.uber.org/fx"
)

type ActivitiesOpts struct {
	fx.In

	ChannelRepository channel.Repository
	TwitchClient      twitchclient.TwitchClient
}

func NewActivities(opts ActivitiesOpts) *Activities {
	return &Activities{
		channelRepository: opts.ChannelRepository,
		twitchClient:      opts.TwitchClient,
	}
}

type Activities struct {
	channelRepository channel.Repository
	twitchClient      twitchclient.TwitchClient
}

var ErrUnknownService = fmt.Errorf("unknown service")

func (c *Activities) GetChannelInformation(ctx context.Context, channelID uuid.UUID) (
	*domain.PlatformChannelInformation,
	error,
) {
	channelEntity, err := c.channelRepository.GetById(ctx, channelID)
	if err != nil {
		return nil, err
	}

	if channelEntity.Service == domain.StreamingServiceTwitch {
		twitchChannel, err := c.twitchClient.GetChannelInformation(channelEntity.ChannelID)
		if err != nil {
			return nil, fmt.Errorf("failed to get twitch channel information: %w", err)
		}

		return twitchChannel, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnknownService, channelEntity.Service)
}
