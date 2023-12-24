package chat

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

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	query, args, err := repository.Sq.
		Select("id", "chat_id", "service").
		From("chats").
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	chat := Chat{}
	err = c.pg.QueryRow(ctx, query, args...).Scan(&chat.ID, &chat.ChatID, &chat.Service)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &domain.Chat{
		ID:      chat.ID,
		Service: domain.ChatService(chat.Service),
		ChatID:  chat.ChatID,
	}, nil
}

func (c *Pgx) GetByChatServiceAndChatID(
	ctx context.Context,
	service ChatService,
	chatID string,
) (*domain.Chat, error) {
	query, args, err := repository.Sq.
		Select("id", "chat_id", "service").
		From("chats").
		Where(
			"chat_id = ? AND service = ?",
			chatID,
			service,
		).ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	chat := Chat{}
	err = c.pg.QueryRow(ctx, query, args...).Scan(&chat.ID, &chat.ChatID, &chat.Service)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &domain.Chat{
		ID:      chat.ID,
		Service: domain.ChatService(chat.Service),
		ChatID:  chat.ChatID,
	}, nil
}

func (c *Pgx) GetAll(ctx context.Context) ([]domain.Chat, error) {
	query, args, err := repository.Sq.
		Select("id", "chat_id", "service").
		From("chats").
		ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	rows, err := c.pg.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var chats []Chat
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.ID, &chat.ChatID, &chat.Service)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	domainChats := make([]domain.Chat, len(chats))
	for i, chat := range chats {
		domainChats[i] = domain.Chat{
			ID:      chat.ID,
			Service: domain.ChatService(chat.Service),
			ChatID:  chat.ChatID,
		}
	}

	return domainChats, nil
}

func (c *Pgx) Create(ctx context.Context, user domain.Chat) error {
	query, args, err := repository.Sq.
		Insert("chats").
		Columns(
			"id",
			"chat_id",
			"service",
		).Values(
		user.ID,
		user.ChatID,
		user.Service,
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

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := repository.Sq.Delete("chats").Where(
		"id = ?",
		id,
	).ToSql()
	if err != nil {
		return repository.ErrBadQuery
	}

	rows, err := c.pg.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return errors.Join(err, ErrCannotDelete)
	}

	return nil
}
