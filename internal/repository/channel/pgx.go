package channel

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/repository"
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

const tableName = "channels"

func (c *Pgx) GetById(ctx context.Context, id uuid.UUID) (*domain.Channel, error) {
	channel := Channel{}

	query, args, err := repository.Sq.
		Select("id", "channel_id", "service").
		From(tableName).
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(&channel.ID, &channel.ChannelID, &channel.Service)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &domain.Channel{
		ID:        channel.ID,
		ChannelID: channel.ChannelID,
		Service:   domain.StreamingService(channel.Service),
	}, nil
}

func (c *Pgx) GetByStreamServiceAndID(
	ctx context.Context,
	service StreamingService,
	id string,
) (*domain.Channel, error) {
	query, args, err := repository.Sq.
		Select("id", "channel_id", "service").
		From(tableName).
		Where(
			"channel_id = ? AND service = ?",
			id,
			service,
		).ToSql()
	if err != nil {
		return nil, err
	}

	channel := Channel{}
	err = c.pg.QueryRow(ctx, query, args...).Scan(&channel.ID, &channel.ChannelID, &channel.Service)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &domain.Channel{
		ID:        channel.ID,
		ChannelID: channel.ChannelID,
		Service:   domain.StreamingService(channel.Service),
	}, nil
}

func (c *Pgx) GetAll(ctx context.Context) ([]domain.Channel, error) {
	query, args, err := repository.Sq.
		Select("id", "channel_id", "service").
		From(tableName).
		ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	rows, err := c.pg.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var channels []Channel
	defer rows.Close()
	for rows.Next() {
		channel := Channel{}
		err = rows.Scan(&channel.ID, &channel.ChannelID, &channel.Service)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	resultChannels := make([]domain.Channel, len(channels))
	for i, channel := range channels {
		resultChannels[i] = domain.Channel{
			ID:        channel.ID,
			ChannelID: channel.ChannelID,
			Service:   domain.StreamingService(channel.Service),
		}
	}

	return resultChannels, nil
}

func (c *Pgx) Create(ctx context.Context, channel domain.Channel) error {
	query, args, err := repository.Sq.Insert(tableName).Columns(
		"id",
		"channel_id",
		"service",
	).Values(
		channel.ID,
		channel.ChannelID,
		channel.Service,
	).ToSql()
	if err != nil {
		return repository.ErrBadQuery
	}

	_, err = c.pg.Exec(ctx, query, args...)
	if err != nil {
		// TODO: viyasnit kak luchshe
		return errors.Join(err, ErrCannotCreate)
	}

	return nil
}

// func (c *Pgx) Update(ctx context.Context, channel Channel) error {
// 	query, args, err := repository.Sq.Update("channel").Set(
// 		"channel_id",
// 		channel.ChannelID,
// 	).Set(
// 		"service",
// 		channel.Service,
// 	).Where(
// 		"id = ?",
// 		channel.ID,
// 	).ToSql()
//
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = c.pg.Exec(ctx, query, args...)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := repository.Sq.
		Delete(tableName).
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return repository.ErrBadQuery
	}

	res, err := c.pg.Exec(ctx, query, args...)
	if err != nil {
		// TODO: viyasnit kak luchshe
		return errors.Join(err, ErrCannotDelete)
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
