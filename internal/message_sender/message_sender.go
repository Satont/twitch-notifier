package message_sender

import (
	"context"
	"github.com/satont/twitch-notifier/internal/db/db_models"
)

type TgParseMode string

const (
	TgParseModeMD TgParseMode = "markdown"
)

type MessageOpts struct {
	Chat     *db_models.Chat
	Text     string
	ImageURL string

	TgParseMode TgParseMode
}

type MessageSenderInterface interface {
	SendMessage(ctx context.Context, opts *MessageOpts) error
}
