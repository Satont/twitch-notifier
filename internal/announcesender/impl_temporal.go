package announcesender

import (
	"context"
	"fmt"

	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/i18n/localizer"
	"github.com/satont/twitch-notifier/internal/messagesender"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

type Opts struct {
	Logger logger.Logger
}

func New(opts Opts) (*AnnounceSenderTemporal, error) {
	temporalClient, err := client.Dial(
		client.Options{
			Logger: log.NewStructuredLogger(opts.Logger.GetSlog()),
		},
	)
	if err != nil {
		return nil, err
	}

	return &AnnounceSenderTemporal{
		client: temporalClient,
	}, nil
}

var _ AnnounceSender = (*AnnounceSenderTemporal)(nil)

type AnnounceSenderTemporal struct {
	client client.Client
}

func (c *AnnounceSenderTemporal) SendOnline(ctx context.Context, opts ChannelOnlineOpts) error {
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

	return nil
}

func (c *AnnounceSenderTemporal) SendOffline(ctx context.Context, opts ChannelOfflineOpts) error {
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

		if !chatSettings.OfflineNotifications {
			continue
		}

		localizedString := c.localizer.MustLocalize(
			localizer.WithKey("offline"),
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
					Text: localizedString,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}
	}

	return nil
}

func (c *AnnounceSenderTemporal) SendTitleChange(
	ctx context.Context,
	opts ChannelTitleChangeOpts,
) error {
	followers, err := c.followRepository.GetByChannelID(ctx, opts.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}

	for _, follower := range followers {
		followerChat, err := c.chatRepository.GetByID(ctx, follower.ChatID)
		if err != nil {
			return fmt.Errorf("failed to get followerChat: %w", err)
		}

		chatSettings, err := c.chatSettingsRepository.GetByChatID(ctx, follower.ChatID)
		if err != nil {
			return fmt.Errorf("failed to get followerChat settings: %w", err)
		}

		if !chatSettings.TitleChangeNotifications {
			continue
		}

		localizedString := c.localizer.MustLocalize(
			localizer.WithKey("offline"),
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
					Text: localizedString,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}
	}

	return nil
}

func (c *AnnounceSenderTemporal) SendCategoryChange(
	ctx context.Context,
	opts ChannelCategoryChangeOpts,
) error {
	followers, err := c.followRepository.GetByChannelID(ctx, opts.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}

	for _, follower := range followers {
		followerChat, err := c.chatRepository.GetByID(ctx, follower.ChatID)
		if err != nil {
			return fmt.Errorf("failed to get followerChat: %w", err)
		}

		chatSettings, err := c.chatSettingsRepository.GetByChatID(ctx, follower.ChatID)
		if err != nil {
			return fmt.Errorf("failed to get followerChat settings: %w", err)
		}

		if !chatSettings.CategoryChangeNotifications {
			continue
		}

		localizedString := c.localizer.MustLocalize(
			localizer.WithKey("offline"),
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
					Text: localizedString,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}
	}

	return nil
}

func (c *AnnounceSenderTemporal) SendTitleAndCategoryChange(
	ctx context.Context,
	opts ChannelTitleAndCategoryChangeOpts,
) error {
	followers, err := c.followRepository.GetByChannelID(ctx, opts.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}

	for _, follower := range followers {
		followerChat, err := c.chatRepository.GetByID(ctx, follower.ChatID)
		if err != nil {
			return fmt.Errorf("failed to get followerChat: %w", err)
		}

		chatSettings, err := c.chatSettingsRepository.GetByChatID(ctx, follower.ChatID)
		if err != nil {
			return fmt.Errorf("failed to get followerChat settings: %w", err)
		}

		if !chatSettings.CategoryAndTitleNotifications {
			continue
		}

		localizedString := c.localizer.MustLocalize(
			localizer.WithKey("offline"),
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
					Text: localizedString,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}
	}

	return nil
}
