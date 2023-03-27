package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/internal/services/config"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/message_sender"
	"github.com/satont/twitch-notifier/internal/services/telegram"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/internal/services/twitch_streams_cheker"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/pkg/i18n"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	cfg, err := config.NewConfig(nil)
	if err != nil {
		log.Fatalln(err)
	}

	var logger *zap.Logger

	logger, _ = zap.NewDevelopment()

	zap.ReplaceGlobals(logger)

	client, err := ent.Open("postgres", cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// Run the auto migration tool.
	//if err := client.Schema.Create(context.Background()); err != nil {
	//	log.Fatalf("failed creating schema resources: %v", err)
	//}

	twitchService, err := twitch.NewTwitchService(cfg.TwitchClientId, cfg.TwitchClientSecret)
	if err != nil {
		log.Fatalln(err)
	}

	i18, err := i18n.NewI18n(filepath.Join(wd, "locales"))
	if err != nil {
		log.Fatalln(err)
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
	fmt.Println("Closing...")
	cancel()
	_ = client.Close()
}
