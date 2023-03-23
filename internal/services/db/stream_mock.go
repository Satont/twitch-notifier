package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/mock"
)

type StreamMock struct {
	mock.Mock
}

func (s *StreamMock) GetByID(ctx context.Context, streamId string) (*db_models.Stream, error) {
	args := s.Called(ctx, streamId)

	return args.Get(0).(*db_models.Stream), args.Error(1)
}

func (s *StreamMock) GetLatestByChannelID(ctx context.Context, channelEntityID uuid.UUID) (*db_models.Stream, error) {
	args := s.Called(ctx, channelEntityID)

	return args.Get(0).(*db_models.Stream), args.Error(1)
}

func (s *StreamMock) GetManyByChannelID(ctx context.Context, channelEntityID uuid.UUID, limit int) ([]*db_models.Stream, error) {
	args := s.Called(ctx, channelEntityID, limit)

	return args.Get(0).([]*db_models.Stream), args.Error(1)
}

func (s *StreamMock) UpdateOneByStreamID(ctx context.Context, streamID string, updateQuery *StreamUpdateQuery) (*db_models.Stream, error) {
	args := s.Called(ctx, streamID, updateQuery)

	return args.Get(0).(*db_models.Stream), args.Error(1)
}

func (s *StreamMock) CreateOneByChannelID(ctx context.Context, channelEntityID uuid.UUID, updateQuery *StreamUpdateQuery) (*db_models.Stream, error) {
	args := s.Called(ctx, channelEntityID, updateQuery)

	return args.Get(0).(*db_models.Stream), args.Error(1)
}
