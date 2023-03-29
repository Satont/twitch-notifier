package commands

import (
	"context"
	"fmt"
	db_models2 "github.com/satont/twitch-notifier/internal/db/db_models"
	tg_types2 "github.com/satont/twitch-notifier/internal/telegram/types"
	"github.com/satont/twitch-notifier/internal/types"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/test_utils"
	"github.com/satont/twitch-notifier/internal/test_utils/mocks"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/stretchr/testify/assert"
)

func TestLiveCommand_GetList(t *testing.T) {
	t.Parallel()

	chat := &db_models2.Chat{
		ID:     uuid.New(),
		ChatID: "1",
	}

	ctx := context.Background()

	sessionManager := tg_types2.NewMockedSessionManager()
	sessionManager.On("Get", ctx).Return(&tg_types2.Session{
		Chat: chat,
	})

	followMock := &mocks.DbFollowMock{}
	twitchMock := &mocks.TwitchApiMock{}

	var now = func() time.Time {
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	follows := []*db_models2.Follow{
		{
			ID:        uuid.UUID{},
			ChannelID: uuid.UUID{},
			ChatID:    uuid.UUID{},
			Channel: &db_models2.Channel{
				ChannelID: "1",
			},
			Chat: nil,
		},
		{
			ID:        uuid.UUID{},
			ChannelID: uuid.UUID{},
			ChatID:    uuid.UUID{},
			Channel: &db_models2.Channel{
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
				followMock.On("GetByChatID", ctx, chat.ID, 0, 0).Return([]*db_models2.Follow{}, nil)
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
				CommandOpts: &tg_types2.CommandOpts{
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

func TestLiveCommand_HandleCommand(t *testing.T) {
	t.Parallel()

	chat := &db_models2.Chat{
		ID:     uuid.New(),
		ChatID: "1",
	}

	ctx := context.Background()

	sessionMock := tg_types2.NewMockedSessionManager()
	followMock := &mocks.DbFollowMock{}
	twitchMock := &mocks.TwitchApiMock{}

	var now = func() time.Time {
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	follows := []*db_models2.Follow{
		{
			ID:        uuid.UUID{},
			ChannelID: uuid.UUID{},
			ChatID:    uuid.UUID{},
			Channel: &db_models2.Channel{
				ChannelID: "1",
			},
			Chat: nil,
		},
		{
			ID:        uuid.UUID{},
			ChannelID: uuid.UUID{},
			ChatID:    uuid.UUID{},
			Channel: &db_models2.Channel{
				ChannelID: "2",
			},
			Chat: nil,
		},
	}

	sessionMock.On("Get", ctx).Return(&tg_types2.Session{
		Chat: chat,
	})
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
		{
			UserID:    "2",
			UserLogin: "sadisnamenya",
			UserName:  "SadisNaMenya",
			GameName:  "Dota 2",
			Title:     "Dotka",
			StartedAt: now(),
		},
	}, nil)

	expectedString1 := "🟢 [Satont](https://twitch.tv/satont) - 0 👁️️\n🎮 Dota 2\n📝 Playing dota\n⌛"
	expectedString2 := "🟢 [SadisNaMenya](https://twitch.tv/sadisnamenya) - 0 👁️️\n🎮 Dota 2\n📝 Dotka\n"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		query, err := url.ParseQuery(string(body))
		assert.NoError(t, err)

		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(
			t,
			fmt.Sprintf("/bot%s/sendMessage", test_utils.TelegramClientToken),
			r.URL.Path,
		)
		assert.Contains(t, query.Get("text"), expectedString1)
		assert.Contains(t, query.Get("text"), expectedString2)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
	}))
	defer server.Close()

	telegramClient := test_utils.NewTelegramClient(server)

	cmd := &LiveCommand{
		CommandOpts: &tg_types2.CommandOpts{
			SessionManager: sessionMock,
			Services: &types.Services{
				Follow: followMock,
				Twitch: twitchMock,
			},
		},
	}

	err := cmd.HandleCommand(ctx, &tgb.MessageUpdate{
		Client: telegramClient,
		Message: &tg.Message{
			Chat: tg.Chat{
				ID: 1,
			},
		},
	})
	assert.NoError(t, err)

	sessionMock.AssertExpectations(t)
	followMock.AssertExpectations(t)
	twitchMock.AssertExpectations(t)
}
