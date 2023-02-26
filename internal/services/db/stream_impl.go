package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/stream"
	"time"
)

type streamService struct {
	entClient *ent.Client
}

func (s *streamService) GetByID(streamId string) (*ent.Stream, error) {
	str, err := s.entClient.Stream.Query().Where(stream.IDEQ(streamId)).Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return str, nil
}

func (s *streamService) GetLatestByChannelID(channelEntityID uuid.UUID) (*ent.Stream, error) {
	str, err := s.entClient.Stream.
		Query().
		Where(stream.HasChannelWith(channel.IDEQ(channelEntityID))).
		Order(ent.Desc(stream.FieldStartedAt)).
		First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return str, nil
}

func (s *streamService) GetManyByChannelID(channelEntityID uuid.UUID, limit int) ([]*ent.Stream, error) {
	streams, err := s.entClient.Stream.
		Query().
		Where(stream.HasChannelWith(channel.IDEQ(channelEntityID))).
		Order(ent.Desc(stream.FieldStartedAt)).
		Limit(limit).
		All(context.Background())

	if err != nil {
		return nil, err
	}

	return streams, err
}

func (s *streamService) UpdateOneByStreamID(streamID string, updateQuery *StreamUpdateQuery) (*ent.Stream, error) {
	stream, err := s.GetByID(streamID)
	if err != nil {
		return nil, err
	}
	if stream == nil {
		return nil, nil
	}

	query := s.entClient.Stream.UpdateOne(stream)

	if updateQuery.IsLive != nil && *updateQuery.IsLive {
		query.SetStartedAt(time.Now().UTC())
	}

	if updateQuery.IsLive != nil && !*updateQuery.IsLive {
		query.SetEndedAt(time.Now().UTC())
	}

	if updateQuery.Category != nil {
		query.AppendCategories([]string{*updateQuery.Category})
	}

	if updateQuery.Title != nil {
		query.AppendTitles([]string{*updateQuery.Title})
	}

	return query.Save(context.Background())
}

func (s *streamService) CreateOneByChannelID(channelEntityID uuid.UUID, data *StreamUpdateQuery) (*ent.Stream, error) {
	query := s.entClient.Stream.Create().
		SetChannelID(channelEntityID)

	query.SetStartedAt(time.Now().UTC())

	if data.Title != nil {
		query.SetTitles([]string{*data.Title})
	}

	if data.Category != nil {
		query.SetCategories([]string{*data.Category})
	}

	return query.Save(context.Background())
}

func NewStreamService(entClient *ent.Client) *streamService {
	return &streamService{
		entClient: entClient,
	}
}
