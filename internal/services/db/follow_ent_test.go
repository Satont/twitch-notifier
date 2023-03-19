package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFollowService_Create(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatEntRepository(entClient)
	channelsService := NewChannelEntService(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", db_models.ChatServiceTelegram)
	assert.NoError(t, err)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	table := []struct {
		name      string
		chatID    uuid.UUID
		channelID uuid.UUID
		wantErr   bool
	}{
		{
			name:      "Create follow",
			chatID:    newChat.ID,
			channelID: newChannel.ID,
			wantErr:   false,
		},
		{
			name:      "Should fail if follow already exists",
			chatID:    newChat.ID,
			channelID: newChannel.ID,
			wantErr:   true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			f, err := service.Create(ctx, tt.channelID, tt.chatID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, newChannel.ID, f.ChannelID, "Expects channel_id to be equal.")
				assert.Equal(t, newChat.ID, f.ChatID, "Expects chat_id to be equal.")
			}
		})
	}
}

func TestFollowService_Delete(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatEntRepository(entClient)
	channelsService := NewChannelEntService(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", db_models.ChatServiceTelegram)
	assert.NoError(t, err)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	foll, err := service.Create(ctx, newChannel.ID, newChat.ID)

	table := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "Delete follow",
			id:      foll.ID,
			wantErr: false,
		},
		{
			name:    "Should fail if follow does not exist",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFollowService_GetByChatAndChannel(t *testing.T) {
	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatEntRepository(entClient)
	channelsService := NewChannelEntService(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", db_models.ChatServiceTelegram)
	assert.NoError(t, err)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	_, err = service.Create(ctx, newChannel.ID, newChat.ID)
	assert.NoError(t, err)

	table := []struct {
		name      string
		chatID    uuid.UUID
		channelID uuid.UUID
		wantNil   bool
	}{
		{
			name:      "Get follow",
			chatID:    newChat.ID,
			channelID: newChannel.ID,
			wantNil:   false,
		},
		{
			name:      "Should fail if follow does not exist",
			chatID:    uuid.New(),
			channelID: uuid.New(),
			wantNil:   true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			f, err := service.GetByChatAndChannel(ctx, tt.channelID, tt.chatID)
			assert.NoError(t, err)

			if tt.wantNil {
				assert.Nil(t, f)
			} else {
				assert.Equal(t, newChannel.ID, f.ChannelID, "Expects channel_id to be equal.")
				assert.Equal(t, newChat.ID, f.ChatID, "Expects chat_id to be equal.")
			}
		})
	}
}

func TestFollowService_GetByChannelID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatEntRepository(entClient)
	channelsService := NewChannelEntService(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", db_models.ChatServiceTelegram)
	assert.NoError(t, err)

	channelsIds := make([]uuid.UUID, 0)

	for i := 0; i < 5; i++ {
		ch, err := channelsService.Create(
			ctx,
			fmt.Sprintf("%v", i),
			db_models.ChannelServiceTwitch,
		)
		assert.NoError(t, err)
		channelsIds = append(channelsIds, ch.ID)
	}

	for _, channelID := range channelsIds {
		f, err := service.Create(ctx, channelID, newChat.ID)
		assert.NoError(t, err)
		assert.Equal(t, channelID, f.ChannelID, "Expects channel_id to be equal.")
		assert.Equal(t, newChat.ID, f.ChatID, "Expects chat_id to be equal.")
	}

	for _, channelID := range channelsIds {
		follows, err := service.GetByChannelID(ctx, channelID)
		assert.NoError(t, err)

		for _, foll := range follows {
			assert.Equal(t, channelID, foll.ChannelID, "Expects channel_id to be equal.")
			assert.Equal(t, newChat.ID, foll.ChatID, "Expects chat_id to be equal.")
		}
	}
}

func TestFollowService_GetByChatID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	chService := NewChatEntRepository(entClient)
	channelsService := NewChannelEntService(entClient)
	service := NewFollowService(entClient)

	newChat, err := chService.Create(ctx, "1", db_models.ChatServiceTelegram)
	assert.NoError(t, err)

	channelsIds := make([]uuid.UUID, 0)

	for i := 0; i < 5; i++ {
		ch, err := channelsService.Create(
			ctx,
			fmt.Sprintf("%v", i),
			db_models.ChannelServiceTwitch,
		)
		assert.NoError(t, err)
		channelsIds = append(channelsIds, ch.ID)
	}

	for _, channelID := range channelsIds {
		f, err := service.Create(ctx, channelID, newChat.ID)
		assert.NoError(t, err)
		assert.Equal(t, channelID, f.ChannelID, "Expects channel_id to be equal.")
		assert.Equal(t, newChat.ID, f.ChatID, "Expects chat_id to be equal.")
	}

	follows, err := service.GetByChatID(ctx, newChat.ID, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, follows, 5)

	for _, foll := range follows {
		assert.Equal(t, newChat.ID, foll.ChatID, "Expects chat_id to be equal.")
	}

	followsPaginated, err := service.GetByChatID(ctx, newChat.ID, 0, 2)
	assert.NoError(t, err)
	assert.Len(t, followsPaginated, 3)
}
