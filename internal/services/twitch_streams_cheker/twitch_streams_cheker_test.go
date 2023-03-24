package twitch_streams_cheker

import (
	"context"
	"fmt"
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
	followMock := &db.FollowMock{}

	ctx := context.Background()

	dbChannel := &db_models.Channel{ID: uuid.New(), ChannelID: "1"}
	dbStream := &db_models.Stream{ID: "123", Titles: []string{"title"}, Categories: []string{"Dota 2"}}
	dbChat := &db_models.Chat{ID: uuid.New(), ChatID: "1"}
	dbFollow := &db_models.Follow{ID: uuid.New(), ChatID: dbChat.ID, Chat: dbChat, Channel: dbChannel, ChannelID: dbChannel.ID}
	twitchChannelInfo := &helix.ChannelInformation{BroadcasterID: "1", BroadcasterName: "Satont"}
	twitchStream := &helix.Stream{
		ID:       "123",
		GameName: "Dota 2",
		Title:    "title",
		UserID:   "1",
	}

	table := []struct {
		name       string
		setupMocks func()
	}{
		{
			name: "stream becomes offline, should call UpdateOneByStreamID with correct args",
			setupMocks: func() {
				twitchMock.On("GetChannelsByUserIds", []string{"1"}).Return([]helix.ChannelInformation{
					*twitchChannelInfo,
				}, nil)
				channelsMock.On("GetAll", ctx).Return([]*db_models.Channel{
					dbChannel,
				}, nil)
				followMock.On("GetByChannelID", ctx, dbChannel.ID).Return([]*db_models.Follow{dbFollow}, nil)
				twitchMock.On("GetStreamsByUserIds", []string{"1"}).Return([]helix.Stream{}, nil)
				streamMock.On("GetLatestByChannelID", ctx, dbChannel.ID).Return(dbStream, nil)
				followMock.On("GetByChannelID", ctx, dbChannel.ID).Return([]*db_models.Follow{dbFollow}, nil)
				streamMock.On("UpdateOneByStreamID", ctx, dbStream.ID, &db.StreamUpdateQuery{
					IsLive: lo.ToPtr(false),
				}).Return((*db_models.Stream)(nil), nil)
				senderMock.
					On("SendMessage",
						ctx,
						dbChat,
						&message_sender.MessageOpts{
							Text: fmt.Sprintf("Stream of %s is offline", twitchChannelInfo.BroadcasterName),
						},
					).
					Return(nil)
			},
		},
		{
			name: "stream becomes online, should call CreateOneByChannelID with correct args",
			setupMocks: func() {
				twitchMock.On("GetChannelsByUserIds", []string{"1"}).Return([]helix.ChannelInformation{
					*twitchChannelInfo,
				}, nil)
				channelsMock.On("GetAll", ctx).Return([]*db_models.Channel{
					dbChannel,
				}, nil)
				followMock.On("GetByChannelID", ctx, dbChannel.ID).Return([]*db_models.Follow{dbFollow}, nil)
				twitchMock.On("GetStreamsByUserIds", []string{"1"}).Return([]helix.Stream{
					*twitchStream,
				}, nil)
				streamMock.On("GetLatestByChannelID", ctx, dbChannel.ID).Return((*db_models.Stream)(nil), nil)
				streamMock.On("CreateOneByChannelID", ctx, dbChannel.ID, &db.StreamUpdateQuery{
					StreamID: "123",
					IsLive:   lo.ToPtr(true),
					Category: lo.ToPtr("Dota 2"),
					Title:    lo.ToPtr("title"),
				}).Return((*db_models.Stream)(nil), nil)
				senderMock.
					On("SendMessage",
						ctx,
						dbChat,
						&message_sender.MessageOpts{
							Text: fmt.Sprintf("Stream of %s is online", twitchChannelInfo.BroadcasterName),
						},
					).
					Return(nil)
			},
		},
		{
			name: "stream is still online, we should update category",
			setupMocks: func() {
				newHelixStream := &helix.Stream{
					ID:       "123",
					GameName: "Just Chatting",
					Title:    "title",
					UserID:   "1",
				}

				twitchMock.On("GetChannelsByUserIds", []string{"1"}).Return([]helix.ChannelInformation{
					*twitchChannelInfo,
				}, nil)
				channelsMock.On("GetAll", ctx).Return([]*db_models.Channel{
					dbChannel,
				}, nil)
				followMock.On("GetByChannelID", ctx, dbChannel.ID).Return([]*db_models.Follow{dbFollow}, nil)
				twitchMock.On("GetStreamsByUserIds", []string{"1"}).Return([]helix.Stream{
					*newHelixStream,
				}, nil)
				streamMock.On("GetLatestByChannelID", ctx, dbChannel.ID).Return(dbStream, nil)
				streamMock.On("UpdateOneByStreamID", ctx, dbStream.ID, &db.StreamUpdateQuery{
					Category: lo.ToPtr("Just Chatting"),
				}).Return((*db_models.Stream)(nil), nil)
				senderMock.
					On("SendMessage",
						ctx,
						dbChat,
						&message_sender.MessageOpts{
							Text: fmt.Sprintf("Changed category %s", twitchChannelInfo.BroadcasterName),
						},
					).
					Return(nil)
			},
		},
		{
			name: "stream is still online, we should update title",
			setupMocks: func() {
				newHelixStream := &helix.Stream{
					ID:       "123",
					GameName: "Dota 2",
					Title:    "title1",
					UserID:   "1",
				}

				twitchMock.On("GetChannelsByUserIds", []string{"1"}).Return([]helix.ChannelInformation{
					*twitchChannelInfo,
				}, nil)
				channelsMock.On("GetAll", ctx).Return([]*db_models.Channel{
					dbChannel,
				}, nil)
				followMock.On("GetByChannelID", ctx, dbChannel.ID).Return([]*db_models.Follow{dbFollow}, nil)
				twitchMock.On("GetStreamsByUserIds", []string{"1"}).Return([]helix.Stream{
					*newHelixStream,
				}, nil)
				streamMock.On("GetLatestByChannelID", ctx, dbChannel.ID).Return(dbStream, nil)
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
					Follow:  followMock,
				},
				sender: senderMock,
			}

			checker.check(ctx)

			channelsMock.AssertExpectations(t)
			twitchMock.AssertExpectations(t)
			senderMock.AssertExpectations(t)
			streamMock.AssertExpectations(t)
			followMock.AssertExpectations(t)

			channelsMock.ExpectedCalls = nil
			twitchMock.ExpectedCalls = nil
			streamMock.ExpectedCalls = nil
			senderMock.ExpectedCalls = nil
		})
	}
}
