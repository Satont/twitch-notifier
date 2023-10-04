package message_sender

import (
	"context"
	"strconv"

	"github.com/mr-linch/go-tg"
	"github.com/satont/twitch-notifier/internal/db/db_models"
)

type MessageSender struct {
	telegram *tg.Client
}

func (m *MessageSender) SendMessage(ctx context.Context, chat *db_models.Chat, opts *MessageOpts) error {
	if chat.Service == db_models.ChatServiceTelegram {
		chatId, err := strconv.Atoi(chat.ChatID)
		if err != nil {
			return err
		}

		var keyboard *tg.InlineKeyboardMarkup
		if opts.Buttons != nil && len(opts.Buttons) > 0 {
			keyboard = &tg.InlineKeyboardMarkup{
				InlineKeyboard: make([][]tg.InlineKeyboardButton, 0, len(opts.Buttons)),
			}

			for _, row := range opts.Buttons {
				var buttons []tg.InlineKeyboardButton
				for _, button := range row {
					if button.SkipInGroup {
						continue
					}

					buttons = append(
						buttons, tg.InlineKeyboardButton{
							Text:         button.Text,
							CallbackData: button.CallbackData,
						},
					)
				}

				if len(buttons) != 0 {
					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, buttons)
				}
			}
		}

		if opts.ImageURL != "" {
			query := m.telegram.
				SendPhoto(tg.ChatID(chatId), tg.FileArg{URL: opts.ImageURL}).
				Caption(opts.Text)

			if opts.ParseMode != nil {
				query = query.ParseMode(*opts.ParseMode)
			}

			if keyboard != nil && keyboard.InlineKeyboard != nil && len(keyboard.InlineKeyboard) > 0 {
				query = query.ReplyMarkup(keyboard)
			}

			return query.DoVoid(ctx)
		} else {
			query := m.telegram.
				SendMessage(tg.ChatID(chatId), opts.Text).
				DisableWebPagePreview(true)

			if keyboard != nil && keyboard.InlineKeyboard != nil && len(keyboard.InlineKeyboard) > 0 {
				query = query.ReplyMarkup(keyboard)
			}

			if opts.ParseMode != nil {
				query = query.ParseMode(*opts.ParseMode)
			}

			return query.DoVoid(ctx)
		}
	}

	return nil
}

func NewMessageSender(telegram *tg.Client) MessageSenderInterface {
	return &MessageSender{
		telegram: telegram,
	}
}
