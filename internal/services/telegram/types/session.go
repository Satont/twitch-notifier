package tg_types

import (
	"context"

	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/satont/twitch-notifier/internal/services/types"
)

type SessionManager[T comparable] interface {
	SetEqualFunc(fn func(t T, t2 T) bool)
	Setup(opt session.ManagerOption, opts ...session.ManagerOption)
	Get(ctx context.Context) *T
	Reset(session *T)
	Filter(fn func(t *T) bool) tgb.Filter
	Wrap(next tgb.Handler) tgb.Handler
}

type Menu struct {
	CurrentPage int
	TotalPages  int
}

type Session struct {
	Chat  *db_models.Chat
	Scene string

	FollowsMenu *Menu
}

type CommandOpts struct {
	Services       *types.Services
	Router         Router
	SessionManager SessionManager[Session]
}

type MiddlewareOpts struct {
	Services       *types.Services
	SessionManager SessionManager[Session]
}
