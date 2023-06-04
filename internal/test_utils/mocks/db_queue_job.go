package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"github.com/stretchr/testify/mock"
)

type DbQueueMock struct {
	mock.Mock
}

func (s *DbQueueMock) AddJob(ctx context.Context, job *db.QueueJobCreateOpts) (*db_models.QueueJob, error) {
	args := s.Called(ctx, job)

	return args.Get(0).(*db_models.QueueJob), args.Error(1)
}

func (s *DbQueueMock) RemoveJobById(ctx context.Context, id uuid.UUID) error {
	args := s.Called(ctx, id)

	return args.Error(0)
}

func (s *DbQueueMock) GetUnprocessedJobsByQueueName(ctx context.Context, queueName string) ([]db_models.QueueJob, error) {
	args := s.Called(ctx, queueName)

	return args.Get(0).([]db_models.QueueJob), args.Error(1)
}

func (s *DbQueueMock) UpdateJob(ctx context.Context, id uuid.UUID, data *db.QueueJobUpdateOpts) error {
	args := s.Called(ctx, id, data)

	return args.Error(0)
}
