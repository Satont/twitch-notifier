package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/mock"
)

type FollowMock struct {
	mock.Mock
}

func (f *FollowMock) Create(
	ctx context.Context,
	channelID uuid.UUID,
	chatID uuid.UUID,
) (*db_models.Follow, error) {
	args := f.Called(ctx, channelID, chatID)

	return args.Get(0).(*db_models.Follow), args.Error(1)
}

func (f *FollowMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := f.Called(ctx, id)

	return args.Error(0)
}

func (f *FollowMock) GetByChatAndChannel(
	ctx context.Context,
	channelId uuid.UUID,
	chatId uuid.UUID,
) (*db_models.Follow, error) {
	args := f.Called(ctx, channelId, chatId)

	return args.Get(0).(*db_models.Follow), args.Error(1)
}

func (f *FollowMock) GetByChannelID(ctx context.Context, channelId uuid.UUID) ([]*db_models.Follow, error) {
	args := f.Called(ctx, channelId)

	return args.Get(0).([]*db_models.Follow), args.Error(1)
}
