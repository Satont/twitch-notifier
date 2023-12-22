package twitch

import (
	"github.com/satont/twitch-notifier/internal/streams_handler"
)

type Opts struct {
	Handler streams_handler.StreamsHandler
}

func New(opts Opts) *Impl {
	return &Impl{
		handler: opts.Handler,
	}
}

var _ Watcher = (*Impl)(nil)

type Impl struct {
	handler streams_handler.StreamsHandler
}

func (c *Impl) Start() error {
	// TODO implement me
	panic("implement me")
}
