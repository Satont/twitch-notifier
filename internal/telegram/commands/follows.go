package commands

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/telegram/types"
	"go.uber.org/zap"
	"math"
	"strings"
)

type FollowsCommand struct {
	*tgtypes.CommandOpts
}

const followsMaxRows = 3
const followsPerRow = 3

func (c *FollowsCommand) newKeyboard(
	ctx context.Context,
	maxRows, perRow int,
) (*tg.InlineKeyboardMarkup, error) {
	session := c.SessionManager.Get(ctx)

	limit := maxRows * perRow
	offset := (session.FollowsMenu.CurrentPage - 1) * limit

	if offset < 0 {
		offset = 0
	}

	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](perRow)

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
	if len(follows) == 0 {
		markup := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)
		return &markup, nil
	}

	totalFollows, err := c.Services.Follow.CountByChatID(ctx, session.Chat.ID)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	session.FollowsMenu.TotalPages = int(math.Ceil(float64(totalFollows) / float64(limit)))
	//spew.Dump(totalFollows)
	//spew.Dump(session.FollowsMenu)
	//spew.Dump(session.FollowsMenu.CurrentPage)

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

	if session.FollowsMenu.CurrentPage > 1 ||
		session.FollowsMenu.CurrentPage < session.FollowsMenu.TotalPages {
		paginationRow = layout.Row()

		// Add "Prev" button
		if session.FollowsMenu.CurrentPage > 1 {
			paginationRow.Insert(tg.NewInlineKeyboardButtonCallback(
				"«",
				"channels_unfollow_prev_page",
			))
		}

		// Add "Next" button
		if session.FollowsMenu.CurrentPage < session.FollowsMenu.TotalPages {
			paginationRow.Insert(tg.NewInlineKeyboardButtonCallback(
				"»",
				"channels_unfollow_next_page",
			))
		}
	}

	markup := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	return &markup, nil
}

func (c *FollowsCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	session := c.SessionManager.Get(ctx)

	session.FollowsMenu.TotalPages = 1
	session.FollowsMenu.CurrentPage = 1

	keyboard, err := c.newKeyboard(ctx, followsMaxRows, followsPerRow)
	if err != nil {
		zap.S().Error(err)
		return msg.Answer("internal error").DoVoid(ctx)
	}

	totalFollows, err := c.Services.Follow.CountByChatID(ctx, session.Chat.ID)

	return msg.
		Answer(c.Services.I18N.Translate(
			"commands.follows.total",
			session.Chat.Settings.ChatLanguage.String(),
			map[string]string{"count": fmt.Sprintf("%v", totalFollows)},
		)).
		ReplyMarkup(keyboard).DoVoid(ctx)
}

func (c *FollowsCommand) handleUnfollow(
	ctx context.Context,
	chat *db_models.Chat,
	input string,
) error {
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
		if err == db_models.FollowNotFoundError {
			return msg.Answer().Text("already unfollowed").DoVoid(ctx)
		}

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

var (
	followsCommandFilter = tgb.Command(
		"follows",
		tgb.WithCommandAlias("unfollow"))
	followsPrevPageQuery = tgb.TextEqual("channels_unfollow_prev_page")
	followsNextPageQuery = tgb.TextEqual("channels_unfollow_next_page")
	followUnfollowQuery  = tgb.TextHasPrefix("channels_unfollow_")
)

func NewFollowsCommand(opts *tgtypes.CommandOpts) {
	cmd := &FollowsCommand{
		CommandOpts: opts,
	}

	messageFilter := []tgb.Filter{
		channelsAdminFilter,
		followsCommandFilter,
	}

	opts.Router.Message(cmd.HandleCommand, messageFilter...)
	opts.Router.ChannelPost(cmd.HandleCommand, messageFilter...)

	opts.Router.CallbackQuery(cmd.prevPageQuery, channelsAdminFilter, followsPrevPageQuery)
	opts.Router.CallbackQuery(cmd.nextPageQuery, channelsAdminFilter, followsNextPageQuery)
	opts.Router.CallbackQuery(cmd.unfollowQuery, channelsAdminFilter, followUnfollowQuery)
}
