package follow

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (Follow, error) {
	follow := Follow{}

	query, args, err := repository.Sq.
		Select(
			"id",
			"chat_id",
			"channel_id",
			"created_at",
		).
		From("follows").
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return follow, err
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(
		&follow.ID,
		&follow.ChatID,
		&follow.ChannelID,
		&follow.CreatedAt,
	)

	return follow, err
}

func (c *Pgx) GetByChatID(ctx context.Context, chatID uuid.UUID) ([]Follow, error) {
	query, args, err := repository.Sq.
		Select(
			"id",
			"chat_id",
			"channel_id",
			"created_at",
		).
		From("follows").
		Where(
			"chat_id = ?",
			chatID,
		).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pg.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var follows []Follow
	for rows.Next() {
		follow := Follow{}
		err = rows.Scan(
			&follow.ID,
			&follow.ChatID,
			&follow.ChannelID,
			&follow.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	return follows, err
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID uuid.UUID) ([]Follow, error) {
	query, args, err := repository.Sq.
		Select(
			"id",
			"chat_id",
			"channel_id",
			"created_at",
		).
		From("follows").
		Where(
			"channel_id = ?",
			channelID,
		).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pg.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var follows []Follow
	for rows.Next() {
		follow := Follow{}
		err = rows.Scan(
			&follow.ID,
			&follow.ChatID,
			&follow.ChannelID,
			&follow.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	return follows, err
}

func (c *Pgx) Create(ctx context.Context, follow Follow) error {
	query, args, err := repository.Sq.
		Insert("follows").
		Columns(
			"id",
			"chat_id",
			"channel_id",
		).
		Values(
			follow.ID,
			follow.ChatID,
			follow.ChannelID,
		).ToSql()
	if err != nil {
		return err
	}

	_, err = c.pg.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := repository.Sq.
		Delete("follows").
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return err
	}

	_, err = c.pg.Exec(ctx, query, args...)
	return err
}
