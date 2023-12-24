package chatsettings

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

const tableName = "chat_settings"

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatSettings, error) {
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
		From(tableName).
		Where(
			"id = ?",
			id,
		).ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	settings := ChatSettings{}
	err = c.pg.QueryRow(ctx, query, args...).Scan(
		&settings.ID,
		&settings.ChatID,
		&settings.Language,
		&settings.CategoryChangeNotifications,
		&settings.TitleChangeNotifications,
		&settings.OfflineNotifications,
		&settings.CategoryAndTitleNotifications,
		&settings.ShowThumbnail,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &domain.ChatSettings{
		ID:                            settings.ID,
		ChatID:                        settings.ChatID,
		Language:                      domain.Language(settings.Language),
		CategoryChangeNotifications:   settings.CategoryChangeNotifications,
		TitleChangeNotifications:      settings.TitleChangeNotifications,
		OfflineNotifications:          settings.OfflineNotifications,
		CategoryAndTitleNotifications: settings.CategoryAndTitleNotifications,
		ShowThumbnail:                 settings.ShowThumbnail,
	}, err
}

func (c *Pgx) GetByChatID(ctx context.Context, chatID uuid.UUID) (*domain.ChatSettings, error) {
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
		From(tableName).
		Where(
			"chat_id = ?",
			chatID,
		).ToSql()
	if err != nil {
		return nil, repository.ErrBadQuery
	}

	settings := ChatSettings{}
	err = c.pg.QueryRow(ctx, query, args...).Scan(
		&settings.ID,
		&settings.ChatID,
		&settings.Language,
		&settings.CategoryChangeNotifications,
		&settings.TitleChangeNotifications,
		&settings.OfflineNotifications,
		&settings.CategoryAndTitleNotifications,
		&settings.ShowThumbnail,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &domain.ChatSettings{
		ID:                            settings.ID,
		ChatID:                        settings.ChatID,
		Language:                      domain.Language(settings.Language),
		CategoryChangeNotifications:   settings.CategoryChangeNotifications,
		TitleChangeNotifications:      settings.TitleChangeNotifications,
		OfflineNotifications:          settings.OfflineNotifications,
		CategoryAndTitleNotifications: settings.CategoryAndTitleNotifications,
		ShowThumbnail:                 settings.ShowThumbnail,
	}, err
}

func (c *Pgx) Create(ctx context.Context, chatSettings domain.ChatSettings) error {
	query, args, err := repository.Sq.Insert(tableName).Columns(
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
		chatSettings.CategoryChangeNotifications,
		chatSettings.TitleChangeNotifications,
		chatSettings.OfflineNotifications,
		chatSettings.CategoryAndTitleNotifications,
		chatSettings.ShowThumbnail,
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

func (c *Pgx) Update(ctx context.Context, chatSettings domain.ChatSettings) error {
	query, args, err := repository.Sq.
		Update(tableName).
		Set(
			"language",
			chatSettings.Language,
		).
		Set(
			"game_change_notifications",
			chatSettings.CategoryChangeNotifications,
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
			chatSettings.CategoryAndTitleNotifications,
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
		return repository.ErrBadQuery
	}

	_, err = c.pg.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := repository.Sq.Delete(tableName).Where(
		"id = ?",
		id,
	).ToSql()
	if err != nil {
		return repository.ErrBadQuery
	}

	rows, err := c.pg.Exec(ctx, query, args...)
	if err != nil {
		return ErrCannotDelete
	}

	if rows.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
