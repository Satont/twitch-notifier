package db

import (
	"context"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChannelEntService_GetByIdOrCreate(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	channelService := NewChannelEntService(entClient)

	channel, err := channelService.GetByIdOrCreate(context.Background(), "123", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)
	assert.Equal(t, "123", channel.ChannelID)
	assert.Equal(t, db_models.ChannelServiceTwitch, channel.Service)
	assert.False(t, channel.IsLive)
	assert.Nil(t, channel.Title)
	assert.Nil(t, channel.Category)
	assert.Nil(t, channel.UpdatedAt)
}

func TestChannelEntService_GetByID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	channelService := NewChannelEntService(entClient)

	table := []struct {
		name          string
		channelID     string
		service       db_models.ChannelService
		wantErr       bool
		createChannel bool
	}{
		{
			name:      "channel not found",
			channelID: "123",
			service:   db_models.ChannelServiceTwitch,
			wantErr:   true,
		},
		{
			name:          "channel found",
			channelID:     "321",
			service:       db_models.ChannelServiceTwitch,
			wantErr:       false,
			createChannel: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createChannel {
				_, err := channelService.Create(context.Background(), tt.channelID, tt.service)
				assert.NoError(t, err)
			}

			channel, err := channelService.GetByID(context.Background(), tt.channelID, tt.service)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, db_models.ChannelNotFoundError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.channelID, channel.ChannelID)
				assert.Equal(t, tt.service, channel.Service)
			}
		})
	}
}

func TestChannelEntService_Create(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	channelService := NewChannelEntService(entClient)

	table := []struct {
		name    string
		channel string
		service db_models.ChannelService
		wantErr bool
	}{
		{
			name:    "channel should be created",
			channel: "123",
			service: db_models.ChannelServiceTwitch,
		},
		{
			name:    "should fail create because channel exists",
			channel: "123",
			service: db_models.ChannelServiceTwitch,
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			channel, err := channelService.Create(context.Background(), tt.channel, tt.service)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.channel, channel.ChannelID)
				assert.Equal(t, tt.service, channel.Service)
				assert.False(t, channel.IsLive)
				assert.Nil(t, channel.Title)
				assert.Nil(t, channel.Category)
				assert.Nil(t, channel.UpdatedAt)
			}
		})
	}
}

func TestChannelEntService_Update(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	channelService := NewChannelEntService(entClient)

	table := []struct {
		name          string
		channelID     string
		service       db_models.ChannelService
		wantErr       bool
		createChannel bool
	}{
		{
			name:          "channel should be update",
			channelID:     "123",
			service:       db_models.ChannelServiceTwitch,
			createChannel: true,
		},
		{
			name:      "should fail update because channel not exists",
			channelID: "321",
			service:   db_models.ChannelServiceTwitch,
			wantErr:   true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createChannel {
				_, err := channelService.Create(context.Background(), tt.channelID, tt.service)
				assert.NoError(t, err)
			}

			channel, err := channelService.Update(
				context.Background(),
				tt.channelID,
				tt.service,
				&ChannelUpdateQuery{
					IsLive:   lo.ToPtr(true),
					Category: lo.ToPtr("Category"),
					Title:    lo.ToPtr("Title"),
				},
			)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.channelID, channel.ChannelID)
				assert.Equal(t, tt.service, channel.Service)
				assert.True(t, channel.IsLive)
				assert.Equal(t, "Title", *channel.Title)
				assert.Equal(t, "Category", *channel.Category)
				assert.NotNil(t, channel.UpdatedAt)
			}
		})
	}
}
