package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
	"github.com/stretchr/testify/mock"
)

type FollowMock struct {
	mock.Mock
}

func (f FollowMock) Create(ctx context.Context, channelID uuid.UUID, chatID uuid.UUID) (*ent.Follow, error) {
	args := f.Called(ctx, channelID, chatID)

	return args.Get(0).(*ent.Follow), args.Error(1)
}

func (f FollowMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := f.Called(ctx, id)

	return args.Error(0)
}
