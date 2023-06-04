package queue

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/satont/workers-pool"
)

type Job[T any] struct {
	ID         uuid.UUID
	Arguments  T
	CreatedAt  time.Time
	TTL        time.Duration
	MaxRetries int

	retries int
}

type UpdateData struct {
	JobID      uuid.UUID
	Retries    int
	FailReason string
}
type AnyFunc[T any] func(ctx context.Context, args T) error
type UpdateTook func(data *UpdateData)

type Queue[T any] struct {
	channel     chan *Job[T]
	workersPool *gopool.Pool
	f           AnyFunc[T]
	updateHook  UpdateTook
}
type Opts[T any] struct {
	Run        AnyFunc[T]
	PoolSize   int
	UpdateHook UpdateTook
}

func New[T any](opts Opts[T]) *Queue[T] {
	if opts.PoolSize == 0 {
		opts.PoolSize = 1
	}

	q := &Queue[T]{
		channel:     make(chan *Job[T]),
		workersPool: gopool.NewPool(opts.PoolSize),
		f:           opts.Run,
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
			q.updateHook(&UpdateData{
				JobID:   job.ID,
				Retries: job.retries,
			})
		}

		return
	}

	err := q.f(ctx, job.Arguments)

	var failReason string
	if err != nil {
		failReason = err.Error()
	}

	if err != nil && job.retries <= job.MaxRetries {
		job.retries++
		q.Push(job)
	}

	if q.updateHook != nil {
		q.updateHook(&UpdateData{
			Retries:    job.retries,
			FailReason: failReason,
		})
	}
}
