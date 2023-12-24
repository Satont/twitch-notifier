package temporal

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/thumbnailchecker"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.uber.org/fx"
)

const queueName = "thumbnail_checker"

type ImplOpts struct {
	fx.In

	Logger   logger.Logger
	Workflow *Workflow
}

func NewImpl(opts ImplOpts) (*Temporal, error) {
	cl, err := client.Dial(
		client.Options{
			Logger: log.NewStructuredLogger(opts.Logger.GetSlog()),
		},
	)
	if err != nil {
		return nil, err
	}

	return &Temporal{
		client:   cl,
		workflow: opts.Workflow,
	}, nil
}

var _ thumbnailchecker.ThumbnailChecker = (*Temporal)(nil)

type Temporal struct {
	client   client.Client
	workflow *Workflow
}

func (c *Temporal) ValidateThumbnail(ctx context.Context, thumbnailUrl string) error {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("Thumnail check %s", uuid.NewString()),
		TaskQueue: queueName,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 100,
		},
	}

	we, err := c.client.ExecuteWorkflow(ctx, workflowOptions, c.workflow, thumbnailUrl)
	if err != nil {
		return err
	}

	err = we.Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Temporal) TransformSizes(url string, width int, height int) string {
	thumbNail := url
	thumbNail = strings.Replace(thumbNail, "{width}", strconv.Itoa(width), 1)
	thumbNail = strings.Replace(thumbNail, "{height}", strconv.Itoa(height), 1)

	return thumbNail
}
