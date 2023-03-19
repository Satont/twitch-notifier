package telegram

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/satont/twitch-notifier/internal/services/telegram/commands"
	"github.com/satont/twitch-notifier/internal/services/telegram/middlewares"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/types"
	"go.uber.org/zap"
)

type telegramService struct {
	services *types.Services
	poller   *tgb.Poller
}

func NewTelegram(ctx context.Context, token string, services *types.Services) *telegramService {
	client := tg.New(token)

	var sessionManager = session.NewManager(tg_types.Session{})

	router := tgb.NewRouter().
		Use(sessionManager).
		//Use(&middlewares.LoggMiddleware{
		//	Services: services,
		//}).
		Use(&middlewares.ChatMiddleware{
			MiddlewareOpts: &tg_types.MiddlewareOpts{
				Services:       services,
				SessionManager: sessionManager,
			}})

	commandOpts := &tg_types.CommandOpts{
		Services:       services,
		Router:         router,
		SessionManager: sessionManager,
	}

	router.Message(func(ctx context.Context, update *tgb.MessageUpdate) error {
		sessionManager.Get(ctx).Scene = ""
		return nil
	}, tgb.Command("cancel"))

	commands.NewStartCommand(commandOpts)
	commands.NewFollowCommand(commandOpts)

	poller := tgb.NewPoller(router, client)

	me, err := client.GetMe().Do(ctx)
	if err != nil {
		zap.S().Fatalw("failed to get bot info", "err", err)
	}

	zap.S().Infow("Telegram bot started", "id", me.ID, "username", me.Username)

	return &telegramService{
		poller:   poller,
		services: services,
	}
}

func (t *telegramService) StartPolling(ctx context.Context) {
	go t.poller.Run(ctx)
}
