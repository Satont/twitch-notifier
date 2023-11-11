package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/satont/twitch-notifier/internal/message_sender"
	"go.uber.org/zap"
)

func (processor *RedisTaskProcessor) ProcessTaskSendPrivateMessage(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload TaskSendPrivateMessagePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	err := processor.messageSender.SendMessage(ctx, &message_sender.MessageOpts{
		ChatID:      payload.ChatID,
		ChatService: payload.ChatService,
		Text:        payload.Text,
		ImageURL:    payload.ImageURL,
		TgParseMode: payload.TgParseMode,
		Buttons:     payload.Buttons,
		SkipButtons: payload.SkipButtons,
	})
	if err != nil {
		return err
	}

	return nil
}

func (distributor *redisTaskDistributor) DistributeSendPrivateMessage(
	ctx context.Context,
	payload *TaskSendPrivateMessagePayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSendPrivateMessage, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}

	distributor.logger.Info("task sent", zap.String("id", info.ID))

	return nil
}
