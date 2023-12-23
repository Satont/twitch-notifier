package announcesender

import (
	"context"
	"fmt"

	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/i18n/localizer"
	"github.com/satont/twitch-notifier/internal/messagesender"
	"github.com/satont/twitch-notifier/internal/repository/channel"
	"github.com/satont/twitch-notifier/internal/repository/chat"
	"github.com/satont/twitch-notifier/internal/repository/chatsettings"
	"github.com/satont/twitch-notifier/internal/repository/follow"
	"github.com/satont/twitch-notifier/internal/thumbnailchecker"
	"go.temporal.io/sdk/workflow"
)

type AnnounceSenderWorkflows struct {
	localizer              localizer.Localizer
	followRepository       follow.Repository
	chatRepository         chat.Repository
	chatSettingsRepository chatsettings.Repository
	channelRepository      channel.Repository
	messageSender          messagesender.MessageSender
	thumbnailChecker       thumbnailchecker.ThumbnailChecker
}

func (c *AnnounceSenderWorkflows) SendOnline(
	workflowCtx workflow.Context,
	opts ChannelOnlineOpts,
) error {
	ctx := context.Background()

	err := c.thumbnailChecker.ValidateThumbnail(ctx, opts.ThumbnailURL)
	if err != nil {
		return fmt.Errorf("failed to check thumbnail: %w", err)
	}

	// db get
	followers, err := c.followRepository.GetByChannelID(ctx, opts.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}

	for _, follower := range followers {
		followerChat, err := c.chatRepository.GetByID(ctx, follower.ChatID)
		if err != nil {
			return fmt.Errorf("failed to get followerChat: %w", err)
		}

		chatSettings, err := c.chatSettingsRepository.GetByChatID(ctx, followerChat.ID)
		if err != nil {
			return fmt.Errorf("failed to get followerChat settings: %w", err)
		}

		localizedString := c.localizer.MustLocalize(
			localizer.WithKey("online"),
			localizer.WithLanguage(chatSettings.Language),
			localizer.WithAttribute("channelName", opts.ChannelID),
			localizer.WithAttribute("follower", follower),
		)

		if followerChat.Service == domain.ChatServiceTelegram {
			err = c.messageSender.SendMessageTelegram(
				ctx,
				messagesender.TelegramOpts{
					ServiceChatID: messagesender.MessageTarget{
						ServiceChatID: followerChat.ChatID,
					},
					Text:     localizedString,
					ImageURL: opts.ThumbnailURL,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}
	}
}
