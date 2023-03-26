package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/mock"
)

type DbFollowMock struct {
	mock.Mock
}

func (f *DbFollowMock) Create(
	ctx context.Context,
	channelID uuid.UUID,
	chatID uuid.UUID,
) (*db_models.Follow, error) {
	args := f.Called(ctx, channelID, chatID)

	return args.Get(0).(*db_models.Follow), args.Error(1)
}

func (f *DbFollowMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := f.Called(ctx, id)

	return args.Error(0)
}

func (f *DbFollowMock) GetByChatAndChannel(
	ctx context.Context,
	channelId uuid.UUID,
	chatId uuid.UUID,
) (*db_models.Follow, error) {
	args := f.Called(ctx, channelId, chatId)

	return args.Get(0).(*db_models.Follow), args.Error(1)
}

func (f *DbFollowMock) GetByChannelID(ctx context.Context, channelId uuid.UUID) ([]*db_models.Follow, error) {
	args := f.Called(ctx, channelId)

	return args.Get(0).([]*db_models.Follow), args.Error(1)
}

func (f *DbFollowMock) GetByChatID(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*db_models.Follow, error) {
	args := f.Called(ctx, chatID, limit, offset)

	return args.Get(0).([]*db_models.Follow), args.Error(1)
}

func (f *DbFollowMock) CountByChatID(ctx context.Context, chatID uuid.UUID) (int, error) {
	args := f.Called(ctx, chatID)

	return args.Int(0), args.Error(1)
}
