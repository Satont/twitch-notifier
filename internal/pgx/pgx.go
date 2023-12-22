package pgx

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satont/twitch-notifier/pkg/logger"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Logger logger.Logger
}

func New(opts Opts) (*pgxpool.Pool, error) {
	pgx, err := pgxpool.New(
		context.Background(),
		"postgres://postgres:postgres@localhost:5432/twitch_notifier",
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
				// return pgx.Ping(ctx)
				opts.Logger.Info("Connected to postgres")
				return nil
			},
		},
	)

	return pgx, nil
}
