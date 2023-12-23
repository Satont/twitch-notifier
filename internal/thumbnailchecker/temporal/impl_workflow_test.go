package temporal

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	t.Parallel()

	activity := &Activity{}
	workflow := &Workflow{
		activity: activity,
	}

	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activity implementation
	env.OnActivity(
		activity.ThumbnailCheckerTemporalActivity,
		mock.Anything,
		"https://twitch.tv/thumbNail",
	).Return(nil)

	env.ExecuteWorkflow(workflow.Workflow, "https://twitch.tv/thumbNail")

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func Test_WorkflowError(t *testing.T) {
	t.Parallel()

	activity := &Activity{}
	workflow := &Workflow{
		activity: activity,
	}

	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activity implementation
	env.
		OnActivity(
			activity.ThumbnailCheckerTemporalActivity,
			mock.Anything,
			"https://twitch.tv/thumbNail",
		).
		Times(activityMaximumAttempts).
		Return(ErrInvalidThumbnail)

	env.ExecuteWorkflow(workflow.Workflow, "https://twitch.tv/thumbNail")

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}
