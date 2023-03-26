package twitch_streams_cheker

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/internal/test_utils/mocks"
	"github.com/satont/twitch-notifier/pkg/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTwitchStreamChecker_check(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	assert.NoError(t, err)

	channelsMock := &mocks.DbChannelMock{}
	twitchMock := &mocks.TwitchApiMock{}
	senderMock := &mocks.MessageSenderMock{}
	streamMock := &mocks.DbStreamMock{}
	followMock := &mocks.DbFollowMock{}
	i18, err := i18n.NewI18n(filepath.Join(wd, "..", "..", "..", "locales"))
	assert.NoError(t, err)

	ctx := context.Background()

	dbChannel := &db_models.Channel{ID: uuid.New(), ChannelID: "1"}
	dbStream := &db_models.Stream{ID: "123", Titles: []string{"title"}, Categories: []string{"Dota 2"}}
	dbChat := &db_models.Chat{
		ID:       uuid.New(),
		ChatID:   "1",
		Settings: &db_models.ChatSettings{ChatLanguage: db_models.ChatLanguageEn},
	}
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
						mock.Anything,
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
						mock.Anything,
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
						mock.Anything,
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
					I18N:    i18,
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
			followMock.ExpectedCalls = nil
		})
	}
}
