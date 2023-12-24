package main

import (
	announcesender "github.com/satont/twitch-notifier/internal/announcesender/temporal"
	i18nstore "github.com/satont/twitch-notifier/internal/i18n/store"
	messagesender "github.com/satont/twitch-notifier/internal/messagesender/fx"
	"github.com/satont/twitch-notifier/internal/pgx"
	repositories "github.com/satont/twitch-notifier/internal/repository/fx"
	thumbnailchecker "github.com/satont/twitch-notifier/internal/thumbnailchecker/temporal"
	"github.com/satont/twitch-notifier/internal/twitchclient"
	"github.com/satont/twitch-notifier/internal/twitchclient/twitchclientimpl"
	"github.com/satont/twitch-notifier/pkg/config"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// fx.NopLogger,
		fx.Provide(
			config.New,
			logger.FxOption,
			fx.Annotate(
				i18nstore.New,
				fx.As(new(i18nstore.I18nStore)),
			),
			pgx.New,
			fx.Annotate(
				twitchclientimpl.New,
				fx.As(new(twitchclient.TwitchClient)),
			),
		),
		repositories.Module,
		thumbnailchecker.Module,
		messagesender.Module,
		announcesender.Module,
		fx.Invoke(pgx.New),
	).Run()
}
