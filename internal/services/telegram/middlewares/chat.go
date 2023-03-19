package middlewares

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/satont/twitch-notifier/internal/services/telegram/types"
	"go.uber.org/zap"
)

type ChatMiddleware struct {
	*tg_types.MiddlewareOpts
}

func (c *ChatMiddleware) Wrap(next tgb.Handler) tgb.Handler {
	return tgb.HandlerFunc(func(ctx context.Context, update *tgb.Update) error {
		chatId := fmt.Sprintf("%v", update.Chat().ID)
		user, err := c.Services.Chat.GetByID(ctx, chatId, db_models.ChatServiceTelegram)
		if err != nil {
			zap.L().Error("failed to get chat", zap.Error(err))
			return nil
		}

		if user == nil {
			user, err = c.Services.Chat.Create(ctx, chatId, db_models.ChatServiceTelegram)
			if err != nil {
				zap.L().Error("failed to create chat", zap.Error(err))
				return nil
			}
		}

		c.SessionManager.Get(ctx).Chat = user

		return next.Handle(ctx, update)
	})
}
