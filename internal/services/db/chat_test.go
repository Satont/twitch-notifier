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

func teardownTest(entClient *ent.Client) error {
	return entClient.Close()
}

func TestChatService_GetByID(t *testing.T) {
	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	//settings, err := entClient.ChatSettings.Create().Save(context.Background())
	//if err != nil {
	//	t.Fatal(err)
	//}

	_, err = entClient.Chat.
		Create().
		SetID(uuid.New()).
		SetChatID("1").
		SetService(chat.ServiceTelegram).
		//SetSettings(settings).
		Save(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	c := &chatService{
		entClient: entClient,
	}
	got, err := c.GetByID("1", chat.ServiceTelegram)
	if err != nil {
		t.Errorf("%v error = %v", "GetByID", err)
	}
	assert.Equal(t, "1", got.ChatID, "Expects chat_id to be 1.")

	got, err = c.GetByID("2", chat.ServiceTelegram)
	assert.Nil(t, err, "Expects err to be not nil")
	assert.Nil(t, got, "Expects got to be nil")
}

func TestChatService_Create(t *testing.T) {
	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	c := &chatService{
		entClient: entClient,
	}
	got, err := c.Create("1", chat.ServiceTelegram)
	assert.Equal(t, "1", got.ChatID, "Expects chat_id to be 1.")
	assert.Nil(t, err, "Expects got to be nil")
	settings := got.QuerySettings().OnlyX(context.Background())
	assert.NotNil(t, settings, "Expects settings to be not nil")
	assert.Equal(t, true, settings.GameChangeNotification, "Expects game_change_notification to be true")

	got, err = c.Create("1", chat.ServiceTelegram)
	assert.Nil(t, got, "Expects got to be nil")
	assert.NotNil(t, err, "Expects err to be not nil")
}

func TestChatService_GetFollowsByID(t *testing.T) {
	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	channel, err := entClient.Channel.
		Create().
		SetID(uuid.New()).
		SetService(channel2.ServiceTwitch).
		SetChannelID("1").
		Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	ch, err := entClient.Chat.
		Create().
		SetID(uuid.New()).
		SetChatID("1").
		SetService(chat.ServiceTelegram).
		Save(context.Background())

	_, err = entClient.Follow.Create().SetID(uuid.New()).SetChannel(channel).SetChat(ch).Save(context.Background())

	c := &chatService{
		entClient: entClient,
	}
	got, err := c.GetFollowsByID("1", chat.ServiceTelegram)
	assert.Equal(t, 1, len(got), "Expects got to be 1")
	assert.Nil(t, err, "Expects got to be nil")

	got, err = c.GetFollowsByID("2", chat.ServiceTelegram)
	assert.Equal(t, 0, len(got), "Expects got to be empty slice")
	assert.Nil(t, err, "Expects err to be not nil")
}
