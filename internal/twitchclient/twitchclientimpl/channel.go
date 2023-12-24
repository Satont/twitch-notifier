package twitchclientimpl

import (
	"errors"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/internal/domain"
)

func (c *TwitchHelixImpl) GetChannelInformation(channelID string) (
	*domain.PlatformChannelInformation,
	error,
) {
	channelsReq, err := c.client.GetChannelInformation(
		&helix.GetChannelInformationParams{
			BroadcasterIDs: []string{channelID},
		},
	)
	if err != nil {
		return nil, err
	}
	if channelsReq.ErrorMessage != "" {
		return nil, fmt.Errorf("twitch api error: %s", channelsReq.ErrorMessage)
	}

	if len(channelsReq.Data.Channels) == 0 {
		return nil, fmt.Errorf("channel not found: %w", errors.New("not found"))
	}

	channel := channelsReq.Data.Channels[0]
	return &domain.PlatformChannelInformation{
		BroadcasterID:   channel.BroadcasterID,
		BroadcasterName: channel.BroadcasterName,
		GameName:        channel.GameName,
		Title:           channel.Title,
		ChannelLink:     fmt.Sprintf("https://twitch.tv/%s", channel.BroadcasterName),
	}, nil
}
