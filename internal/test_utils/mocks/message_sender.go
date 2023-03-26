package mocks

import (
	"context"

	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/satont/twitch-notifier/internal/services/message_sender"
	"github.com/stretchr/testify/mock"
)

type MessageSenderMock struct {
	mock.Mock
}

func (m *MessageSenderMock) SendMessage(ctx context.Context, chat *db_models.Chat, opts *message_sender.MessageOpts) error {
	args := m.Called(ctx, chat, opts)

	return args.Error(0)
}
