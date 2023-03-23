package message_sender

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
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
			err := m.telegram.
				SendPhoto(tg.ChatID(chatId), tg.FileArg{URL: opts.ImageURL}).
				Caption(opts.Text).
				DoVoid(ctx)
			return err
		} else {
			err := m.telegram.SendMessage(tg.ChatID(chatId), opts.Text).DoVoid(ctx)
			return err
		}
	}

	return nil
}

func NewMessageSender(telegram *tg.Client) MessageSenderInterface {
	return &MessageSender{
		telegram: telegram,
	}
}
