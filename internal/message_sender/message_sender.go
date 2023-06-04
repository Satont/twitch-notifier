package message_sender

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/satont/twitch-notifier/internal/db/db_models"
)

type MessageOpts struct {
	Chat      *db_models.Chat
	Text      string
	ImageURL  string
	ParseMode *tg.ParseMode
}

type MessageSenderInterface interface {
	SendMessage(ctx context.Context, opts *MessageOpts) error
}
