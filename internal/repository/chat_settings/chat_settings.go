package chat_settings

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
)

type ChatSettings struct {
	ID                        uuid.UUID
	ChatID                    uuid.UUID
	Language                  domain.Language
	GameChangeNotifications   bool
	TitleChangeNotifications  bool
	OfflineNotifications      bool
	GameAndTitleNotifications bool
	ShotThumbnail             bool
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (ChatSettings, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID) (ChatSettings, error)
	GetAll(ctx context.Context) ([]ChatSettings, error)
	Create(ctx context.Context, chatSettings ChatSettings) error
	Update(ctx context.Context, chatSettings ChatSettings) error
	Delete(ctx context.Context, id uuid.UUID) error
}
