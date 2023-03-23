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

const followsMaxRows = 3
const followsPerRow = 5

func (c *FollowsCommand) newKeyboard(ctx context.Context, maxRows, perRow int) (*tg.InlineKeyboardMarkup, error) {
	session := c.SessionManager.Get(ctx)

	limit := maxRows * perRow
	offset := (session.FollowsMenu.CurrentPage - 1) * limit

	follows, err := c.Services.Follow.GetByChatID(
		ctx,
		session.Chat.ID,
		limit,
		offset,
	)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	totalFollows, err := c.Services.Follow.CountByChatID(ctx, session.Chat.ID)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	session.FollowsMenu.TotalPages = totalFollows / limit

	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](perRow)

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

	var paginationRow *tg.ButtonLayout[tg.InlineKeyboardButton]

	if session.FollowsMenu.CurrentPage > 1 || session.FollowsMenu.CurrentPage < session.FollowsMenu.TotalPages {
		paginationRow = layout.Row()
	}

	fmt.Println(session.FollowsMenu.CurrentPage, session.FollowsMenu.TotalPages)
	if session.FollowsMenu.CurrentPage > 1 && paginationRow != nil {
		paginationRow.Insert(tg.NewInlineKeyboardButtonCallback(
			"«",
			"channels_unfollow_prev_page",
		))
	}

	if session.FollowsMenu.CurrentPage < session.FollowsMenu.TotalPages && paginationRow != nil {
		paginationRow.Insert(tg.NewInlineKeyboardButtonCallback(
			"»",
			"channels_unfollow_next_page",
		))
	}

	markup := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	return &markup, nil
}

func (c *FollowsCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	session := c.SessionManager.Get(ctx)

	session.FollowsMenu.TotalPages = 0
	session.FollowsMenu.CurrentPage = 1

	keyboard, err := c.newKeyboard(ctx, followsMaxRows, followsPerRow)
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

func (c *FollowsCommand) prevPageQuery(ctx context.Context, msg *tgb.CallbackQueryUpdate) error {
	session := c.SessionManager.Get(ctx)

	if session.FollowsMenu.CurrentPage > 0 {
		session.FollowsMenu.CurrentPage--
	}

	keyboard, err := c.newKeyboard(ctx, followsMaxRows, followsPerRow)
	if err != nil {
		zap.S().Error(err)
	}

	return msg.Client.
		EditMessageReplyMarkup(msg.Message.Chat.ID, msg.Message.ID).
		ReplyMarkup(*keyboard).
		DoVoid(ctx)
}

func (c *FollowsCommand) nextPageQuery(ctx context.Context, msg *tgb.CallbackQueryUpdate) error {
	session := c.SessionManager.Get(ctx)

	if session.FollowsMenu.CurrentPage+1 <= session.FollowsMenu.TotalPages {
		session.FollowsMenu.CurrentPage++
	}

	keyboard, err := c.newKeyboard(ctx, followsMaxRows, followsPerRow)
	if err != nil {
		zap.S().Error(err)
	}

	return msg.Client.
		EditMessageReplyMarkup(msg.Message.Chat.ID, msg.Message.ID).
		ReplyMarkup(*keyboard).
		DoVoid(ctx)
}

func NewFollowsCommand(opts *tgtypes.CommandOpts) {
	cmd := &FollowsCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(
		cmd.HandleCommand,
		tgb.Command(
			"follows",
			tgb.WithCommandAlias("unfollow"),
		),
	)
	opts.Router.CallbackQuery(cmd.prevPageQuery, tgb.TextEqual("channels_unfollow_prev_page"))
	opts.Router.CallbackQuery(cmd.nextPageQuery, tgb.TextEqual("channels_unfollow_next_page"))
	opts.Router.CallbackQuery(cmd.unfollowQuery, tgb.TextHasPrefix("channels_unfollow_"))
}
