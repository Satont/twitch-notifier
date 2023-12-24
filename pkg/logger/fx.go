package logger

import (
	"github.com/satont/twitch-notifier/pkg/config"
	"go.uber.org/fx"
)

var FxOption = fx.Annotate(
	func(cfg config.Config) *Impl {
		return New(
			Opts{
				Env: cfg.AppEnv,
			},
		)
	},
	fx.As(new(Logger)),
)
