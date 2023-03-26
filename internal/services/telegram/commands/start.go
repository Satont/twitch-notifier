package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
)

type StartCommand struct {
	*tg_types.CommandOpts
}

func (c *StartCommand) createCheckMark(value bool) string {
	if value {
		return "✅"
	}

	return "❌"
}

func (c *StartCommand) buildKeyboard(ctx context.Context) (*tg.InlineKeyboardMarkup, error) {
	chat := c.SessionManager.Get(ctx).Chat
	if chat == nil {
		return nil, errors.New("chat is nil")
	}

	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](1)

	gameChangeNotificationsButton := c.Services.I18N.Translate(
		"commands.start.game_change_notification_setting.button",
		chat.Settings.ChatLanguage.String(),
		nil,
	)
	offlineNotificationsButton := c.Services.I18N.Translate(
		"commands.start.offline_notification.button",
		chat.Settings.ChatLanguage.String(),
		nil,
	)

	layout.Add(
		tg.NewInlineKeyboardButtonCallback(
			fmt.Sprintf(
				"%s %s",
				c.createCheckMark(chat.Settings.GameChangeNotification),
				gameChangeNotificationsButton,
			),
			"start_game_change_notification_setting",
		),
		tg.NewInlineKeyboardButtonCallback(
			fmt.Sprintf(
				"%s %s",
				c.createCheckMark(chat.Settings.OfflineNotification),
				offlineNotificationsButton,
			),
			"start_offline_notification",
		),
		tg.NewInlineKeyboardButtonCallback(
			c.Services.I18N.Translate(
				"commands.start.language.button",
				chat.Settings.ChatLanguage.String(),
				nil,
			),
			"language_picker",
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

var (
	startCommandFilter = tgb.Command("start",
		tgb.WithCommandAlias("help"),
		tgb.WithCommandAlias("info"),
		tgb.WithCommandAlias("settings"),
	)
	startMenuFilter = tgb.TextEqual("start_command_menu")
)

func (c *StartCommand) handleCallback(ctx context.Context, msg *tgb.CallbackQueryUpdate) error {
	keyboard, err := c.buildKeyboard(ctx)
	if err != nil {
		return msg.Answer().Text("internal error").DoVoid(ctx)
	}

	return msg.Client.
		EditMessageReplyMarkup(msg.Message.Chat.ID, msg.Message.ID).
		ReplyMarkup(*keyboard).
		DoVoid(ctx)
}

func NewStartCommand(opts *tg_types.CommandOpts) {
	cmd := &StartCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, startCommandFilter)
	opts.Router.CallbackQuery(cmd.handleCallback, startMenuFilter)
}
