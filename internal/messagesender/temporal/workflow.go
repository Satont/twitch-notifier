package temporal

import (
	"github.com/satont/twitch-notifier/internal/messagesender"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
)

type WorkflowOpts struct {
	fx.In

	Logger logger.Logger
}

func NewWorkflow(opts WorkflowOpts) *Workflow {
	return &Workflow{}
}

type Workflow struct {
}

func (c *Workflow) SendTelegram(ctx workflow.Context, opts messagesender.TelegramOpts) error {
	return nil
}
