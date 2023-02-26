package commands

import (
	"context"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
)

type FollowCommand struct {
	*tg_types.CommandOpts
}

func (c *FollowCommand) createFollow(ctx context.Context, chat *ent.Chat, input string) (*ent.Follow, error) {
	twitchChannel, err := c.Services.Twitch.GetUser("", input)
	if err != nil {
		return nil, err
	}

	spew.Dump(twitchChannel)
	if twitchChannel == nil {
		return nil, errors.New("channel not found")
	}

	dbChannel, err := c.Services.Channel.GetByIdOrCreate(ctx, twitchChannel.ID, channel.ServiceTwitch)
	if err != nil {
		return nil, err
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

	return msg.Answer("Created").DoVoid(ctx)
}

func NewFollowCommand(opts *tg_types.CommandOpts) {
	cmd := &FollowCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, tgb.Command("follow"))
}
