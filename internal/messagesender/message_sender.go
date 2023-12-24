package messagesender

import (
	"context"
)

//go:generate go run go.uber.org/mock/mockgen -source=message_sender.go -destination=mocks/mock.go

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
