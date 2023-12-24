package twitchclientimpl

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/internal/twitchclient"
	"github.com/satont/twitch-notifier/pkg/config"
	"go.uber.org/fx"
)

type TwitchHelixImplOpts struct {
	fx.In

	Config config.Config
}

func New(opts TwitchHelixImplOpts) (*TwitchHelixImpl, error) {
	client, err := helix.NewClient(
		&helix.Options{
			ClientID:      opts.Config.TwitchClientID,
			ClientSecret:  opts.Config.TwitchClientSecret,
			RateLimitFunc: helixRateLimitCallback,
		},
	)
	if err != nil {
		return nil, err
	}

	return &TwitchHelixImpl{
		client: client,
	}, nil
}

type TwitchHelixImpl struct {
	client *helix.Client
}

var _ twitchclient.TwitchClient = (*TwitchHelixImpl)(nil)
