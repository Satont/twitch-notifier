package commands

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"github.com/satont/twitch-notifier/internal/telegram/types"
	"go.uber.org/zap"
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

func (c *StartCommand) buildKeyboard(ctx context.Context) *tg.InlineKeyboardMarkup {
	chat := c.SessionManager.Get(ctx).Chat

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

	return &markup
}

func (c *StartCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	session := c.SessionManager.Get(ctx)

	keyBoard := c.buildKeyboard(ctx)

	description := c.Services.I18N.Translate(
		"bot.description",
		session.Chat.Settings.ChatLanguage.String(),
		nil,
	)

	return msg.Answer(description).ReplyMarkup(keyBoard).DoVoid(ctx)
}

func (c *StartCommand) handleCallback(ctx context.Context, msg *tgb.CallbackQueryUpdate) error {
	keyboard := c.buildKeyboard(ctx)

	return msg.Client.
		EditMessageReplyMarkup(msg.Message.Chat.ID, msg.Message.ID).
		ReplyMarkup(*keyboard).
		DoVoid(ctx)
}

func (c *StartCommand) handleGameNotificationSettings(
	ctx context.Context,
	msg *tgb.CallbackQueryUpdate,
) error {
	chat := c.SessionManager.Get(ctx).Chat

	chat.Settings.GameChangeNotification = !chat.Settings.GameChangeNotification

	_, err := c.Services.Chat.Update(
		ctx,
		chat.ChatID,
		db_models.ChatServiceTelegram,
		&db.ChatUpdateQuery{
			Settings: &db.ChatUpdateSettingsQuery{
				GameChangeNotification: &chat.Settings.GameChangeNotification,
			},
		},
	)
	if err != nil {
		zap.S().Error(err)
		return msg.Answer().Text("internal error").DoVoid(ctx)
	}

	keyboard := c.buildKeyboard(ctx)

	return msg.Client.
		EditMessageReplyMarkup(msg.Message.Chat.ID, msg.Message.ID).
		ReplyMarkup(*keyboard).
		DoVoid(ctx)
}

func (c *StartCommand) handleOfflineNotificationSettings(
	ctx context.Context,
	msg *tgb.CallbackQueryUpdate,
) error {
	chat := c.SessionManager.Get(ctx).Chat

	chat.Settings.OfflineNotification = !chat.Settings.OfflineNotification

	_, err := c.Services.Chat.Update(
		ctx,
		chat.ChatID,
		db_models.ChatServiceTelegram,
		&db.ChatUpdateQuery{
			Settings: &db.ChatUpdateSettingsQuery{
				OfflineNotification: &chat.Settings.OfflineNotification,
			},
		},
	)
	if err != nil {
		zap.S().Error(err)
		return msg.Answer().Text("internal error").DoVoid(ctx)
	}

	keyboard := c.buildKeyboard(ctx)

	return msg.Client.
		EditMessageReplyMarkup(msg.Message.Chat.ID, msg.Message.ID).
		ReplyMarkup(*keyboard).
		DoVoid(ctx)
}

var (
	startCommandFilter = tgb.Command("start",
		tgb.WithCommandAlias("help"),
		tgb.WithCommandAlias("info"),
		tgb.WithCommandAlias("settings"),
	)
	startMenuFilter                     = tgb.TextEqual("start_command_menu")
	gameChangeNotificationSettingFilter = tgb.TextEqual("start_game_change_notification_setting")
	offlineNotificationSettingFilter    = tgb.TextEqual("start_offline_notification")
)

func NewStartCommand(opts *tg_types.CommandOpts) {
	cmd := &StartCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, startCommandFilter)
	opts.Router.CallbackQuery(cmd.handleCallback, startMenuFilter)
	opts.Router.CallbackQuery(cmd.handleGameNotificationSettings, gameChangeNotificationSettingFilter)
	opts.Router.CallbackQuery(cmd.handleOfflineNotificationSettings, offlineNotificationSettingFilter)
}
