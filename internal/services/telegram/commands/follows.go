package commands

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"go.uber.org/zap"
	"strings"
)

type FollowsCommand struct {
	*tgtypes.CommandOpts
}

func (c *FollowsCommand) newKeyboard(follows []*db_models.Follow) (*tg.InlineKeyboardMarkup, error) {
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](3)

	channelsIds := lo.Map(follows, func(follow *db_models.Follow, _ int) string {
		return follow.Channel.ChannelID
	})

	channels, err := c.Services.Twitch.GetChannelsByUserIds(channelsIds)

	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		layout.Insert(tg.NewInlineKeyboardButtonCallback(
			channel.BroadcasterName,
			fmt.Sprintf("channels_unfollow_%s", channel.BroadcasterID),
		))
	}

	markup := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	return &markup, nil
}

func (c *FollowsCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	follows, err := c.Services.Follow.GetByChatID(ctx, c.SessionManager.Get(ctx).Chat.ID)
	if err != nil {
		zap.S().Error(err)
		return msg.Answer("internal error").DoVoid(ctx)
	}

	keyboard, err := c.newKeyboard(follows)
	if err != nil {
		zap.S().Error(err)
		return msg.Answer("internal error").DoVoid(ctx)
	}

	return msg.Answer("your follows").ReplyMarkup(keyboard).DoVoid(ctx)
}

func (c *FollowsCommand) handleUnfollow(ctx context.Context, chat *db_models.Chat, input string) error {
	channelID := strings.Replace(input, "channels_unfollow_", "", 1)

	channel, err := c.Services.Channel.GetByID(ctx, channelID, db_models.ChannelServiceTwitch)
	if err != nil {
		return err
	}

	follow, err := c.Services.Follow.GetByChatAndChannel(ctx, channel.ID, chat.ID)
	if err != nil {
		return err
	}

	return c.Services.Follow.Delete(ctx, follow.ID)
}

func (c *FollowsCommand) unfollowQuery(ctx context.Context, msg *tgb.CallbackQueryUpdate) error {
	chat := c.SessionManager.Get(ctx).Chat

	if err := c.handleUnfollow(ctx, chat, msg.CallbackQuery.Data); err != nil {
		return msg.Answer().Text("internal error").DoVoid(ctx)
	}

	return msg.Answer().Text("unfollowed").DoVoid(ctx)
}

func NewFollowsCommand(opts *tgtypes.CommandOpts) {
	cmd := &FollowsCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, tgb.Command("follows"))
	opts.Router.CallbackQuery(cmd.unfollowQuery, tgb.TextHasPrefix("channels_unfollow_"))
}
