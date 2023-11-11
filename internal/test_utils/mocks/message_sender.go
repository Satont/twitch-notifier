package mocks

import (
	"context"

	"github.com/satont/twitch-notifier/internal/message_sender"

	"github.com/stretchr/testify/mock"
)

type MessageSenderMock struct {
	mock.Mock
}

func (m *MessageSenderMock) SendMessage(
	ctx context.Context,
	opts *message_sender.MessageOpts,
) error {
	args := m.Called(ctx, opts)

	return args.Error(0)
}
