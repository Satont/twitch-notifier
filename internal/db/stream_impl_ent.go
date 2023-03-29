package db

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/stream"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"time"
)

type StreamEntService struct {
	entClient *ent.Client
}

func (s *StreamEntService) convertEntity(stream *ent.Stream) *db_models.Stream {
	return &db_models.Stream{
		ID:         stream.ID,
		ChannelID:  stream.ChannelID,
		Titles:     stream.Titles,
		Categories: stream.Categories,
		StartedAt:  stream.StartedAt,
		UpdatedAt:  stream.UpdatedAt,
		EndedAt:    stream.EndedAt,
	}
}

func (s *StreamEntService) GetByID(ctx context.Context, streamID string) (*db_models.Stream, error) {
	str, err := s.entClient.Stream.Query().Where(stream.IDEQ(streamID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return s.convertEntity(str), nil
}

func (s *StreamEntService) GetLatestByChannelID(
	ctx context.Context,
	channelEntityID uuid.UUID,
) (*db_models.Stream, error) {
	str, err := s.entClient.Stream.
		Query().
		Where(stream.ChannelIDEQ(channelEntityID), stream.EndedAtIsNil()).
		Order(ent.Desc(stream.FieldStartedAt)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return s.convertEntity(str), nil
}

func (s *StreamEntService) GetManyByChannelID(
	ctx context.Context,
	channelEntityID uuid.UUID,
	limit int,
) ([]*db_models.Stream, error) {
	streams, err := s.entClient.Stream.
		Query().
		Where(stream.HasChannelWith(channel.IDEQ(channelEntityID))).
		Order(ent.Desc(stream.FieldStartedAt)).
		Limit(limit).
		All(ctx)

	if err != nil {
		return nil, err
	}

	convertedStreams := make([]*db_models.Stream, len(streams))
	for i, str := range streams {
		convertedStreams[i] = s.convertEntity(str)
	}

	return convertedStreams, err
}

func (s *StreamEntService) UpdateOneByStreamID(
	ctx context.Context,
	streamID string,
	updateQuery *StreamUpdateQuery,
) (*db_models.Stream, error) {
	str, err := s.GetByID(ctx, streamID)
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, errors.New("stream not found")
	}

	query := s.entClient.Stream.UpdateOneID(str.ID)

	if updateQuery.IsLive != nil && *updateQuery.IsLive {
		query.SetStartedAt(time.Now().UTC())
	}

	if updateQuery.IsLive != nil && !*updateQuery.IsLive {
		query.SetEndedAt(time.Now().UTC())
	}

	if updateQuery.Category != nil {
		str.Categories = append(str.Categories, *updateQuery.Category)
		query.SetCategories(str.Categories)
	}

	if updateQuery.Title != nil {
		str.Titles = append(str.Titles, *updateQuery.Title)
		query.SetTitles(str.Titles)
	}

	newStream, err := query.Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.convertEntity(newStream), nil
}

func (s *StreamEntService) CreateOneByChannelID(
	ctx context.Context,
	channelEntityID uuid.UUID,
	data *StreamUpdateQuery,
) (*db_models.Stream, error) {
	query := s.entClient.Stream.Create()

	query.SetChannelID(channelEntityID)

	query.SetStartedAt(time.Now().UTC())
	query.SetID(data.StreamID)

	if data.Title != nil {
		query.SetTitles(pq.StringArray{*data.Title})
	}

	if data.Category != nil {
		query.SetCategories(pq.StringArray{*data.Category})
	}

	str, err := query.Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.convertEntity(str), nil
}

func NewStreamEntService(entClient *ent.Client) *StreamEntService {
	return &StreamEntService{
		entClient: entClient,
	}
}
