package temporal

import (
	"github.com/satont/twitch-notifier/internal/messagesender"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
)

type WorkflowOpts struct {
	fx.In
}

func NewWorkflow(opts WorkflowOpts) *Workflow {
	return &Workflow{}
}

type Workflow struct {
}

func (c *Workflow) SendTelegram(ctx workflow.Context, opts messagesender.TelegramOpts)
