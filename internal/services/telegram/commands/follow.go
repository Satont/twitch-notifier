package commands

import (
	"context"
	"errors"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
)

type FollowCommand struct {
	*tgtypes.CommandOpts
}

var (
	channelNotFoundError = errors.New("channel not found")
	followAlreadyExists  = errors.New("follow already exists")
)

func (c *FollowCommand) createFollow(ctx context.Context, chat *ent.Chat, input string) (*ent.Follow, error) {
	twitchChannel, err := c.Services.Twitch.GetUser("", input)
	if err != nil {
		return nil, err
	}

	if twitchChannel == nil {
		return nil, channelNotFoundError
	}

	dbChannel, err := c.Services.Channel.GetByIdOrCreate(ctx, twitchChannel.ID, channel.ServiceTwitch)
	if err != nil {
		return nil, err
	}

	existedFollow, err := c.Services.Follow.GetByChatAndChannel(ctx, chat.ID, dbChannel.ID)
	if err != nil {
		return nil, err
	}
	if existedFollow != nil {
		return nil, followAlreadyExists
	}

	follow, err := c.Services.Follow.Create(ctx, dbChannel.ID, chat.ID)
	if err != nil {
		return nil, err
	}

	return follow, nil
}

func (c *FollowCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	chat := c.SessionManager.Get(ctx).Chat

	_, err := c.createFollow(ctx, chat, msg.Text)
	if err != nil {
		return err
	}

	message := c.Services.I18N.Translate(
		"commands.follow.success",
		chat.Edges.Settings.ChatLanguage.String(),
		map[string]string{
			"streamer": msg.Text,
		},
	)

	return msg.Answer(message).DoVoid(ctx)
}

func NewFollowCommand(opts *tgtypes.CommandOpts) {
	cmd := &FollowCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, tgb.Command("follow"))
}
