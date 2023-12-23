package fx

import (
	"github.com/satont/twitch-notifier/internal/thumbnailchecker"
	"github.com/satont/twitch-notifier/internal/thumbnailchecker/temporal"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		temporal.NewActivity,
		temporal.NewWorkflow,
		fx.Annotate(temporal.NewImpl, fx.As(new(thumbnailchecker.ThumbnailChecker))),
	),
	fx.Invoke(temporal.NewWorker),
)
