package db

import (
	"context"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFollowCreate(t *testing.T) {
	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatEntRepository(entClient)
	channelsService := NewChannelEntRepository(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", db_models.ChatServiceTelegram)
	assert.NoError(t, err)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	f, err := service.Create(ctx, newChannel.ID, newChat.ID)
	assert.NoError(t, err)
	assert.Equal(t, newChannel.ID, f.ChannelID, "Expects channel_id to be equal.")
	assert.Equal(t, newChat.ID, f.ChannelID, "Expects chat_id to be equal.")
}

func TestGetByChatAndChannel(t *testing.T) {
	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatEntRepository(entClient)
	channelsService := NewChannelEntRepository(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", db_models.ChatServiceTelegram)
	assert.NoError(t, err)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	_, err = service.Create(ctx, newChannel.ID, newChat.ID)
	assert.NoError(t, err)

	f, err := service.GetByChatAndChannel(ctx, newChat.ID, newChannel.ID)
	assert.NoError(t, err)
	assert.Equal(t, newChannel.ID, f.ChannelID, "Expects channel_id to be equal.")
	assert.Equal(t, newChat.ID, f.ChatID, "Expects chat_id to be equal.")
}
