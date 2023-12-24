package temporal

import (
	"context"

	"github.com/satont/twitch-notifier/pkg/config"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

type WorkerOpts struct {
	fx.In
	LC fx.Lifecycle

	Workflow *Workflow
	Activity *Activity
	Logger   logger.Logger
	Config   config.Config
}

func NewWorker(opts WorkerOpts) error {
	temporalClient, err := client.Dial(
		client.Options{
			Logger:   log.NewStructuredLogger(opts.Logger.GetSlog()),
			HostPort: opts.Config.TemporalUrl,
		},
	)
	if err != nil {
		return err
	}

	w := worker.New(temporalClient, queueName, worker.Options{})

	w.RegisterWorkflow(opts.Workflow.SendTelegram)
	w.RegisterActivity(opts.Activity.SendTelegram)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return w.Start()
			},
			OnStop: func(ctx context.Context) error {
				w.Stop()
				temporalClient.Close()
				return nil
			},
		},
	)

	return nil
}
