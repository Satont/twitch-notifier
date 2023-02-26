package db

import (
	"context"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/stretchr/testify/assert"
	"testing"
)

func FollowCreateTest(t *testing.T) {
	entClient, err := setupTest()
	assert.Nil(t, err, "Expects err to be nil")
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatService(entClient)
	channelsService := NewChannelService(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", chat.ServiceTelegram)
	assert.Nil(t, err, "Expects err to be nil")

	newChannel, err := channelsService.Create(ctx, "1", channel.ServiceTwitch)
	assert.Nil(t, err, "Expects err to be nil")

	f, err := service.Create(ctx, newChannel.ID, newChat.ID)
	assert.Nil(t, err, "Expects err to be nil")
	assert.Equal(t, newChannel.ID, f.Edges.Channel.ID, "Expects channel_id to be equal.")
	assert.Equal(t, newChat.ID, f.Edges.Chat.ID, "Expects chat_id to be equal.")
}
