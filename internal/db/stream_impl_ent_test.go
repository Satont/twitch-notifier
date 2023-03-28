package db

import (
	"context"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStreamEntService_GetByID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	channelsService := NewChannelEntService(entClient)
	service := NewStreamEntService(entClient)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)
	assert.Equal(t, "1", newChannel.ChannelID, "Expects channel_id to be equal.")

	_, err = channelsService.Create(ctx, "2", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	table := []struct {
		name      string
		channelID string
		wantNil   bool
		create    bool
		streamID  string
	}{
		{
			name:      "Get stream by id",
			channelID: newChannel.ChannelID,
			wantNil:   false,
			create:    true,
			streamID:  "1",
		},
		{
			name:      "Should return nil if stream not found",
			channelID: "2",
			wantNil:   true,
			create:    false,
			streamID:  "2",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			if tt.create {
				_, err = service.CreateOneByChannelID(ctx, newChannel.ID, &StreamUpdateQuery{
					IsLive:   nil,
					Category: nil,
					Title:    nil,
					StreamID: tt.streamID,
				})
				assert.NoError(t, err)
			}

			stream, err := service.GetByID(ctx, tt.streamID)
			if tt.wantNil {
				assert.Nil(t, stream)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, newChannel.ID, stream.ChannelID, "Expects channel_id to be equal.")
				assert.Equal(t, tt.streamID, stream.ID, "Expects stream_id to be equal.")
			}
		})
	}
}

func TestStreamEntService_GetLatestByChannelID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)
	defer teardownTest(entClient)

	ctx := context.Background()

	channelsService := NewChannelEntService(entClient)
	service := NewStreamEntService(entClient)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	table := []struct {
		name           string
		channelID      string
		wantNil        bool
		wantedStreamID string
		clearTable     bool
		before         func()
	}{
		{
			name:           "Get latest stream by channel id",
			channelID:      newChannel.ChannelID,
			wantNil:        false,
			wantedStreamID: "321",
			clearTable:     true,
			before: func() {
				_, _ = service.CreateOneByChannelID(ctx, newChannel.ID, &StreamUpdateQuery{
					StreamID: "321",
					IsLive:   lo.ToPtr(true),
					Category: lo.ToPtr("Category"),
					Title:    lo.ToPtr("Title"),
				})
			},
		},
		{
			name:           "Should return nil if stream not found",
			channelID:      newChannel.ChannelID,
			wantNil:        true,
			wantedStreamID: "2",
			clearTable:     true,
			before:         func() {},
		},
		{
			name:      "Should return correct stream",
			channelID: newChannel.ChannelID,
			wantNil:   false,
			before: func() {
				_, _ = service.CreateOneByChannelID(ctx, newChannel.ID, &StreamUpdateQuery{
					StreamID: "321",
					IsLive:   lo.ToPtr(false),
					Category: lo.ToPtr("Category"),
					Title:    lo.ToPtr("Title"),
				})
				_, _ = service.CreateOneByChannelID(ctx, newChannel.ID, &StreamUpdateQuery{
					StreamID: "4321",
					IsLive:   lo.ToPtr(true),
					Category: lo.ToPtr("Category"),
					Title:    lo.ToPtr("Title"),
				})
			},
			wantedStreamID: "4321",
			clearTable:     true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			stream, err := service.GetLatestByChannelID(ctx, newChannel.ID)
			assert.NoError(t, err)

			if tt.wantNil {
				assert.Nil(t, stream)
			} else {
				assert.NotNil(t, stream)
				assert.Equal(t, newChannel.ID, stream.ChannelID, "Expects channel_id to be equal.")
				assert.Equal(t, tt.wantedStreamID, stream.ID, "Expects stream_id to be equal.")
				assert.Nil(t, stream.EndedAt, "Expects is_live to be equal.")
				assert.Contains(t, stream.Categories, "Category", "Expects category to be equal.")
				assert.Contains(t, stream.Titles, "Title", "Expects title to be equal.")
			}

			if tt.clearTable {
				_, err = entClient.Stream.Delete().Exec(ctx)
				assert.NoError(t, err)
			}
		})
	}

}

func TestStreamEntService_GetManyByChannelID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)

	defer teardownTest(entClient)

	ctx := context.Background()

	channelsService := NewChannelEntService(entClient)
	service := NewStreamEntService(entClient)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)

	assert.NoError(t, err)

	_, err = service.CreateOneByChannelID(ctx, newChannel.ID, &StreamUpdateQuery{
		StreamID: "123",
		IsLive:   lo.ToPtr(true),
		Category: nil,
		Title:    nil,
	})
	assert.NoError(t, err)

	_, err = service.CreateOneByChannelID(ctx, newChannel.ID, &StreamUpdateQuery{
		StreamID: "321",
		IsLive:   lo.ToPtr(true),
		Category: lo.ToPtr("Category"),
		Title:    nil,
	})
	assert.NoError(t, err)

	streams, err := service.GetManyByChannelID(ctx, newChannel.ID, 100)
	assert.NoError(t, err)

	assert.Len(t, streams, 2, "Expects streams length to be equal.")
	assert.Equal(t, "321", streams[0].ID, "Expects stream_id to be equal.")
	assert.Contains(t, streams[0].Categories, "Category", "Expects category to be equal.")
	assert.Equal(t, "123", streams[1].ID, "Expects stream_id to be equal.")

	_, err = entClient.Stream.Delete().Exec(ctx)
	assert.NoError(t, err)

	streams, err = service.GetManyByChannelID(ctx, newChannel.ID, 100)
	assert.NoError(t, err)
	assert.Len(t, streams, 0, "Expects streams length to be equal.")
}

func TestStreamEntService_UpdateOneByStreamID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)

	defer teardownTest(entClient)

	ctx := context.Background()

	channelsService := NewChannelEntService(entClient)
	service := NewStreamEntService(entClient)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	_, err = service.CreateOneByChannelID(ctx, newChannel.ID, &StreamUpdateQuery{
		StreamID: "123",
		IsLive:   lo.ToPtr(true),
		Category: nil,
		Title:    nil,
	})
	assert.NoError(t, err)

	newStream, err := service.UpdateOneByStreamID(ctx, "123", &StreamUpdateQuery{
		IsLive:   lo.ToPtr(false),
		Title:    lo.ToPtr("Title"),
		Category: lo.ToPtr("Category"),
	})
	assert.NoError(t, err)

	assert.Equal(t, "123", newStream.ID, "Expects stream_id to be equal.")
	assert.Equal(t, newChannel.ID, newStream.ChannelID, "Expects channel_id to be equal.")
	assert.Equal(t, "Title", newStream.Titles[0], "Expects title to be equal.")
	assert.Equal(t, "Category", newStream.Categories[0], "Expects category to be equal.")
	assert.NotNil(t, newStream.EndedAt, "Expects ended_at to be not nil.")

	stream, err := service.UpdateOneByStreamID(ctx, "321", &StreamUpdateQuery{})
	assert.Error(t, err)
	assert.Nil(t, stream)
}

func TestStreamEntService_CreateOneByChannelID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	assert.NoError(t, err)

	defer teardownTest(entClient)

	ctx := context.Background()

	channelsService := NewChannelEntService(entClient)
	service := NewStreamEntService(entClient)

	newChannel, err := channelsService.Create(ctx, "1", db_models.ChannelServiceTwitch)
	assert.NoError(t, err)

	newStream, err := service.CreateOneByChannelID(ctx, newChannel.ID, &StreamUpdateQuery{
		StreamID: "123",
		IsLive:   lo.ToPtr(true),
		Category: nil,
		Title:    nil,
	})
	assert.NoError(t, err)

	assert.Equal(t, "123", newStream.ID, "Expects stream_id to be equal.")
	assert.Equal(t, newChannel.ID, newStream.ChannelID, "Expects channel_id to be equal.")
	assert.Nil(t, newStream.EndedAt, "Expects ended_at to be nil.")
	assert.NotNil(t, newStream.StartedAt, "Expects started_at to be not nil.")
	assert.Len(t, newStream.Categories, 0, "Expects categories length to be equal.")
}
