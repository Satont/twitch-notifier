package pgx

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satont/twitch-notifier/pkg/config"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Logger logger.Logger
	Config config.Config
}

func New(opts Opts) (*pgxpool.Pool, error) {
	pgx, err := pgxpool.New(
		context.Background(),
		opts.Config.PostgresUrl,
	)
	if err != nil {
		return nil, err
	}

	opts.LC.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				pgx.Close()
				return nil
			},
			OnStart: func(ctx context.Context) error {
				err := pgx.Ping(ctx)
				if err != nil {
					return err
				}

				opts.Logger.Info("Connected to postgres")
				return nil
			},
		},
	)

	return pgx, nil
}
