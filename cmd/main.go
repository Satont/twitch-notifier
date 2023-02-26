package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satont/twitch-notifier/ent"
	channel2 "github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/satont/twitch-notifier/internal/services/config"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/internal/services/twitch_streams_cheker"
	"github.com/satont/twitch-notifier/internal/services/types"
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

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	settings, _ := client.ChatSettings.Create().Save(context.Background())

	c, err := client.Chat.
		Create().
		SetID(uuid.New()).
		SetChatID("1").
		SetService(chat.ServiceTelegram).
		SetChatSettings(settings).
		Save(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	ch, err := client.Channel.Create().SetChannelID("1").SetService(channel2.ServiceTwitch).Save(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	client.Follow.Create().SetChat(c).SetChannel(ch).Save(context.Background())

	f := client.Follow.Query().WithChat(func(query *ent.ChatQuery) {
		query.WithChatSettings()
	}).WithChannel().FirstX(context.Background())

	spew.Dump(f.Edges)

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

	twitch_streams_cheker.NewTwitchStreamCheker(services.Twitch).StartPolling()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
}
