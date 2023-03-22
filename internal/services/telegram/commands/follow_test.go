package commands

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	"github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/internal/test_utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFollowService(t *testing.T) {
	t.Parallel()

	mockedTwitch := &twitch.MockedService{}
	channelsMock := &db.ChannelMock{}
	followsMock := &db.FollowMock{}

	userLogin := "fukushine"
	user := &helix.User{
		ID:          "1",
		Login:       userLogin,
		DisplayName: "Fukushine",
	}

	ctx := context.Background()

	chat := &db_models.Chat{
		ID: uuid.New(),
	}
	chann := &db_models.Channel{
		ID:        uuid.New(),
		ChannelID: "1",
	}
	f := &db_models.Follow{}

	follow := &FollowCommand{
		&tg_types.CommandOpts{
			Services: &types.Services{
				Twitch:  mockedTwitch,
				Channel: channelsMock,
				Follow:  followsMock,
			},
		},
	}

	// table tests
	table := []struct {
		name       string
		input      string
		want       *db_models.Follow
		wantErr    bool
		setupMocks func()
	}{
		{
			name:    "Should fail because of GetUser error",
			input:   "fukushine2",
			want:    nil,
			wantErr: true,
			setupMocks: func() {
				mockedTwitch.On("GetUser", "", "fukushine2").Return((*helix.User)(nil), nil)
			},
		},
		{
			name:    "Should create",
			input:   userLogin,
			want:    f,
			wantErr: false,
			setupMocks: func() {
				mockedTwitch.
					On("GetUser", "", userLogin).Return(user, nil)
				channelsMock.
					On("GetByIdOrCreate", ctx, user.ID, db_models.ChannelServiceTwitch).Return(chann, nil)
				followsMock.
					On("Create", ctx, chann.ID, chat.ID).Return(f, nil)
			},
		},
		{
			name:    "Should fail because follow exists",
			input:   userLogin,
			want:    nil,
			wantErr: true,
			setupMocks: func() {
				mockedTwitch.
					On("GetUser", "", userLogin).Return(user, nil)
				channelsMock.
					On("GetByIdOrCreate", ctx, user.ID, db_models.ChannelServiceTwitch).Return(chann, nil)
				followsMock.
					On("Create", ctx, chann.ID, chat.ID).Return((*db_models.Follow)(nil), db_models.FollowAlreadyExistsError)
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := follow.createFollow(ctx, chat, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockedTwitch.AssertExpectations(t)
			channelsMock.AssertExpectations(t)
			followsMock.AssertExpectations(t)

			mockedTwitch.ExpectedCalls = nil
			channelsMock.ExpectedCalls = nil
			followsMock.ExpectedCalls = nil
		})
	}
}

func TestFollowCommand_HandleCommand(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	sessionService := test_utils.NewMockedSessionManager()

	sessionService.On("Get", ctx).Return(&tg_types.Session{
		Chat: &db_models.Chat{ChatID: "123"},
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(
			t,
			fmt.Sprintf("/bot%s/sendMessage", test_utils.TelegramClientToken),
			r.URL.Path,
		)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
	}))

	tgClient := test_utils.NewTelegramClient(server)

	followCommand := &FollowCommand{
		&tg_types.CommandOpts{
			SessionManager: sessionService,
			Services:       &types.Services{},
		},
	}

	assert.Equal(t, "", sessionService.Get(ctx).Scene)
	err := followCommand.HandleCommand(ctx, &tgb.MessageUpdate{
		Client: tgClient,
		Message: &tg.Message{
			Chat: tg.Chat{
				ID: 123,
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "follow", sessionService.Get(ctx).Scene)

	sessionService.AssertExpectations(t)
}
