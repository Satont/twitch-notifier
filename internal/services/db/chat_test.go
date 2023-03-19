package db

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satont/twitch-notifier/ent"
	channel2 "github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTest() (*ent.Client, error) {
	entClient, err := ent.Open("sqlite3", "file:tests?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, err
	}
	if err := entClient.Schema.Create(context.Background()); err != nil {
		return nil, err
	}
	return entClient, nil
}

func teardownTest(entClient *ent.Client) {
	entClient.Close()
}

func TestChatService_GetByID(t *testing.T) {
	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	_, err = entClient.Chat.
		Create().
		SetID(uuid.New()).
		SetChatID("1").
		SetService(chat.ServiceTelegram).
		Save(ctx)

	if err != nil {
		t.Fatal(err)
	}

	c := &chatService{
		entClient: entClient,
	}
	got, err := c.GetByID(ctx, "1", chat.ServiceTelegram)
	assert.NoError(t, err)
	assert.Equal(t, "1", got.ChatID, "Expects chat_id to be 1.")

	got, err = c.GetByID(ctx, "2", chat.ServiceTelegram)
	assert.NoError(t, err)
	assert.Nil(t, got, "Expects got to be nil")
}

func TestChatService_Create(t *testing.T) {
	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	c := &chatService{
		entClient: entClient,
	}
	got, err := c.Create(ctx, "1", chat.ServiceTelegram)
	assert.NoError(t, err)
	assert.Equal(t, "1", got.ChatID, "Expects chat_id to be 1.")

	settings := got.QuerySettings().OnlyX(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, settings.GameChangeNotification, "Expects game_change_notification to be true")

	got, err = c.Create(ctx, "1", chat.ServiceTelegram)
	assert.Nil(t, got)
	assert.Error(t, err)
}

func TestChatService_GetFollowsByID(t *testing.T) {
	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	channel, err := entClient.Channel.
		Create().
		SetID(uuid.New()).
		SetService(channel2.ServiceTwitch).
		SetChannelID("1").
		Save(ctx)
	assert.NoError(t, err)

	ch, err := entClient.Chat.
		Create().
		SetID(uuid.New()).
		SetChatID("1").
		SetService(chat.ServiceTelegram).
		Save(ctx)

	_, err = entClient.Follow.Create().SetID(uuid.New()).SetChannel(channel).SetChat(ch).Save(context.Background())

	c := &chatService{
		entClient: entClient,
	}
	got, err := c.GetFollowsByID(ctx, "1", chat.ServiceTelegram)
	assert.Equal(t, 1, len(got), "Expects got to be 1")
	assert.NoError(t, err)

	got, err = c.GetFollowsByID(ctx, "2", chat.ServiceTelegram)
	assert.Equal(t, 0, len(got), "Expects got to be empty slice")
	assert.NoError(t, err)
}
