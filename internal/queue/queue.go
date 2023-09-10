package queue

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/satont/workers-pool"
)

type JobStatus string

const (
	JobStatusPending JobStatus = "pending"
	JobStatusFailed  JobStatus = "failed"
	JobStatusDone    JobStatus = "done"
)

func (c JobStatus) String() string {
	return string(c)
}

func (JobStatus) Values() []string {
	return []string{JobStatusPending.String(), JobStatusFailed.String(), JobStatusDone.String()}
}

type Job[T any] struct {
	ID         uuid.UUID
	Arguments  T
	CreatedAt  time.Time
	TTL        time.Duration
	MaxRetries int
	Status     JobStatus

	Retries int
}

type UpdateData struct {
	JobID      uuid.UUID
	Retries    int
	FailReason string

	Status JobStatus
}
type AnyFunc[T any] func(ctx context.Context, args T) error
type UpdateHook func(ctx context.Context, data *UpdateData)

type Queue[T any] struct {
	channel     chan *Job[T]
	workersPool *gopool.Pool
	run         AnyFunc[T]
	updateHook  UpdateHook
}
type Opts[T any] struct {
	Run        AnyFunc[T]
	PoolSize   int
	UpdateHook UpdateHook
}

func New[T any](opts Opts[T]) *Queue[T] {
	if opts.PoolSize == 0 {
		opts.PoolSize = 1
	}

	q := &Queue[T]{
		channel:     make(chan *Job[T]),
		workersPool: gopool.NewPool(opts.PoolSize),
		run:         opts.Run,
		updateHook:  opts.UpdateHook,
	}

	return q
}

func (q *Queue[T]) Push(item *Job[T]) {
	if item.MaxRetries == 0 {
		item.MaxRetries = 3
	}

	q.channel <- item
}

func (q *Queue[T]) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case job := <-q.channel:
				q.workersPool.Submit(func() {
					q.process(ctx, job)
				})
			}
		}
	}()
}

func (q *Queue[T]) process(ctx context.Context, job *Job[T]) {
	if job.TTL != 0 && time.Now().After(job.CreatedAt.Add(job.TTL)) {
		if q.updateHook != nil {
			q.updateHook(ctx, &UpdateData{
				JobID:   job.ID,
				Retries: job.Retries,
				Status:  JobStatusDone,
			})
		}

		return
	}

	err := q.run(ctx, job.Arguments)

	var failReason string
	if err != nil {
		failReason = err.Error()
	}

	doRetry := job.Retries <= job.MaxRetries
	if err != nil && doRetry {
		job.Retries++
		q.Push(job)
	}

	if q.updateHook != nil {
		var status JobStatus
		if doRetry {
			status = JobStatusPending
		} else if err != nil {
			status = JobStatusFailed
		} else {
			status = JobStatusDone
		}

		q.updateHook(ctx, &UpdateData{
			JobID:      job.ID,
			Retries:    job.Retries,
			FailReason: failReason,
			Status:     status,
		})
	}
}
