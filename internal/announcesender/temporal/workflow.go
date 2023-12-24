package temporal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/satont/twitch-notifier/internal/announcesender"
	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/i18n/localizer"
	"github.com/satont/twitch-notifier/internal/messagesender"
	"github.com/satont/twitch-notifier/internal/repository/chat"
	"github.com/satont/twitch-notifier/internal/repository/chatsettings"
	"github.com/satont/twitch-notifier/internal/repository/follow"
	"github.com/satont/twitch-notifier/internal/thumbnailchecker"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
)

type WorkflowOpts struct {
	fx.In

	Localizer              localizer.Localizer
	FollowRepository       follow.Repository
	ChatRepository         chat.Repository
	ChatSettingsRepository chatsettings.Repository
	MessageSender          messagesender.MessageSender
	ThumbnailChecker       thumbnailchecker.ThumbnailChecker
	Activities             *Activities
}

func NewWorkflow(opts WorkflowOpts) *Workflow {
	return &Workflow{
		localizer:              opts.Localizer,
		messageSender:          opts.MessageSender,
		thumbnailChecker:       opts.ThumbnailChecker,
		followRepository:       opts.FollowRepository,
		chatRepository:         opts.ChatRepository,
		chatSettingsRepository: opts.ChatSettingsRepository,
		activities:             opts.Activities,
	}
}

type Workflow struct {
	localizer        localizer.Localizer
	messageSender    messagesender.MessageSender
	thumbnailChecker thumbnailchecker.ThumbnailChecker

	followRepository       follow.Repository
	chatRepository         chat.Repository
	chatSettingsRepository chatsettings.Repository
	activities             *Activities
}

func (c *Workflow) SendOnline(
	workflowCtx workflow.Context,
	opts announcesender.ChannelOnlineOpts,
) error {
	ctx := context.Background()

	logger := workflow.GetLogger(workflowCtx)
	logger.Info("Sending online message", "channelId", opts.ChannelID.String())

	var channelInformation *domain.PlatformChannelInformation
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(
			workflowCtx,
			workflow.ActivityOptions{
				TaskQueue: queueName,
				RetryPolicy: &temporal.RetryPolicy{
					MaximumInterval: 15 * time.Second,
					MaximumAttempts: 10,
				},
			},
		),
		c.activities.GetChannelInformation,
		opts.ChannelID,
	).Get(workflowCtx, &channelInformation)
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	// TODO: execute this as child workflow
	err = c.thumbnailChecker.ValidateThumbnail(ctx, opts.ThumbnailURL)
	if err != nil {
		return fmt.Errorf("failed to check thumbnail: %w", err)
	}

	// db get
	followers, err := c.followRepository.GetByChannelID(ctx, opts.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}

	logger.Info("Got followers", "followers", len(followers))

	var followersSendErrorsErrors []error
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
			localizer.WithKey("notifications.streams.nowOnline"),
			localizer.WithLanguage(chatSettings.Language),
			localizer.WithAttribute("channelLink", channelInformation.ChannelLink),
			localizer.WithAttribute("category", channelInformation.GameName),
			localizer.WithAttribute("title", channelInformation.Title),
		)

		logger.Info(
			"Sending message",
			"followerChat",
			followerChat.ID,
			"service",
			followerChat.Service,
		)

		if followerChat.Service == domain.ChatServiceTelegram {
			err = c.messageSender.SendMessageTelegram(
				ctx,
				messagesender.TelegramOpts{
					ServiceChatID: followerChat.ChatID,
					Text:          localizedString,
					ImageURL:      opts.ThumbnailURL,
				},
			)
			if err != nil {
				followersSendErrorsErrors = append(followersSendErrorsErrors, err)
			}
		}
	}

	if len(followersSendErrorsErrors) > 0 {
		logger.Error("Failed to send messages", "errors", followersSendErrorsErrors)
		return errors.Join(followersSendErrorsErrors...)
	}

	return nil
}
