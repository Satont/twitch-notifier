package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twitch-notifier/internal/message_sender"
	"go.uber.org/zap"
)

const (
	QueueDefault = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendPrivateMessage(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server        *asynq.Server
	logger        *zap.Logger
	messageSender message_sender.MessageSenderInterface
}

func NewRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	logger *zap.Logger,
	messageSender message_sender.MessageSenderInterface,
) *RedisTaskProcessor {
	redisLogger := NewLogger(logger)
	redis.SetLogger(redisLogger)

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueDefault: 5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(
				func(ctx context.Context, task *asynq.Task, err error) {
					logger.Sugar().Error(err)
				},
			),
			Logger: redisLogger,
		},
	)

	return &RedisTaskProcessor{
		server:        server,
		logger:        logger,
		messageSender: messageSender,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendPrivateMessage, processor.ProcessTaskSendPrivateMessage)

	return processor.server.Start(mux)
}

const TaskSendPrivateMessage = "task:send_private_message"

type TaskSendPrivateMessagePayload struct {
	ChatID      string
	ChatService string
	Text        string
	ImageURL    string

	TgParseMode message_sender.TgParseMode
	Buttons     [][]message_sender.KeyboardButton
	SkipButtons bool
}
