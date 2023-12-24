package temporal

import (
	"github.com/satont/twitch-notifier/internal/announcesender"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewActivities,
		NewWorkflow,
		fx.Annotate(NewImpl, fx.As(new(announcesender.AnnounceSender))),
	),
	fx.Invoke(NewWorker),
)
