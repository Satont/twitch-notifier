package message_sender

import (
	"context"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
)

type MessageOpts struct {
	Text     string
	ImageURL string
}

type MessageSenderInterface interface {
	SendMessage(ctx context.Context, chat *db_models.Chat, opts *MessageOpts) error
}
