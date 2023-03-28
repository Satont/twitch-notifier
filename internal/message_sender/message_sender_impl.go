package message_sender

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"strconv"
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

		if opts.ImageURL != "" {
			query := m.telegram.
				SendPhoto(tg.ChatID(chatId), tg.FileArg{URL: opts.ImageURL}).
				Caption(opts.Text)

			if opts.ParseMode != nil {
				query = query.ParseMode(*opts.ParseMode)
			}

			return query.DoVoid(ctx)
		} else {
			query := m.telegram.
				SendMessage(tg.ChatID(chatId), opts.Text).
				DisableWebPagePreview(true)

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
