package localizer

import (
	"errors"

	"github.com/satont/twitch-notifier/internal/domain"
)

var ErrKeyIsEmpty = errors.New("key is empty")
var ErrTranslateError = errors.New("failed to translate")

//go:generate go run go.uber.org/mock/mockgen -source=localizer.go -destination=mocks/mock.go

type Localizer interface {
	Localize(opts ...Option) (string, error)
	MustLocalize(opts ...Option) string
}

type Options struct {
	key        string
	language   domain.Language
	attributes map[string]any
}

type (
	Option interface {
		apply(options *Options)
	}

	applyFunc func(options *Options)
)
