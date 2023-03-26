package commands

import (
	"context"
	"errors"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
)

type StartCommand struct {
	*tg_types.CommandOpts
}

func (c *StartCommand) buildKeyboard(ctx context.Context) (*tg.InlineKeyboardMarkup, error) {
	chat := c.SessionManager.Get(ctx).Chat
	if chat == nil {
		return nil, errors.New("chat is nil")
	}

	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](1)

	layout.Add(
		tg.NewInlineKeyboardButtonCallback(
			c.Services.I18N.Translate(
				"commands.start.game_change_notification_setting",
				chat.Settings.ChatLanguage.String(),
				nil,
			),
			"start_game_change_notification_setting",
		),
		tg.NewInlineKeyboardButtonCallback(
			c.Services.I18N.Translate(
				"commands.start.offline_notification",
				chat.Settings.ChatLanguage.String(),
				nil,
			),
			"start_offline_notification",
		),
		tg.NewInlineKeyboardButtonURL("Github", "https://github.com/Satont/twitch-notifier"),
	)

	markup := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	return &markup, nil
}

func (c *StartCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	keyBoard, err := c.buildKeyboard(ctx)
	if err != nil {
		return msg.Answer("internal error").DoVoid(ctx)
	}

	description := c.Services.I18N.Translate("bot.description", "en", nil)

	return msg.Answer(description).ReplyMarkup(keyBoard).DoVoid(ctx)
}

var startCommandFilter = tgb.Command("start",
	tgb.WithCommandAlias("help"),
	tgb.WithCommandAlias("info"),
	tgb.WithCommandAlias("settings"),
)

func NewStartCommand(opts *tg_types.CommandOpts) {
	cmd := &StartCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, startCommandFilter)
}
