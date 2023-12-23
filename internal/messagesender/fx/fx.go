package fx

import (
	"github.com/satont/twitch-notifier/internal/messagesender"
	"github.com/satont/twitch-notifier/internal/messagesender/temporal"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		temporal.NewActivity,
		temporal.NewWorkflow,
		fx.Annotate(temporal.NewImpl, fx.As(new(messagesender.MessageSender))),
	),
	fx.Invoke(temporal.NewWorker),
)
