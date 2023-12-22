package chat_settings

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(
		NewPgx,
		fx.As(new(Repository)),
	),
)
