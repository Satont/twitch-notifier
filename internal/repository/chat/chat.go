package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
)

type Chat struct {
	ID      uuid.UUID
	Service domain.ChatService
	ChatID  string
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Chat, error)
	GetByChatServiceAndChatID(ctx context.Context, service domain.ChatService, chatID string) (
		Chat,
		error,
	)
	GetAll(ctx context.Context) ([]Chat, error)
	Create(ctx context.Context, user Chat) error
	// Update(ctx context.Context, user Chat) error
	Delete(ctx context.Context, id uuid.UUID) error
}
