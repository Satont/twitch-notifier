package db

import (
	"context"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChannelServiceEntGetByIdOrCreate(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	channelService := NewChannelEntRepository(entClient)

	channel, err := channelService.GetByIdOrCreate(context.Background(), "123", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)
	assert.Equal(t, "123", channel.ChannelID)
	assert.Equal(t, db_models.ChannelServiceTwitch, channel.Service)
	assert.False(t, channel.IsLive)
	assert.Nil(t, channel.Title)
	assert.Nil(t, channel.Category)
	assert.Nil(t, channel.UpdatedAt)
}
