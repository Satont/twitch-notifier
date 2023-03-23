package message_sender

import (
	"context"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) SendMessage(ctx context.Context, chat *db_models.Chat, opts *MessageOpts) error {
	args := m.Called(ctx, chat, opts)

	return args.Error(0)
}
