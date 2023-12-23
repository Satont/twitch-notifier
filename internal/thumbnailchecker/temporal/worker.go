package temporal

import (
	"context"

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
}

func NewWorker(opts WorkerOpts) error {
	// The client and worker are heavyweight objects that should be created once per process.
	temporalClient, err := client.Dial(
		client.Options{
			Logger: log.NewStructuredLogger(opts.Logger.GetSlog()),
		},
	)
	if err != nil {
		return err
	}

	w := worker.New(temporalClient, queueName, worker.Options{})

	w.RegisterWorkflow(opts.Workflow)
	w.RegisterActivity(opts.Activity)

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
