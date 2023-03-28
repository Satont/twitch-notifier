package middlewares

import (
	"context"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/types"
	"go.uber.org/zap"
	"time"
)

type LoggMiddleware struct {
	Services *types.Services
}

func (c *LoggMiddleware) Wrap(next tgb.Handler) tgb.Handler {
	return tgb.HandlerFunc(func(ctx context.Context, update *tgb.Update) error {
		defer func(started time.Time) {
			zap.L().
				Info("update handled", zap.Duration("duration", time.Since(started)))
		}(time.Now())

		return next.Handle(ctx, update)
	})
}
