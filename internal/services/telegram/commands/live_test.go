package commands

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/stretchr/testify/assert"
)

func TestLiveCommand_GetList(t *testing.T) {
	t.Parallel()

	chat := &db_models.Chat{
		ID:     uuid.New(),
		ChatID: "1",
	}

	ctx := context.Background()

	sessionManager := tg_types.NewMockedSessionManager()
	sessionManager.On("Get", ctx).Return(&tg_types.Session{
		Chat: chat,
	})

	followMock := &db.FollowMock{}
	twitchMock := &twitch.Mock{}

	var now = func() time.Time {
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	follows := []*db_models.Follow{
		{
			ID:        uuid.UUID{},
			ChannelID: uuid.UUID{},
			ChatID:    uuid.UUID{},
			Channel: &db_models.Channel{
				ChannelID: "1",
			},
			Chat: nil,
		},
		{
			ID:        uuid.UUID{},
			ChannelID: uuid.UUID{},
			ChatID:    uuid.UUID{},
			Channel: &db_models.Channel{
				ChannelID: "2",
			},
			Chat: nil,
		},
	}

	table := []struct {
		name       string
		setupMocks func()
		wantErr    bool
		wants      any
	}{
		{
			name: "Should return empty list if no follows",
			setupMocks: func() {
				followMock.On("GetByChatID", ctx, chat.ID, 0, 0).Return([]*db_models.Follow{}, nil)
			},
			wantErr: false,
			wants:   []*liveChannel(nil),
		},
		{
			name: "Should return empty list if no channels online",
			setupMocks: func() {
				followMock.On("GetByChatID", ctx, chat.ID, 0, 0).
					Return(follows, nil)
				twitchMock.
					On(
						"GetStreamsByUserIds",
						[]string{"1", "2"},
					).Return([]helix.Stream{}, nil)
			},
			wantErr: false,
			wants:   []*liveChannel(nil),
		},
		{
			name: "Should return one channel",
			setupMocks: func() {
				followMock.On("GetByChatID", ctx, chat.ID, 0, 0).
					Return(follows, nil)
				twitchMock.
					On(
						"GetStreamsByUserIds",
						[]string{"1", "2"},
					).Return([]helix.Stream{
					{
						UserID:    "1",
						UserLogin: "satont",
						UserName:  "Satont",
						GameName:  "Dota 2",
						Title:     "Playing dota",
						StartedAt: now(),
					},
				}, nil)
			},
			wantErr: false,
			wants: []*liveChannel{
				{
					Name:      "Satont",
					Login:     "satont",
					StartedAt: now(),
					Title:     "Playing dota",
					Category:  "Dota 2",
				},
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			command := &LiveCommand{
				CommandOpts: &tg_types.CommandOpts{
					SessionManager: sessionManager,
					Services: &types.Services{
						Follow: followMock,
						Twitch: twitchMock,
					},
				},
			}

			list, err := command.getList(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.wants, list)

			followMock.AssertExpectations(t)
			twitchMock.AssertExpectations(t)

			followMock.ExpectedCalls = nil
			twitchMock.ExpectedCalls = nil
		})
	}
}

// func Test_NewLiveCommand(t *testing.T) {
// 	t.Parallel()

// 	router := &tg_types.MockedRouter{}

// 	router.On("Message", mock.Anything, []tgb.Filter{liveCommandFilter}).Return(router)

// 	NewLiveCommand(&tg_types.CommandOpts{Router: router})

// 	router.AssertExpectations(t)
// }
