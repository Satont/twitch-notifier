package logger

import (
	"go.uber.org/fx"
)

// TODO: read from config
func NewFx() func() *Impl {
	return func() *Impl {
		return New(
			Opts{
				Env: "development",
			},
		)
	}
}

var Module = fx.Annotate(
	NewFx,
	fx.As(new(Logger)),
)
