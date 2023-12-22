package streams_handler

import (
	"context"

	"github.com/satont/twitch-notifier/internal/announce_sender"
	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/repository/channel"
)

type TwitchHandlerOpts struct {
	AnnounceSender    announce_sender.AnnounceSender
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
	announceSender    announce_sender.AnnounceSender
	channelRepository channel.Repository
}

func (c *TwitchHandler) Online(ctx context.Context, opts ChannelOnlineOpts) error {
	channel, err := c.channelRepository.GetByStreamServiceAndID(
		ctx,
		domain.StreamingServiceTwitch,
		opts.ChannelID,
	)
	if err != nil {
		return err
	}

	return c.announceSender.SendOnline(
		ctx,
		announce_sender.ChannelOnlineOpts{
			ChannelID:    channel.ID,
			Category:     opts.Category.Name,
			Title:        opts.Title,
			ThumbnailURL: opts.ThumbnailURL,
		},
	)
}

func (c *TwitchHandler) Offline(ctx context.Context, opts ChannelOfflineOpts) error {
	channel, err := c.channelRepository.GetByStreamServiceAndID(
		ctx,
		domain.StreamingServiceTwitch,
		opts.ChannelID,
	)
	if err != nil {
		return err
	}

	return c.announceSender.SendOffline(
		ctx,
		announce_sender.ChannelOfflineOpts{
			ChannelID: channel.ID,
		},
	)
}

func (c *TwitchHandler) MetadataChange(ctx context.Context, opts ChannelMetaDataChangedOpts) error {
	channel, err := c.channelRepository.GetByStreamServiceAndID(
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
			announce_sender.ChannelTitleAndCategoryChangeOpts{
				ChannelID:   channel.ID,
				OldCategory: opts.OldCategory.Name,
				NewCategory: opts.NewCategory.Name,
				OldTitle:    opts.OldTitle,
			},
		)
	}

	if opts.OldCategory.Name != opts.NewCategory.Name {
		return c.announceSender.SendCategoryChange(
			ctx,
			announce_sender.ChannelCategoryChangeOpts{
				ChannelID:   channel.ID,
				OldCategory: opts.OldCategory.Name,
				NewCategory: opts.NewCategory.Name,
			},
		)
	}

	if opts.OldTitle != opts.NewTitle {
		return c.announceSender.SendTitleChange(
			ctx,
			announce_sender.ChannelTitleChangeOpts{
				ChannelID: channel.ID,
				OldTitle:  opts.OldTitle,
				NewTitle:  opts.NewTitle,
			},
		)
	}

	return nil
}
