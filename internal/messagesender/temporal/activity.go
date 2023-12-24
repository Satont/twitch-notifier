package temporal

import (
	"context"

	"github.com/satont/twitch-notifier/internal/messagesender"
	"go.uber.org/fx"
)

type ActivityOpts struct {
	fx.In
}

func NewActivity() *Activity {
	return &Activity{}
}

type Activity struct {
}

func (c *Activity) SendTelegram(ctx context.Context, opts messagesender.TelegramOpts) error {
	return nil
}
