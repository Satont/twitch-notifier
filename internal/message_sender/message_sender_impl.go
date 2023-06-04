package message_sender

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/mr-linch/go-tg"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"github.com/satont/twitch-notifier/internal/queue"
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
		MaxRetries: lo.ToPtr(3),
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
	return &MessageSender{
		telegram: telegram,
		que: queue.New[*MessageOpts](queue.Opts[*MessageOpts]{
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
			UpdateHook: func(data *queue.UpdateData) {
				dbQueue.UpdateJob(context.Background(), data.JobID, &db.QueueJobUpdateOpts{
					Retries:    &data.Retries,
					FailReason: &data.FailReason,
				})
			},
		}),
		dbQueue: dbQueue,
	}
}
