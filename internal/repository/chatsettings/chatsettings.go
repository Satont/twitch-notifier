package chatsettings

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/domain"
)

type ChatSettings struct {
	ID                            uuid.UUID
	ChatID                        uuid.UUID
	Language                      Language
	CategoryChangeNotifications   bool
	TitleChangeNotifications      bool
	OfflineNotifications          bool
	CategoryAndTitleNotifications bool
	ShowThumbnail                 bool
}

type Language string

func (l Language) String() string {
	return string(l)
}

const (
	LanguageEN Language = "en"
	LanguageRU Language = "ru"
	LanguageUA Language = "ua"
)

var ErrNotFound = errors.New("chatsettings not found")
var ErrCannotCreate = errors.New("cannot create chatsettings")
var ErrCannotDelete = errors.New("cannot delete chatsettings")

//go:generate go run go.uber.org/mock/mockgen -source=chatsettings.go -destination=mocks/mock.go

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatSettings, error)
	GetByChatID(ctx context.Context, chatID uuid.UUID) (*domain.ChatSettings, error)
	Create(ctx context.Context, chatSettings domain.ChatSettings) error
	Update(ctx context.Context, chatSettings domain.ChatSettings) error
	Delete(ctx context.Context, id uuid.UUID) error
}
