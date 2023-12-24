package temporal

import (
	"time"

	"github.com/satont/twitch-notifier/internal/messagesender"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
)

type WorkflowOpts struct {
	fx.In

	Logger   logger.Logger
	Activity *Activity
}

func NewWorkflow(opts WorkflowOpts) *Workflow {
	return &Workflow{
		logger:   opts.Logger,
		activity: opts.Activity,
	}
}

type Workflow struct {
	logger   logger.Logger
	activity *Activity
}

const telegramActivityMaximumAttempts = 5

func (c *Workflow) SendTelegram(ctx workflow.Context, opts messagesender.TelegramOpts) error {
	ao := workflow.ActivityOptions{
		TaskQueue:           queueName,
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: telegramActivityMaximumAttempts,
			InitialInterval: 2 * time.Second,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	log := workflow.GetLogger(ctx)
	log.Info("Sending message", "chatId", opts.ServiceChatID)

	err := workflow.ExecuteActivity(
		ctx,
		c.activity.SendTelegram,
		opts,
	).Get(
		ctx,
		nil,
	)
	if err != nil {
		log.Error("Send failed", "Error", err)
		return err
	}

	log.Info("Message sent")

	return nil
}
