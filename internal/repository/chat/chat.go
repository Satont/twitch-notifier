package chat

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
)

type Chat struct {
	ID      uuid.UUID
	Service ChatService
	ChatID  string
}

type ChatService string

func (c ChatService) String() string {
	return string(c)
}

const (
	ChatServiceTelegram ChatService = "telegram"
)

var ErrNotFound = errors.New("chat not found")
var ErrCannotCreate = errors.New("cannot create chat")
var ErrCannotDelete = errors.New("cannot delete chat")

//go:generate go run go.uber.org/mock/mockgen -source=chat.go -destination=mocks/mock.go

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
	GetByChatServiceAndChatID(ctx context.Context, service ChatService, chatID string) (
		*domain.Chat,
		error,
	)
	GetAll(ctx context.Context) ([]domain.Chat, error)
	Create(ctx context.Context, user domain.Chat) error
	// Update(ctx context.Context, user Chat) error
	Delete(ctx context.Context, id uuid.UUID) error
}
