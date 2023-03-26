package commands

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
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

func NewLanguagePicker(opts *tgtypes.CommandOpts) {
	picker := &LanguagePicker{opts}

	opts.Router.CallbackQuery(picker.HandleCallback, tgb.TextEqual("language_picker"))
}
