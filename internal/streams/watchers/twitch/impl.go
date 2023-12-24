package twitch

type Opts struct {
	Handler streamshandler.StreamsHandler
}

func New(opts Opts) *Impl {
	return &Impl{
		handler: opts.Handler,
	}
}

var _ Watcher = (*Impl)(nil)

type Impl struct {
	handler streamshandler.StreamsHandler
}

func (c *Impl) Start() error {
	// TODO implement me
	panic("implement me")
}
