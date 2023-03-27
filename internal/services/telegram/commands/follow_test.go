package commands

import (
	"context"
	"errors"
	"fmt"
	i18nmocks "github.com/satont/twitch-notifier/pkg/i18n/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/internal/test_utils"
	"github.com/satont/twitch-notifier/internal/test_utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFollowService(t *testing.T) {
	t.Parallel()

	mockedTwitch := &mocks.TwitchApiMock{}
	channelsMock := &mocks.DbChannelMock{}
	followsMock := &mocks.DbFollowMock{}

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

	sessionService := tg_types.NewMockedSessionManager()

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

func TestFollowCommand_HandleScene(t *testing.T) {
	t.Parallel()

	mockedTwitch := &mocks.TwitchApiMock{}
	channelsMock := &mocks.DbChannelMock{}
	followsMock := &mocks.DbFollowMock{}
	i18nMock := i18nmocks.NewI18nMock()
	sessionMock := tg_types.NewMockedSessionManager()

	ctx := context.Background()

	userLogin := "satont"
	helixUser := &helix.User{
		ID:          "1",
		Login:       userLogin,
		DisplayName: "Satont",
	}

	dbChat := &db_models.Chat{
		ID: uuid.New(),
		Settings: &db_models.ChatSettings{
			ChatLanguage: db_models.ChatLanguageEn,
		},
	}
	dbChannel := &db_models.Channel{
		ID:        uuid.New(),
		ChannelID: "1",
	}

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
	defer server.Close()
	tgMockedServer := test_utils.NewTelegramClient(server)

	tgMsg := &tgb.MessageUpdate{
		Client: tgMockedServer,
		Message: &tg.Message{
			Chat: tg.Chat{ID: 1},
			Text: userLogin,
		},
	}

	sessionMock.On("Get", ctx).Return(&tg_types.Session{
		Chat: dbChat,
	})

	table := []struct {
		name       string
		input      string
		setupMocks func()
	}{
		{
			name:  "Should fail because of GetUser error",
			input: "",
			setupMocks: func() {
				mockedTwitch.
					On("GetUser", "", userLogin).
					Return((*helix.User)(nil), nil)
				i18nMock.On(
					"Translate",
					"commands.follow.errors.streamerNotFound",
					"en",
					map[string]string{"streamer": userLogin},
				).Return("")
			},
		},
		{
			name:  "Should fail because db follow exists",
			input: userLogin,
			setupMocks: func() {
				mockedTwitch.On("GetUser", "", userLogin).Return(helixUser, nil)
				channelsMock.
					On("GetByIdOrCreate", ctx, helixUser.ID, db_models.ChannelServiceTwitch).
					Return(dbChannel, nil)
				followsMock.
					On("Create", ctx, dbChannel.ID, dbChat.ID).
					Return((*db_models.Follow)(nil), db_models.FollowAlreadyExistsError)
				i18nMock.On(
					"Translate",
					"commands.follow.alreadyFollowed",
					"en",
					map[string]string{"streamer": userLogin},
				).Return("")
			},
		},
		{
			name:  "Should fail because db channel cannot be created",
			input: userLogin,
			setupMocks: func() {
				mockedTwitch.On("GetUser", "", userLogin).Return(helixUser, nil)
				channelsMock.
					On("GetByIdOrCreate", ctx, helixUser.ID, db_models.ChannelServiceTwitch).
					Return(dbChannel, errors.New("some error"))
			},
		},
		{
			name:  "Should success",
			input: userLogin,
			setupMocks: func() {
				mockedTwitch.On("GetUser", "", userLogin).Return(helixUser, nil)
				channelsMock.
					On("GetByIdOrCreate", ctx, helixUser.ID, db_models.ChannelServiceTwitch).
					Return(dbChannel, nil)
				followsMock.
					On("Create", ctx, dbChannel.ID, dbChat.ID).
					Return((*db_models.Follow)(nil), nil)
				i18nMock.On(
					"Translate",
					"commands.follow.success",
					"en",
					map[string]string{"streamer": userLogin},
				).Return("")
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			followCommand := &FollowCommand{
				&tg_types.CommandOpts{
					SessionManager: sessionMock,
					Services: &types.Services{
						Twitch:  mockedTwitch,
						Channel: channelsMock,
						Follow:  followsMock,
						I18N:    i18nMock,
					},
				},
			}

			err := followCommand.handleScene(ctx, tgMsg)
			assert.NoError(t, err)

			i18nMock.AssertExpectations(t)
			mockedTwitch.AssertExpectations(t)
			channelsMock.AssertExpectations(t)
			followsMock.AssertExpectations(t)
			sessionMock.AssertExpectations(t)

			mockedTwitch.ExpectedCalls = nil
			channelsMock.ExpectedCalls = nil
			followsMock.ExpectedCalls = nil
			sessionMock.ExpectedCalls = nil
			i18nMock.ExpectedCalls = nil
		})
	}
}
