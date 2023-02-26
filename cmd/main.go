package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/internal/services/config"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/telegram"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/internal/services/twitch_streams_cheker"
	"github.com/satont/twitch-notifier/internal/services/types"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	var logger *zap.Logger

	logger, _ = zap.NewDevelopment()

	zap.ReplaceGlobals(logger)

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	twitchService, err := twitch.NewTwitchService(cfg.TwitchClientId, cfg.TwitchClientSecret)
	if err != nil {
		log.Fatalln(err)
	}

	services := &types.Services{
		Twitch:  twitchService,
		Chat:    db.NewChatService(client),
		Channel: db.NewChannelService(client),
		Follow:  db.NewFollowService(client),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tg := telegram.NewTelegram(cfg.TelegramToken, services)
	tg.StartPolling(ctx)

	checker := twitch_streams_cheker.NewTwitchStreamCheker(services.Twitch)
	checker.StartPolling(ctx)

	logger.Sugar().Info("Started")
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
}
