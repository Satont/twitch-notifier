package commands

import (
	"context"
	"errors"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"go.uber.org/zap"
)

type FollowCommand struct {
	*tgtypes.CommandOpts
}

var (
	channelNotFoundError = errors.New("channel not found")
)

func (c *FollowCommand) createFollow(ctx context.Context, chat *db_models.Chat, input string) (*db_models.Follow, error) {
	twitchChannel, err := c.Services.Twitch.GetUser("", input)
	if err != nil {
		return nil, err
	}

	if twitchChannel == nil {
		return nil, channelNotFoundError
	}

	dbChannel, err := c.Services.Channel.GetByIdOrCreate(ctx, twitchChannel.ID, db_models.ChannelServiceTwitch)
	if err != nil {
		return nil, err
	}

	follow, err := c.Services.Follow.Create(ctx, dbChannel.ID, chat.ID)
	if err != nil {
		return nil, err
	}

	return follow, nil
}

func (c *FollowCommand) handleScene(ctx context.Context, msg *tgb.MessageUpdate) error {
	chat := c.SessionManager.Get(ctx).Chat

	_, err := c.createFollow(ctx, chat, msg.Text)

	c.SessionManager.Get(ctx).Scene = ""

	if errors.Is(err, channelNotFoundError) {
		message := c.Services.I18N.Translate(
			"commands.follow.errors.streamerNotFound",
			chat.Settings.ChatLanguage.String(),
			map[string]string{
				"streamer": msg.Text,
			},
		)
		return msg.Answer(message).DoVoid(ctx)
	} else if errors.Is(err, db_models.FollowAlreadyExistsError) {
		message := c.Services.I18N.Translate(
			"commands.follow.alreadyFollowed",
			chat.Settings.ChatLanguage.String(),
			map[string]string{
				"streamer": msg.Text,
			},
		)
		return msg.Answer(message).DoVoid(ctx)
	} else if err != nil {
		zap.S().Error(err)
		return msg.Answer("Internal error").DoVoid(ctx)
	}

	message := c.Services.I18N.Translate(
		"commands.follow.success",
		chat.Settings.ChatLanguage.String(),
		map[string]string{
			"streamer": msg.Text,
		},
	)

	return msg.Answer(message).DoVoid(ctx)
}

func (c *FollowCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	c.SessionManager.Get(ctx).Scene = "follow"

	return msg.Answer("Enter name").DoVoid(ctx)
}

func NewFollowCommand(opts *tgtypes.CommandOpts) {
	cmd := &FollowCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.handleScene, tgb.FilterFunc(func(ctx context.Context, update *tgb.Update) (bool, error) {
		session := opts.SessionManager.Get(ctx)
		return session.Scene == "follow", nil
	}))
	opts.Router.Message(cmd.HandleCommand, tgb.Command("follow"))
}
