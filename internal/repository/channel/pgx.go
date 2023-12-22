package channel

import (
	"context"

	"github.com/google/uuid"
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

func (c *Pgx) GetById(ctx context.Context, id uuid.UUID) (Channel, error) {
	channel := Channel{}

	query, args, err := repository.Sq.
		Select("id", "channel_id", "service").
		From("channels").
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return channel, err
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(&channel.ID, &channel.ChannelID, &channel.Service)
	if err != nil {
		return channel, err
	}

	return channel, nil
}

func (c *Pgx) GetByStreamServiceAndID(
	ctx context.Context,
	service domain.StreamingService,
	id string,
) (Channel, error) {
	channel := Channel{}

	query, args, err := repository.Sq.
		Select("id", "channel_id", "service").
		From("channels").
		Where(
			"channel_id = ? AND service = ?",
			id,
			service,
		).ToSql()
	if err != nil {
		return channel, err
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(&channel.ID, &channel.ChannelID, &channel.Service)
	if err != nil {
		return channel, err
	}

	return channel, nil
}

func (c *Pgx) GetAll(ctx context.Context) ([]Channel, error) {
	var channels []Channel

	query, args, err := repository.Sq.
		Select("id", "channel_id", "service").
		From("channels").
		ToSql()
	if err != nil {
		return channels, err
	}

	rows, err := c.pg.Query(ctx, query, args...)
	if err != nil {
		return channels, err
	}

	defer rows.Close()
	for rows.Next() {
		channel := Channel{}
		err = rows.Scan(&channel.ID, &channel.ChannelID, &channel.Service)
		if err != nil {
			return channels, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

func (c *Pgx) Create(ctx context.Context, channel Channel) error {
	query, args, err := repository.Sq.Insert("channels").Columns(
		"id",
		"channel_id",
		"service",
	).Values(
		channel.ID,
		channel.ChannelID,
		channel.Service,
	).ToSql()
	if err != nil {
		return err
	}

	_, err = c.pg.Exec(ctx, query, args...)
	if err != nil {
		return err
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
		Delete("channel").
		Where(
			"id = ?",
			id,
		).ToSql()

	if err != nil {
		return err
	}

	_, err = c.pg.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
