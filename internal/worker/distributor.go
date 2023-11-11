package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type TaskDistributor interface {
	DistributeSendPrivateMessage(
		ctx context.Context,
		payload *TaskSendPrivateMessagePayload,
		opts ...asynq.Option,
	) error
}

type redisTaskDistributor struct {
	client *asynq.Client
	logger *zap.Logger
}

func NewRedisTaskDistributor(
	redisOpt asynq.RedisClientOpt,
	logger *zap.Logger,
) *redisTaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &redisTaskDistributor{
		client: client,
		logger: logger,
	}
}
