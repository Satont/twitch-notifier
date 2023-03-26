package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/mock"
)

type DbStreamMock struct {
	mock.Mock
}

func (s *DbStreamMock) GetByID(ctx context.Context, streamId string) (*db_models.Stream, error) {
	args := s.Called(ctx, streamId)

	return args.Get(0).(*db_models.Stream), args.Error(1)
}

func (s *DbStreamMock) GetLatestByChannelID(ctx context.Context, channelEntityID uuid.UUID) (*db_models.Stream, error) {
	args := s.Called(ctx, channelEntityID)

	return args.Get(0).(*db_models.Stream), args.Error(1)
}

func (s *DbStreamMock) GetManyByChannelID(ctx context.Context, channelEntityID uuid.UUID, limit int) ([]*db_models.Stream, error) {
	args := s.Called(ctx, channelEntityID, limit)

	return args.Get(0).([]*db_models.Stream), args.Error(1)
}

func (s *DbStreamMock) UpdateOneByStreamID(
	ctx context.Context,
	streamID string,
	updateQuery *db.StreamUpdateQuery,
) (*db_models.Stream, error) {
	args := s.Called(ctx, streamID, updateQuery)

	return args.Get(0).(*db_models.Stream), args.Error(1)
}

func (s *DbStreamMock) CreateOneByChannelID(
	ctx context.Context,
	channelEntityID uuid.UUID,
	updateQuery *db.StreamUpdateQuery,
) (*db_models.Stream, error) {
	args := s.Called(ctx, channelEntityID, updateQuery)

	return args.Get(0).(*db_models.Stream), args.Error(1)
}
