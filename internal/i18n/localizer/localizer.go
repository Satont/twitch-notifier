package localizer

import (
	"github.com/satont/twitch-notifier/internal/domain"
)

type Localizer interface {
	Localize(opts ...Option) (string, error)
	MustLocalize(opts ...Option) string
}

type localizerOptions struct {
	key        string
	language   domain.Language
	attributes map[string]any
}

type (
	Option interface {
		apply(options *localizerOptions)
	}

	applyFunc func(options *localizerOptions)
)
