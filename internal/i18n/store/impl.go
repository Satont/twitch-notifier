package store

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/pkg/logger"
)

type translation map[string]any

func New(l logger.Logger) (*Store, error) {
	store := &Store{
		locales: make(map[domain.Language]translation),
	}

	err := store.readLocales()
	if err != nil {
		return nil, err
	}

	supported := store.GetSupportedLanguages()
	l.Info("Locales loaded", slog.Any("locales", supported))

	return store, nil
}

type Store struct {
	locales map[domain.Language]translation
}

var _ I18nStore = (*Store)(nil)

func (c *Store) readLocales() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := os.ReadDir(filepath.Join(pwd, "locales"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			panic("locales directory should not contain directories")
		}

		data := make(map[string]any)

		content, err := os.ReadFile(filepath.Join(pwd, "locales", file.Name()))
		if err != nil {
			return err
		}

		if err := json.Unmarshal(content, &data); err != nil {
			return err
		}

		name := strings.Replace(file.Name(), ".json", "", 1)

		c.locales[domain.Language(name)] = data
	}

	return nil
}

var ErrLocaleNotFound = errors.New("locale not found")
var ErrKeyNotFound = errors.New("key not found")

func (c *Store) GetKey(language domain.Language, key string) (string, error) {
	lang := c.locales[language]
	if lang == nil {
		return "", ErrLocaleNotFound
	}

	value, ok := lang[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	return value.(string), nil
}

func (c *Store) GetSupportedLanguages() []domain.Language {
	languages := make([]domain.Language, 0, len(c.locales))

	for language := range c.locales {
		languages = append(languages, language)
	}

	return languages
}
