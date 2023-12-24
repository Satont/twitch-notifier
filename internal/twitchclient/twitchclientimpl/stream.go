package twitchclientimpl

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/internal/twitchclient"
)

func (c *TwitchHelixImpl) GetLiveStream(channelID string) (*twitchclient.Stream, error) {
	streamsReq, err := c.client.GetStreams(
		&helix.StreamsParams{
			UserIDs: []string{channelID},
		},
	)
	if err != nil {
		return nil, err
	}
	if streamsReq.ErrorMessage != "" {
		return nil, fmt.Errorf("failed to get streams: %s", streamsReq.ErrorMessage)
	}

	if len(streamsReq.Data.Streams) == 0 {
		return nil, nil
	}

	stream := streamsReq.Data.Streams[0]
	return &twitchclient.Stream{
		ID:           stream.ID,
		UserID:       stream.UserID,
		UserLogin:    stream.UserLogin,
		UserName:     stream.UserName,
		GameID:       stream.GameID,
		GameName:     stream.GameName,
		TagIDs:       stream.TagIDs,
		Tags:         stream.Tags,
		IsMature:     stream.IsMature,
		Type:         stream.Type,
		Title:        stream.Title,
		ViewerCount:  stream.ViewerCount,
		StartedAt:    stream.StartedAt,
		Language:     stream.Language,
		ThumbnailURL: stream.ThumbnailURL,
	}, nil
}
