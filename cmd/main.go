package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"

	"entgo.io/ent/dialect/sql"
	"github.com/TheZeroSlave/zapsentry"
	"github.com/lib/pq"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/internal/config"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/message_sender"
	"github.com/satont/twitch-notifier/internal/telegram"
	"github.com/satont/twitch-notifier/internal/twitch"
	"github.com/satont/twitch-notifier/internal/twitch_streams_cheker"
	"github.com/satont/twitch-notifier/internal/types"
	"github.com/satont/twitch-notifier/pkg/i18n"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createEnt(cfg *config.Config) (*ent.Client, error) {
	pgConnectionUrl, err := pq.ParseURL(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalln(err)
	}

	drv, err := sql.Open("postgres", pgConnectionUrl)
	if err != nil {
		return nil, err
	}

	db := drv.DB()
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)
	return ent.NewClient(ent.Driver(drv)), nil
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	cfg, err := config.NewConfig(nil)
	if err != nil {
		log.Fatalln(err)
	}

	logger, _ := zap.NewDevelopment()

	if cfg.SentryDsn != "" {
		sentryClient, err := sentry.NewClient(
			sentry.ClientOptions{
				Dsn:           cfg.SentryDsn,
				EnableTracing: true,
			},
		)
		if err != nil {
			log.Fatalln(err)
		}
		logger = modifyToSentryLogger(logger, sentryClient)
		defer sentry.Flush(2 * time.Second)
	}

	zap.ReplaceGlobals(logger)

	client, err := createEnt(cfg)
	if err != nil {
		logger.Sugar().Fatalln("failed opening connection to postgres: %v", err)
	}
	// Run the auto migration tool.
	// if err := client.Schema.Create(context.Background()); err != nil {
	//	log.Fatalf("failed creating schema resources: %v", err)
	// }

	twitchService, err := twitch.NewTwitchService(cfg.TwitchClientId, cfg.TwitchClientSecret)
	if err != nil {
		logger.Sugar().Fatalln(err)
	}

	i18, err := i18n.NewI18n(filepath.Join(wd, "locales"))
	if err != nil {
		logger.Sugar().Fatalln(err)
	}

	services := &types.Services{
		Config:  cfg,
		Twitch:  twitchService,
		Chat:    db.NewChatEntRepository(client),
		Channel: db.NewChannelEntService(client),
		Follow:  db.NewFollowService(client),
		Stream:  db.NewStreamEntService(client),
		I18N:    i18,
	}

	ctx, cancel := context.WithCancel(context.Background())

	tg := telegram.NewTelegram(ctx, cfg.TelegramToken, services)
	tg.StartPolling(ctx)

	sender := message_sender.NewMessageSender(tg.Client)

	checker := twitch_streams_cheker.NewTwitchStreamChecker(services, sender, nil)
	checker.StartPolling(ctx)

	logger.Sugar().Info("Started")
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	logger.Sugar().Info("Closing...")
	cancel()
	_ = client.Close()
}

func modifyToSentryLogger(log *zap.Logger, client *sentry.Client) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, // when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))

	// in case of err it will return noop core. so we can safely attach it
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}

	log = zapsentry.AttachCoreToLogger(core, log)

	// to use breadcrumbs feature - create new scope explicitly
	// and attach after attaching the core
	return log.With(zapsentry.NewScope())
}
