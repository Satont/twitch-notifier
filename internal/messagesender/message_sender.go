package messagesender

import (
	"context"
)

type MessageSender interface {
	SendMessageTelegram(ctx context.Context, opts TelegramOpts) error
	// For add new service we need to implement new method, for example:
	// SendMessageDiscord(ctx context.Context, opts TelegramOpts) error
}

type TelegramOpts struct {
	ServiceChatID string
	Text          string
	ImageURL      string
}