package twitch_streams_cheker

import (
	"context"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/satont/twitch-notifier/internal/services/message_sender"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/internal/services/types"
	"testing"
)

func TestTwitchStreamChecker_check(t *testing.T) {
	t.Parallel()

	channelsMock := &db.ChannelMock{}
	twitchMock := &twitch.Mock{}
	senderMock := &message_sender.Mock{}
	streamMock := &db.StreamMock{}

	ctx := context.Background()

	table := []struct {
		name       string
		setupMocks func()
	}{
		{
			name: "stream becomes offline, should call UpdateOneByStreamID with correct args",
			setupMocks: func() {
				channel := &db_models.Channel{ID: uuid.New(), ChannelID: "1"}
				stream := &db_models.Stream{ID: "123"}

				channelsMock.On("GetAll", ctx).Return([]*db_models.Channel{
					channel,
				}, nil)
				twitchMock.On("GetStreamsByUserIds", []string{"1"}).Return([]helix.Stream{}, nil)
				streamMock.On("GetLatestByChannelID", ctx, channel.ID).Return(stream, nil)
				streamMock.On("UpdateOneByStreamID", ctx, stream.ID, &db.StreamUpdateQuery{
					IsLive: lo.ToPtr(false),
				}).Return((*db_models.Stream)(nil), nil)
			},
		},
		{
			name: "stream becomes online, should call CreateOneByChannelID with correct args",
			setupMocks: func() {
				channel := &db_models.Channel{ID: uuid.New(), ChannelID: "1"}
				helixStream := &helix.Stream{
					ID:       "123",
					GameName: "Dota 2",
					Title:    "title",
					UserID:   "1",
				}

				channelsMock.On("GetAll", ctx).Return([]*db_models.Channel{
					channel,
				}, nil)
				twitchMock.On("GetStreamsByUserIds", []string{"1"}).Return([]helix.Stream{
					*helixStream,
				}, nil)
				streamMock.On("GetLatestByChannelID", ctx, channel.ID).Return((*db_models.Stream)(nil), nil)
				streamMock.On("CreateOneByChannelID", ctx, channel.ID, &db.StreamUpdateQuery{
					StreamID: "123",
					IsLive:   lo.ToPtr(true),
					Category: lo.ToPtr("Dota 2"),
					Title:    lo.ToPtr("title"),
				}).Return((*db_models.Stream)(nil), nil)
			},
		},
		{
			name: "stream is still online, we should update category",
			setupMocks: func() {
				channel := &db_models.Channel{ID: uuid.New(), ChannelID: "1"}
				newHelixStream := &helix.Stream{
					ID:       "123",
					GameName: "Just Chatting",
					Title:    "title",
					UserID:   "1",
				}
				dbStream := &db_models.Stream{
					ID:         "123",
					Titles:     []string{"title"},
					Categories: []string{"Dota 2"},
				}

				channelsMock.On("GetAll", ctx).Return([]*db_models.Channel{
					channel,
				}, nil)
				twitchMock.On("GetStreamsByUserIds", []string{"1"}).Return([]helix.Stream{
					*newHelixStream,
				}, nil)
				streamMock.On("GetLatestByChannelID", ctx, channel.ID).Return(dbStream, nil)
				streamMock.On("UpdateOneByStreamID", ctx, dbStream.ID, &db.StreamUpdateQuery{
					Category: lo.ToPtr("Just Chatting"),
				}).Return((*db_models.Stream)(nil), nil)
			},
		},
		{
			name: "stream is still online, we should update title",
			setupMocks: func() {
				channel := &db_models.Channel{ID: uuid.New(), ChannelID: "1"}
				newHelixStream := &helix.Stream{
					ID:       "123",
					GameName: "Dota 2",
					Title:    "title1",
					UserID:   "1",
				}
				dbStream := &db_models.Stream{
					ID:         "123",
					Titles:     []string{"title"},
					Categories: []string{"Dota 2"},
				}

				channelsMock.On("GetAll", ctx).Return([]*db_models.Channel{
					channel,
				}, nil)
				twitchMock.On("GetStreamsByUserIds", []string{"1"}).Return([]helix.Stream{
					*newHelixStream,
				}, nil)
				streamMock.On("GetLatestByChannelID", ctx, channel.ID).Return(dbStream, nil)
				streamMock.On("UpdateOneByStreamID", ctx, dbStream.ID, &db.StreamUpdateQuery{
					Title: lo.ToPtr("title1"),
				}).Return((*db_models.Stream)(nil), nil)
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			checker := &TwitchStreamChecker{
				services: &types.Services{
					Channel: channelsMock,
					Twitch:  twitchMock,
					Stream:  streamMock,
				},
				sender: senderMock,
			}

			checker.check(ctx)

			channelsMock.AssertExpectations(t)
			twitchMock.AssertExpectations(t)
			senderMock.AssertExpectations(t)
			streamMock.AssertExpectations(t)

			channelsMock.ExpectedCalls = nil
			twitchMock.ExpectedCalls = nil
			streamMock.ExpectedCalls = nil
		})
	}
}
