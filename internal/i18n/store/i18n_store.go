package store

import (
	"github.com/satont/twitch-notifier/internal/domain"
)

//go:generate go run go.uber.org/mock/mockgen -source=i18n_store.go -destination=mocks/mock.go

type I18nStore interface {
	GetKey(language domain.Language, key string) (string, error)
	GetSupportedLanguages() []domain.Language
}
