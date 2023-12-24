package temporal

import (
	"github.com/satont/twitch-notifier/internal/thumbnailchecker"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewActivity,
		NewWorkflow,
		fx.Annotate(NewImpl, fx.As(new(thumbnailchecker.ThumbnailChecker))),
	),
	fx.Invoke(NewWorker),
)
