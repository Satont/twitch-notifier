package message_sender

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/mr-linch/go-tg"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"github.com/satont/twitch-notifier/internal/queue"
	"go.uber.org/zap"
	"strconv"
)

type MessageSender struct {
	telegram *tg.Client
	que      *queue.Queue[*MessageOpts]
	dbQueue  db.QueueJobInterface
}

func (m *MessageSender) SendMessage(ctx context.Context, opts *MessageOpts) error {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(opts)
	if err != nil {
		return err
	}

	dbJob, err := m.dbQueue.AddJob(ctx, &db.QueueJobCreateOpts{
		QueueName:  "send_message",
		Data:       buff.Bytes(),
		MaxRetries: 3,
	})

	if err != nil {
		return err
	}

	job := &queue.Job[*MessageOpts]{
		ID:         dbJob.ID,
		Arguments:  opts,
		CreatedAt:  dbJob.AddedAt,
		MaxRetries: dbJob.MaxRetries,
	}

	m.que.Push(job)

	return nil
}

func NewMessageSender(telegram *tg.Client, dbQueue db.QueueJobInterface) MessageSenderInterface {
	que := queue.New[*MessageOpts](queue.Opts[*MessageOpts]{
		Run: func(ctx context.Context, opts *MessageOpts) error {
			var parseMode tg.ParseMode

			if opts.TgParseMode == TgParseModeMD {
				parseMode = tg.MD
			}

			if opts.Chat.Service == db_models.ChatServiceTelegram {
				chatId, err := strconv.Atoi(opts.Chat.ChatID)
				if err != nil {
					return err
				}

				if opts.ImageURL != "" {
					query := telegram.
						SendPhoto(tg.ChatID(chatId), tg.FileArg{URL: opts.ImageURL}).
						Caption(opts.Text)

					if opts.TgParseMode != "" {
						query = query.ParseMode(parseMode)
					}

					return query.DoVoid(ctx)
				} else {
					query := telegram.
						SendMessage(tg.ChatID(chatId), opts.Text).
						DisableWebPagePreview(true)

					if opts.TgParseMode != "" {
						query = query.ParseMode(parseMode)
					}

					return query.DoVoid(ctx)
				}
			}

			return nil
		},
		PoolSize: 50,
		UpdateHook: func(ctx context.Context, data *queue.UpdateData) {
			dbQueue.UpdateJob(ctx, data.JobID, &db.QueueJobUpdateOpts{
				Retries:    &data.Retries,
				FailReason: &data.FailReason,
				Status:     data.Status,
			})
		},
	})
	que.Run(context.Background())

	jobs, err := dbQueue.GetUnprocessedJobsByQueueName(context.Background(), "send_message")
	if err != nil {
		zap.S().Error(err)
	} else {
		for _, job := range jobs {
			buff := bytes.NewBuffer(job.Data)
			decoder := gob.NewDecoder(buff)
			opts := &MessageOpts{}
			err = decoder.Decode(opts)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			que.Push(&queue.Job[*MessageOpts]{
				ID:         job.ID,
				Arguments:  opts,
				CreatedAt:  job.AddedAt,
				MaxRetries: job.MaxRetries,
				Retries:    job.Retries,
				TTL:        job.TTL,
				Status:     job.Status,
			})
		}
	}

	return &MessageSender{
		telegram,
		que,
		dbQueue,
	}
}
