package db_models

import (
	"github.com/google/uuid"
	"time"
)

type QueueJob struct {
	ID         uuid.UUID
	QueueName  string
	Data       []byte
	MaxRetries int
	Retries    int
	AddedAt    time.Time
	TTL        time.Duration
	FailReason string
}
