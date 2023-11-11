package mocks

import (
	context "context"

	asynq "github.com/hibiken/asynq"
	worker "github.com/satont/twitch-notifier/internal/worker"
	"github.com/stretchr/testify/mock"
)

type TaskDistributorMock struct {
	mock.Mock
}

// DistributeSendPrivateMessage mocks base method.
func (m *TaskDistributorMock) DistributeSendPrivateMessage(
	ctx context.Context,
	payload *worker.TaskSendPrivateMessagePayload,
	opts ...asynq.Option,
) error {
	args := m.Called(ctx, payload, opts)

	return args.Error(0)
}
