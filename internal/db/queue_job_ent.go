package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/queuejob"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"time"
)

type QueueJobEntService struct {
	entClient *ent.Client
}

func NewQueueJobEntService(entClient *ent.Client) QueueJobInterface {
	return &QueueJobEntService{
		entClient: entClient,
	}
}

func (q *QueueJobEntService) convertEntity(job *ent.QueueJob) *db_models.QueueJob {
	return &db_models.QueueJob{
		ID:         job.ID,
		QueueName:  job.QueueName,
		Data:       job.Data,
		Retries:    job.Retries,
		AddedAt:    job.AddedAt,
		TTL:        time.Duration(job.TTL) * time.Millisecond,
		FailReason: job.FailReason,
	}
}

func (q *QueueJobEntService) AddJob(ctx context.Context, job *QueueJobCreateOpts) (*db_models.QueueJob, error) {
	j, err := q.entClient.QueueJob.Create().
		SetID(uuid.New()).
		SetQueueName(job.QueueName).
		SetData(job.Data).
		SetRetries(0).
		SetTTL(int(job.TTL.Milliseconds())).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return q.convertEntity(j), nil
}

func (q *QueueJobEntService) RemoveJobById(ctx context.Context, id uuid.UUID) error {
	return q.entClient.QueueJob.DeleteOneID(id).Exec(ctx)
}

func (q *QueueJobEntService) GetJobsByQueueName(ctx context.Context, queueName string) ([]db_models.QueueJob, error) {
	jobs, err := q.entClient.QueueJob.Query().Where(queuejob.QueueName(queueName)).All(ctx)
	if err != nil {
		return nil, err
	}

	var convertedJobs []db_models.QueueJob
	for _, job := range jobs {
		convertedJobs = append(convertedJobs, *q.convertEntity(job))
	}

	return convertedJobs, nil
}

func (q *QueueJobEntService) UpdateJob(ctx context.Context, id uuid.UUID, data *QueueJobUpdateOpts) error {
	query := q.entClient.QueueJob.UpdateOneID(id)

	if data.Retries != nil {
		query = query.SetRetries(*data.Retries)
	}

	if data.FailReason != nil {
		query = query.SetFailReason(*data.FailReason)
	}

	_, err := query.Save(ctx)

	return err
}
