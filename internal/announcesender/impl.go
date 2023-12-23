package announcesender

import (
	"context"
	"fmt"

	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/i18n/localizer"
	"github.com/satont/twitch-notifier/internal/messagesender"
	"github.com/satont/twitch-notifier/internal/repository/channel"
	"github.com/satont/twitch-notifier/internal/repository/chat"
	"github.com/satont/twitch-notifier/internal/repository/follow"
	"github.com/satont/twitch-notifier/internal/thumbnailchecker"
)

type Opts struct {
	Localizer              localizer.Localizer
	FollowRepository       follow.Repository
	ChatRepository         chat.Repository
	ChatSettingsRepository chatsettings.Repository
	ChannelRepository      channel.Repository
	MessageSender          messagesender.MessageSender
	ThumbnailChecker       thumbnailchecker.ThumbnailChecker
}

func New(opts Opts) *AnnounceSenderImpl {
	return &AnnounceSenderImpl{
		localizer:              opts.Localizer,
		followRepository:       opts.FollowRepository,
		chatRepository:         opts.ChatRepository,
		chatSettingsRepository: opts.ChatSettingsRepository,
		channelRepository:      opts.ChannelRepository,
		messageSender:          opts.MessageSender,
		thumbnailChecker:       opts.ThumbnailChecker,
	}
}

var _ AnnounceSender = (*AnnounceSenderImpl)(nil)

type AnnounceSenderImpl struct {
	localizer              localizer.Localizer
	followRepository       follow.Repository
	chatRepository         chat.Repository
	chatSettingsRepository chatsettings.Repository
	channelRepository      channel.Repository
	messageSender          messagesender.MessageSender
	thumbnailChecker       thumbnailchecker.ThumbnailChecker
}

func (c *AnnounceSenderImpl) SendOnline(ctx context.Context, opts ChannelOnlineOpts) error {
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
				messagesender.Opts{
					Target: messagesender.MessageTarget{
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

func (c *AnnounceSenderImpl) SendOffline(ctx context.Context, opts ChannelOfflineOpts) error {
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
				messagesender.Opts{
					Target: messagesender.MessageTarget{
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

func (c *AnnounceSenderImpl) SendTitleChange(
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
				messagesender.Opts{
					Target: messagesender.MessageTarget{
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

func (c *AnnounceSenderImpl) SendCategoryChange(
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
				messagesender.Opts{
					Target: messagesender.MessageTarget{
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

func (c *AnnounceSenderImpl) SendTitleAndCategoryChange(
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
				messagesender.Opts{
					Target: messagesender.MessageTarget{
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
