package chat

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

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (Chat, error) {
	chat := Chat{}

	query, args, err := repository.Sq.
		Select("id", "chat_id", "service").
		From("chats").
		Where(
			"id = ?",
			id,
		).ToSql()

	if err != nil {
		return chat, err
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(&chat.ID, &chat.ChatID, &chat.Service)
	if err != nil {
		return chat, err
	}

	return chat, nil
}

func (c *Pgx) GetByChatServiceAndChatID(
	ctx context.Context,
	service domain.ChatService,
	chatID string,
) (Chat, error) {
	chat := Chat{}

	query, args, err := repository.Sq.
		Select("id", "chat_id", "service").
		From("chats").
		Where(
			"chat_id = ? AND service = ?",
			chatID,
			service,
		).ToSql()
	if err != nil {
		return chat, err
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(&chat.ID, &chat.ChatID, &chat.Service)
	if err != nil {
		return chat, err
	}

	return chat, nil
}

func (c *Pgx) GetAll(ctx context.Context) ([]Chat, error) {
	var chats []Chat

	query, args, err := repository.Sq.
		Select("id", "chat_id", "service").
		From("chats").
		ToSql()
	if err != nil {
		return chats, err
	}

	rows, err := c.pg.Query(ctx, query, args...)
	if err != nil {
		return chats, err
	}

	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.ID, &chat.ChatID, &chat.Service)
		if err != nil {
			return chats, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (c *Pgx) Create(ctx context.Context, user Chat) error {
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
		return err
	}

	_, err = c.pg.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := repository.Sq.Delete("chats").Where(
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
