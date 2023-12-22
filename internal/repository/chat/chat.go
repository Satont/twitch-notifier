package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
)

type User struct {
	ID      uuid.UUID
	Service domain.ChatService
	ChatID  string
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByChatServiceAndChatID(ctx context.Context, service domain.ChatService, chatID string) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
