package handler

import (
	"context"

	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/repository/channel"
)

type TwitchHandlerOpts struct {
	AnnounceSender    announcesender.AnnounceSender
	ChannelRepository channel.Repository
}

func NewTwitch(opts TwitchHandlerOpts) *TwitchHandler {
	return &TwitchHandler{
		announceSender:    opts.AnnounceSender,
		channelRepository: opts.ChannelRepository,
	}
}

var _ StreamsHandler = (*TwitchHandler)(nil)

type TwitchHandler struct {
	announceSender    announcesender.AnnounceSender
	channelRepository channel.Repository
}

func (c *TwitchHandler) Online(ctx context.Context, opts ChannelOnlineOpts) error {
	streamChannel, err := c.channelRepository.GetByStreamServiceAndID(
		ctx,
		domain.StreamingServiceTwitch,
		opts.ChannelID,
	)
	if err != nil {
		return err
	}

	return c.announceSender.SendOnline(
		ctx,
		announcesender.ChannelOnlineOpts{
			ChannelID:    streamChannel.ID,
			Category:     opts.Category.Name,
			Title:        opts.Title,
			ThumbnailURL: opts.ThumbnailURL,
		},
	)
}

func (c *TwitchHandler) Offline(ctx context.Context, opts ChannelOfflineOpts) error {
	streamChannel, err := c.channelRepository.GetByStreamServiceAndID(
		ctx,
		domain.StreamingServiceTwitch,
		opts.ChannelID,
	)
	if err != nil {
		return err
	}

	return c.announceSender.SendOffline(
		ctx,
		announcesender.ChannelOfflineOpts{
			ChannelID: streamChannel.ID,
		},
	)
}

func (c *TwitchHandler) MetadataChange(ctx context.Context, opts ChannelMetaDataChangedOpts) error {
	streamChannel, err := c.channelRepository.GetByStreamServiceAndID(
		ctx,
		domain.StreamingServiceTwitch,
		opts.ChannelID,
	)
	if err != nil {
		return err
	}

	// TODO: this is bad, need to rethink

	if opts.OldCategory.Name != opts.NewCategory.Name && opts.OldTitle != opts.NewTitle {
		return c.announceSender.SendTitleAndCategoryChange(
			ctx,
			announcesender.ChannelTitleAndCategoryChangeOpts{
				ChannelID:   streamChannel.ID,
				OldCategory: opts.OldCategory.Name,
				NewCategory: opts.NewCategory.Name,
				OldTitle:    opts.OldTitle,
			},
		)
	}

	if opts.OldCategory.Name != opts.NewCategory.Name {
		return c.announceSender.SendCategoryChange(
			ctx,
			announcesender.ChannelCategoryChangeOpts{
				ChannelID:   streamChannel.ID,
				OldCategory: opts.OldCategory.Name,
				NewCategory: opts.NewCategory.Name,
			},
		)
	}

	if opts.OldTitle != opts.NewTitle {
		return c.announceSender.SendTitleChange(
			ctx,
			announcesender.ChannelTitleChangeOpts{
				ChannelID: streamChannel.ID,
				OldTitle:  opts.OldTitle,
				NewTitle:  opts.NewTitle,
			},
		)
	}

	return nil
}
