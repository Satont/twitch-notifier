package localizer

import (
	"errors"
	"fmt"

	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/i18n/store"
)

type LocalizeOpts struct {
	Store store.I18nStore
}

func NewLocalizer(i18nStore store.I18nStore) *Impl {
	return &Impl{store: i18nStore}
}

var _ Localizer = (*Impl)(nil)

type Impl struct {
	store store.I18nStore
}

const defaultLanguage = domain.LanguageEN

var ErrKeyIsEmpty = errors.New("key is empty")

func (c *Impl) Localize(opts ...Option) (string, error) {
	options := &localizerOptions{
		attributes: make(map[string]any),
	}
	for _, opt := range opts {
		opt.apply(options)
	}

	if options.key == "" {
		return "", ErrKeyIsEmpty
	}

	key, err := c.store.GetKey(options.language, options.key)
	if err != nil {
		return "", fmt.Errorf("failed to get key: %w", err)
	}

	return key, nil
}

func (c *Impl) MustLocalize(opts ...Option) string {
	key, err := c.Localize(opts...)
	if err != nil {
		opts = append(
			opts,
			WithLanguage(defaultLanguage),
		)
		key, err = c.Localize(opts...)
		if err != nil {
			panic(err)
		}
	}

	return key
}

func (f applyFunc) apply(s *localizerOptions) { f(s) }

func WithKey(key string) Option {
	return applyFunc(
		func(s *localizerOptions) {
			s.key = key
		},
	)
}

func WithLanguage(language domain.Language) Option {
	return applyFunc(
		func(s *localizerOptions) {
			s.language = language
		},
	)
}

func WithAttribute(key string, value any) Option {
	return applyFunc(
		func(s *localizerOptions) {
			s.attributes[key] = value
		},
	)
}
