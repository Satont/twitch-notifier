package messagesender

import (
	"context"
)

type MessageSender interface {
	SendMessageTelegram(ctx context.Context, opts Opts) error
	// For add new service we need to implement new method, for example:
	// SendMessageDiscord(ctx context.Context, opts Opts) error
}

type Opts struct {
	Target   MessageTarget
	Text     string
	ImageURL string
}

type MessageTarget struct {
	ServiceChatID string
}
