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
)

type telegramService struct {
	services *types.Services
	poller   *tgb.Poller
}

func NewTelegram(token string, services *types.Services) *telegramService {
	client := tg.New(token)

	var sessionManager = session.NewManager(tg_types.Session{})

	router := tgb.NewRouter().
		Use(sessionManager).
		Use(&middlewares.LoggMiddleware{
			Services: services,
		}).
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

	commands.NewStartCommand(commandOpts)
	commands.NewFollowCommand(commandOpts)

	poller := tgb.NewPoller(router, client)

	return &telegramService{
		poller:   poller,
		services: services,
	}
}

func (t *telegramService) StartPolling(ctx context.Context) {
	go t.poller.Run(ctx)
}
