package temporal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/messagesender"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.uber.org/fx"
)

type TemporalOpts struct {
	fx.In

	Workflow *Workflow
	Logger   logger.Logger
}

func NewImpl(opts TemporalOpts) (*Temporal, error) {
	cl, err := client.Dial(client.Options{})
	if err != nil {
		return nil, err
	}

	return &Temporal{
		client:   cl,
		workflow: opts.Workflow,
		logger:   opts.Logger,
	}, nil
}

const queueName = "message-sender"

type Temporal struct {
	client   client.Client
	workflow *Workflow
	logger   logger.Logger
}

var _ messagesender.MessageSender = (*Temporal)(nil)

func (c *Temporal) SendMessageTelegram(ctx context.Context, opts messagesender.TelegramOpts) error {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("MSG: Telegram to %s #%s", uuid.NewString(), opts.ServiceChatID),
		TaskQueue: queueName,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 5,
			InitialInterval: 10 * time.Second,
		},
	}

	we, err := c.client.ExecuteWorkflow(ctx, workflowOptions, c.workflow.SendTelegram, opts)
	if err != nil {
		return err
	}

	err = we.Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
