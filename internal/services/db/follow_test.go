package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFollowCreate(t *testing.T) {
	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatService(entClient)
	channelsService := NewChannelService(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", chat.ServiceTelegram)
	assert.NoError(t, err)

	newChannel, err := channelsService.Create(ctx, "1", channel.ServiceTwitch)
	assert.NoError(t, err)

	f, err := service.Create(ctx, newChannel.ID, newChat.ID)
	assert.NoError(t, err)
	assert.Equal(t, newChannel.ID, f.Edges.Channel.ID, "Expects channel_id to be equal.")
	assert.Equal(t, newChat.ID, f.Edges.Chat.ID, "Expects chat_id to be equal.")
}

func TestGetByChatAndChannel(t *testing.T) {
	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatService(entClient)
	channelsService := NewChannelService(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", chat.ServiceTelegram)
	assert.NoError(t, err)

	newChannel, err := channelsService.Create(ctx, "1", channel.ServiceTwitch)
	assert.NoError(t, err)

	_, err = service.Create(ctx, newChannel.ID, newChat.ID)
	assert.NoError(t, err)

	f, err := service.GetByChatAndChannel(ctx, newChat.ID, newChannel.ID)
	assert.NoError(t, err)
	assert.Equal(t, newChannel.ID, f.Edges.Channel.ID, "Expects channel_id to be equal.")
	assert.Equal(t, newChat.ID, f.Edges.Chat.ID, "Expects chat_id to be equal.")
}
