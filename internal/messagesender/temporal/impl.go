package temporal

import (
	"context"

	"github.com/satont/twitch-notifier/internal/messagesender"
)

func NewTemporal() *Temporal {
	return &Temporal{}
}

type Temporal struct {
}

var _ messagesender.MessageSender = (*Temporal)(nil)

func (m *Temporal) SendMessageTelegram(ctx context.Context, opts messagesender.TelegramOpts) error {
	return nil
}
