package follow

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

const tableName = "follows"

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (domain.Follow, error) {
	follow := Follow{}

	query, args, err := repository.Sq.
		Select(
			"id",
			"chat_id",
			"channel_id",
			"created_at",
		).
		From(tableName).
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return domain.Follow{}, repository.ErrBadQuery
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(
		&follow.ID,
		&follow.ChatID,
		&follow.ChannelID,
		&follow.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Follow{}, ErrNotFound
		}

		return domain.Follow{}, err
	}

	return domain.Follow{
		ID:        follow.ID,
		ChatID:    follow.ChatID,
		ChannelID: follow.ChannelID,
		CreatedAt: follow.CreatedAt,
	}, err
}

func (c *Pgx) GetByChatID(ctx context.Context, chatID uuid.UUID) ([]domain.Follow, error) {
	query, args, err := repository.Sq.
		Select(
			"id",
			"chat_id",
			"channel_id",
			"created_at",
		).
		From(tableName).
		Where(
			"chat_id = ?",
			chatID,
		).ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
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

	domainFollows := make([]domain.Follow, len(follows))
	for i, follow := range follows {
		domainFollows[i] = domain.Follow{
			ID:        follow.ID,
			ChatID:    follow.ChatID,
			ChannelID: follow.ChannelID,
			CreatedAt: follow.CreatedAt,
		}
	}

	return domainFollows, err
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID uuid.UUID) ([]domain.Follow, error) {
	query, args, err := repository.Sq.
		Select(
			"id",
			"chat_id",
			"channel_id",
			"created_at",
		).
		From(tableName).
		Where(
			"channel_id = ?",
			channelID,
		).ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
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

	domainFollows := make([]domain.Follow, len(follows))
	for i, follow := range follows {
		domainFollows[i] = domain.Follow{
			ID:        follow.ID,
			ChatID:    follow.ChatID,
			ChannelID: follow.ChannelID,
			CreatedAt: follow.CreatedAt,
		}
	}

	return domainFollows, err
}

func (c *Pgx) Create(ctx context.Context, follow domain.Follow) error {
	query, args, err := repository.Sq.
		Insert(tableName).
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
		return repository.ErrBadQuery
	}

	_, err = c.pg.Exec(ctx, query, args...)
	if err != nil {
		return errors.Join(err, ErrCannotCreate)
	}
	return err
}

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

	rows, err := c.pg.Exec(ctx, query, args...)
	if err != nil {
		return errors.Join(err, ErrCannotDelete)
	}

	if rows.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
