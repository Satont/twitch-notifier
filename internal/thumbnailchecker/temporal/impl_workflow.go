package temporal

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WorkflowOpts struct {
	Activity *Activity
}

func NewWorkflow(opts WorkflowOpts) *Workflow {
	return &Workflow{
		activity: opts.Activity,
	}
}

type Workflow struct {
	activity *Activity
}

const activityMaximumAttempts = 50

func (c *Workflow) Workflow(ctx workflow.Context, thumbNailUrl string) error {
	ao := workflow.ActivityOptions{
		TaskQueue:           queueName,
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval:        15 * time.Second,
			MaximumAttempts:        activityMaximumAttempts,
			NonRetryableErrorTypes: nil,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Starting check thumbnail validation", "thumbNailUrl", thumbNailUrl)

	err := workflow.ExecuteActivity(
		ctx,
		c.activity.ThumbnailCheckerTemporalActivity,
		thumbNailUrl,
	).Get(
		ctx,
		nil,
	)
	if err != nil {
		logger.Error("Validation failed", "Error", err)
		return err
	}

	logger.Info("Thumbnail is validated")

	return nil
}
