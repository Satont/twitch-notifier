package chat_settings

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

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (ChatSettings, error) {
	settings := ChatSettings{}

	query, args, err := repository.Sq.
		Select(
			"id",
			"chat_id",
			"language",
			"game_change_notifications",
			"title_change_notifications",
			"offline_notifications",
			"game_and_title_notifications",
			"show_thumbnail",
		).
		From("chat_settings").
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return settings, err
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(
		&settings.ID,
		&settings.ChatID,
		&settings.Language,
		&settings.GameChangeNotifications,
		&settings.TitleChangeNotifications,
		&settings.OfflineNotifications,
		&settings.GameAndTitleNotifications,
		&settings.ShowThumbnail,
	)

	return settings, err
}

func (c *Pgx) GetByChatID(ctx context.Context, chatID uuid.UUID) (ChatSettings, error) {
	settings := ChatSettings{}

	query, args, err := repository.Sq.
		Select(
			"id",
			"chat_id",
			"language",
			"game_change_notifications",
			"title_change_notifications",
			"offline_notifications",
			"game_and_title_notifications",
			"show_thumbnail",
		).
		From("chat_settings").
		Where(
			"chat_id = ?",
			chatID,
		).ToSql()
	if err != nil {
		return settings, err
	}

	err = c.pg.QueryRow(ctx, query, args...).Scan(
		&settings.ID,
		&settings.ChatID,
		&settings.Language,
		&settings.GameChangeNotifications,
		&settings.TitleChangeNotifications,
		&settings.OfflineNotifications,
		&settings.GameAndTitleNotifications,
		&settings.ShowThumbnail,
	)

	return settings, err
}

func (c *Pgx) Create(ctx context.Context, chatSettings ChatSettings) error {
	query, args, err := repository.Sq.Insert("chat_settings").Columns(
		"id",
		"chat_id",
		"language",
		"game_change_notifications",
		"title_change_notifications",
		"offline_notifications",
		"game_and_title_notifications",
		"show_thumbnail",
	).Values(
		chatSettings.ID,
		chatSettings.ChatID,
		chatSettings.Language,
		chatSettings.GameChangeNotifications,
		chatSettings.TitleChangeNotifications,
		chatSettings.OfflineNotifications,
		chatSettings.GameAndTitleNotifications,
		chatSettings.ShowThumbnail,
	).ToSql()
	if err != nil {
		return err
	}

	_, err = c.pg.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Update(ctx context.Context, chatSettings ChatSettings) error {
	query, args, err := repository.Sq.
		Update("chat_settings").
		Set(
			"language",
			chatSettings.Language,
		).
		Set(
			"game_change_notifications",
			chatSettings.GameChangeNotifications,
		).
		Set(
			"title_change_notifications",
			chatSettings.TitleChangeNotifications,
		).
		Set(
			"offline_notifications",
			chatSettings.OfflineNotifications,
		).
		Set(
			"game_and_title_notifications",
			chatSettings.GameAndTitleNotifications,
		).
		Set(
			"show_thumbnail",
			chatSettings.ShowThumbnail,
		).
		Where(
			"id = ?",
			chatSettings.ID,
		).ToSql()
	if err != nil {
		return err
	}

	_, err = c.pg.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := repository.Sq.Delete("chat_settings").Where(
		"id = ?",
		id,
	).ToSql()
	if err != nil {
		return err
	}

	_, err = c.pg.Exec(ctx, query, args...)
	return err
}
