package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"github.com/satont/twitch-notifier/internal/queue"
	"time"
)

type QueueJobUpdateOpts struct {
	Retries    *int
	FailReason *string
	Status     queue.JobStatus
}

type QueueJobCreateOpts struct {
	QueueName  string
	Data       []byte
	MaxRetries int
	TTL        time.Duration
}

type QueueJobInterface interface {
	AddJob(ctx context.Context, job *QueueJobCreateOpts) (*db_models.QueueJob, error)
	RemoveJobById(ctx context.Context, id uuid.UUID) error
	GetUnprocessedJobsByQueueName(ctx context.Context, queueName string) ([]db_models.QueueJob, error)
	UpdateJob(ctx context.Context, id uuid.UUID, data *QueueJobUpdateOpts) error
}
