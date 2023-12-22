package store

import (
	"github.com/satont/twitch-notifier/internal/domain"
)

type I18nStore interface {
	GetKey(language domain.Language, key string) (string, error)
	GetSupportedLanguages() []domain.Language
}
