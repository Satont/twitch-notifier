package middlewares

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/satont/twitch-notifier/internal/services/types"
	"go.uber.org/zap"
)

type ChatMiddleware struct {
	Services *types.Services
}

func (c *ChatMiddleware) Wrap(next tgb.Handler) tgb.Handler {
	return tgb.HandlerFunc(func(ctx context.Context, update *tgb.Update) error {
		chatId := fmt.Sprintf("%v", update.Chat().ID)
		user, err := c.Services.Chat.GetByID(chatId, chat.ServiceTelegram)
		if err != nil {
			zap.L().Error("failed to get chat", zap.Error(err))
			return nil
		}

		if user == nil {
			user, err = c.Services.Chat.Create(chatId, chat.ServiceTelegram)
			if err != nil {
				zap.L().Error("failed to create chat", zap.Error(err))
				return nil
			}
		}

		zap.L().Info("chat", zap.Any("user", user))

		return next.Handle(ctx, update)
	})
}
