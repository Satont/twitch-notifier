package stream

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/repository"
	"gopkg.in/guregu/null.v4"
)

func NewPgx(pg *pgxpool.Pool) *Pgx {
	return &Pgx{
		pg: pg,
	}
}

var _ Repository = (*Pgx)(nil)

type Pgx struct {
	pg *pgxpool.Pool
}

func (c *Pgx) GetById(ctx context.Context, id uuid.UUID) (*domain.Stream, error) {
	query, args, err := repository.Sq.
		Select(
			"id",
			"channel_id",
			"titles",
			"categories",
			"started_at",
			"updated_at",
			"ended_at",
		).
		From("streams").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	stream := Stream{}
	err = c.pg.QueryRow(ctx, query, args...).Scan(
		&stream.ID,
		&stream.ChannelID,
		&stream.Titles,
		&stream.Categories,
		&stream.StartedAt,
		&stream.UpdatedAt,
		&stream.EndedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &domain.Stream{
		ID:         stream.ID,
		ChannelID:  stream.ChannelID,
		Titles:     stream.Titles,
		Categories: stream.Categories,
		StartedAt:  stream.StartedAt,
		UpdatedAt:  stream.UpdatedAt,
		EndedAt:    stream.EndedAt.Ptr(),
	}, err
}

func (c *Pgx) GetLatestByChannelId(ctx context.Context, channelId uuid.UUID) (
	*domain.Stream,
	error,
) {
	query, args, err := repository.Sq.
		Select(
			"id",
			"channel_id",
			"titles",
			"categories",
			"started_at",
			"updated_at",
			"ended_at",
		).
		From("streams").
		Where("channel_id = ?", channelId).
		OrderBy("started_at DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	stream := Stream{}
	err = c.pg.QueryRow(ctx, query, args...).Scan(
		&stream.ID,
		&stream.ChannelID,
		&stream.Titles,
		&stream.Categories,
		&stream.StartedAt,
		&stream.UpdatedAt,
		&stream.EndedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &domain.Stream{
		ID:         stream.ID,
		ChannelID:  stream.ChannelID,
		Titles:     stream.Titles,
		Categories: stream.Categories,
		StartedAt:  stream.StartedAt,
		UpdatedAt:  stream.UpdatedAt,
		EndedAt:    stream.EndedAt.Ptr(),
	}, err
}

func (c *Pgx) GetByChannelId(ctx context.Context, channelId uuid.UUID) ([]domain.Stream, error) {
	query, args, err := repository.Sq.
		Select(
			"id",
			"channel_id",
			"titles",
			"categories",
			"started_at",
			"updated_at",
			"ended_at",
		).
		From("streams").
		Where("channel_id = ?", channelId).
		OrderBy("started_at DESC").
		ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	rows, err := c.pg.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var streams []Stream
	for rows.Next() {
		stream := Stream{}
		err = rows.Scan(
			&stream.ID,
			&stream.ChannelID,
			&stream.Titles,
			&stream.Categories,
			&stream.StartedAt,
			&stream.UpdatedAt,
			&stream.EndedAt,
		)
		if err != nil {
			return nil, err
		}
		streams = append(streams, stream)
	}

	domainStreams := make([]domain.Stream, len(streams))
	for i, stream := range streams {
		domainStreams[i] = domain.Stream{
			ID:         stream.ID,
			ChannelID:  stream.ChannelID,
			Titles:     stream.Titles,
			Categories: stream.Categories,
			StartedAt:  stream.StartedAt,
			UpdatedAt:  stream.UpdatedAt,
			EndedAt:    stream.EndedAt.Ptr(),
		}
	}

	return domainStreams, err
}

func (c *Pgx) Create(ctx context.Context, stream domain.Stream) error {
	query, args, err := repository.Sq.
		Insert("streams").
		Columns(
			"id",
			"channel_id",
			"titles",
			"categories",
			"started_at",
			"updated_at",
			"ended_at",
		).
		Values(
			stream.ID,
			stream.ChannelID,
			stream.Titles,
			stream.Categories,
			stream.StartedAt,
			stream.UpdatedAt,
			null.Time{},
		).ToSql()
	if err != nil {
		return repository.ErrBadQuery
	}

	_, err = c.pg.Exec(ctx, query, args...)
	if err != nil {
		return errors.Join(err, ErrCannotCreate)
	}

	return nil
}

func (c *Pgx) Update(ctx context.Context, stream domain.Stream) error {
	query, args, err := repository.Sq.
		Update("streams").
		Set("channel_id", stream.ChannelID).
		Set("titles", stream.Titles).
		Set("categories", stream.Categories).
		Set("started_at", stream.StartedAt).
		Set("updated_at", stream.UpdatedAt).
		Set("ended_at", null.TimeFromPtr(stream.EndedAt)).
		Where("id = ?", stream.ID).
		ToSql()
	if err != nil {
		return repository.ErrBadQuery
	}

	_, err = c.pg.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := repository.Sq.
		Delete("streams").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return repository.ErrBadQuery
	}

	rows, err := c.pg.Exec(ctx, query, args...)
	if err != nil {
		return errors.Join(err, ErrCannotDelete)
	}

	if rows.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
