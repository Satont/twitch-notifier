package tg_types

import (
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/satont/twitch-notifier/internal/services/types"
)

type Session struct {
	PizzaCount int
}

type CommandOpts struct {
	Services       *types.Services
	Router         *tgb.Router
	SessionManager *session.Manager[Session]
}
