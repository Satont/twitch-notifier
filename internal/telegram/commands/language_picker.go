package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/telegram/types"
	"go.uber.org/zap"
	"strings"
)

type LanguagePicker struct {
	*tgtypes.CommandOpts
}

func (c *LanguagePicker) buildKeyboard() (*tg.InlineKeyboardMarkup, error) {
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](1)

	codes := c.Services.I18N.GetLanguagesCodes()

	for _, code := range codes {
		layout.Add(
			tg.NewInlineKeyboardButtonCallback(
				fmt.Sprintf(
					"%s %s",
					c.Services.I18N.Translate("language.emoji", code, nil),
					c.Services.I18N.Translate("language.name", code, nil),
				),
				"language_picker_set_"+code,
			),
		)
	}

	layout.Add(tg.NewInlineKeyboardButtonCallback("«", "start_command_menu"))

	markup := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	return &markup, nil
}

func (c *LanguagePicker) HandleCallback(ctx context.Context, msg *tgb.CallbackQueryUpdate) error {
	keyboard, err := c.buildKeyboard()
	if err != nil {
		return msg.Answer().Text("internal error").DoVoid(ctx)
	}

	return msg.Client.
		EditMessageReplyMarkup(msg.Message.Chat.ID, msg.Message.ID).
		ReplyMarkup(*keyboard).
		DoVoid(ctx)
}

func (c *LanguagePicker) handleSetLanguage(ctx context.Context, msg *tgb.CallbackQueryUpdate) error {
	chat := c.SessionManager.Get(ctx).Chat
	if chat == nil {
		return errors.New("no chat")
	}

	lang := db_models.ChatLanguage(
		strings.TrimPrefix(msg.CallbackQuery.Data, "language_picker_set_"),
	)
	if !db_models.LanguageExists(lang) {
		return errors.New("language not exists")
	}

	_, err := c.Services.Chat.Update(
		ctx,
		msg.Message.Chat.ID.PeerID(),
		db_models.ChatServiceTelegram,
		&db.ChatUpdateQuery{
			Settings: &db.ChatUpdateSettingsQuery{
				ChatLanguage: &lang,
			},
		},
	)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	chat.Settings.ChatLanguage = lang

	err = msg.
		Answer().
		Text(c.Services.I18N.Translate("language.changed", lang.String(), nil)).
		DoVoid(ctx)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	return nil
}

func NewLanguagePicker(opts *tgtypes.CommandOpts) {
	picker := &LanguagePicker{opts}

	opts.Router.CallbackQuery(picker.HandleCallback, tgb.TextEqual("language_picker"))
	opts.Router.CallbackQuery(
		picker.handleSetLanguage,
		tgb.TextHasPrefix("language_picker_set_"),
	)
}
