package commands

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/test_utils"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/internal/test_utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFollowsCommand_handleUnfollow(t *testing.T) {
	t.Parallel()

	type fields struct {
		CommandOpts *tgtypes.CommandOpts
	}
	type args struct {
		ctx   context.Context
		chat  *db_models.Chat
		input string
	}

	//mockedTwitch := &twitch.MockedService{}
	channelsMock := &mocks.DbChannelMock{}
	followsMock := &mocks.DbFollowMock{}

	ctx := context.Background()
	chat := &db_models.Chat{
		ID:     uuid.New(),
		ChatID: "1",
	}

	commandOpts := &tgtypes.CommandOpts{
		Services: &types.Services{
			Channel: channelsMock,
			Follow:  followsMock,
		},
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantedErr  error
		setupMocks func()
	}{
		{
			name:   "should return error if channel not found",
			fields: fields{CommandOpts: commandOpts},
			args: args{
				ctx:   ctx,
				chat:  chat,
				input: "channels_unfollow_1",
			},
			wantErr:   true,
			wantedErr: db_models.ChannelNotFoundError,
			setupMocks: func() {
				channelsMock.
					On("GetByID", ctx, "1", db_models.ChannelServiceTwitch).
					Return((*db_models.Channel)(nil), db_models.ChannelNotFoundError)
			},
		},
		{
			name: "should return error if follow not found",
			fields: fields{
				CommandOpts: commandOpts,
			},
			args: args{
				ctx:   ctx,
				chat:  chat,
				input: "channels_unfollow_1",
			},
			wantErr:   true,
			wantedErr: db_models.FollowNotFoundError,
			setupMocks: func() {
				channelId := uuid.New()
				channelsMock.
					On("GetByID", ctx, "1", db_models.ChannelServiceTwitch).
					Return(&db_models.Channel{
						ID:        channelId,
						ChannelID: "1",
					}, nil)
				followsMock.
					On("GetByChatAndChannel", ctx, channelId, chat.ID).
					Return((*db_models.Follow)(nil), db_models.FollowNotFoundError)
			},
		},
		{
			name: "should return nil",
			fields: fields{
				CommandOpts: commandOpts,
			},
			args: args{
				ctx:   ctx,
				chat:  chat,
				input: "channels_unfollow_1",
			},
			wantErr:   false,
			wantedErr: nil,
			setupMocks: func() {
				channelID := uuid.New()
				followID := uuid.New()
				channelsMock.
					On("GetByID", ctx, "1", db_models.ChannelServiceTwitch).
					Return(&db_models.Channel{
						ID:        channelID,
						ChannelID: "1",
					}, nil)
				followsMock.
					On("GetByChatAndChannel", ctx, channelID, chat.ID).
					Return(&db_models.Follow{
						ID: followID,
					}, nil)
				followsMock.
					On("Delete", ctx, followID).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &FollowsCommand{
				CommandOpts: tt.fields.CommandOpts,
			}

			tt.setupMocks()

			err := c.handleUnfollow(tt.args.ctx, tt.args.chat, tt.args.input)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.wantedErr)
			}

			channelsMock.AssertExpectations(t)

			channelsMock.ExpectedCalls = nil
		})
	}
}

func TestFollowsCommand_HandleCommand(t *testing.T) {
	t.Parallel()

	sessionMock := tgtypes.NewMockedSessionManager()
	followsMock := &mocks.DbFollowMock{}

	ctx := context.Background()
	chat := &db_models.Chat{
		ID:     uuid.New(),
		ChatID: "1",
		Settings: &db_models.ChatSettings{
			ChatLanguage: db_models.ChatLanguageEn,
		},
	}

	session := &tgtypes.Session{
		Chat: chat,
		FollowsMenu: &tgtypes.Menu{
			CurrentPage: 5,
			TotalPages:  10,
		},
	}

	sessionMock.On("Get", ctx).Return(session)
	followsMock.On("GetByChatID", ctx, chat.ID, 9, 0).Return([]*db_models.Follow{}, nil)

	commandOpts := &tgtypes.CommandOpts{
		Services: &types.Services{
			Follow: followsMock,
		},
		SessionManager: sessionMock,
	}

	cmd := &FollowsCommand{CommandOpts: commandOpts}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		query, _ := url.ParseQuery(string(body))

		assert.Greater(t, len(query.Get("text")), 1)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
	}))
	defer server.Close()

	msg := &tgb.MessageUpdate{
		Client: test_utils.NewTelegramClient(server),
		Message: &tg.Message{
			Chat: tg.Chat{ID: 1},
		},
	}

	err := cmd.HandleCommand(ctx, msg)
	assert.NoError(t, err)
}

func TestFollowsCommand_newKeyboard(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	followsMock := &mocks.DbFollowMock{}
	sessionsMock := tgtypes.NewMockedSessionManager()
	twitchMock := &mocks.TwitchApiMock{}

	dbChat := &db_models.Chat{
		ID:     uuid.New(),
		ChatID: "1",
		Settings: &db_models.ChatSettings{
			ChatLanguage: db_models.ChatLanguageEn,
		},
	}

	session := &tgtypes.Session{
		Chat: dbChat,
		FollowsMenu: &tgtypes.Menu{
			CurrentPage: 1,
			TotalPages:  0,
		},
	}

	table := []struct {
		name       string
		setupMocks func()
		asserts    func(t *testing.T, keyboard *tg.InlineKeyboardMarkup)
	}{
		{
			name: "should return keyboard with 1 page and no next and prev buttons",
			setupMocks: func() {
				sessionsMock.On("Get", ctx).Return(session)
				followsMock.On("GetByChatID", ctx, dbChat.ID, 9, 0).
					Return([]*db_models.Follow{
						{
							ID:        uuid.New(),
							ChannelID: uuid.New(),
							ChatID:    dbChat.ID,
							Channel: &db_models.Channel{
								ID:        uuid.New(),
								ChannelID: "1",
							},
						},
					}, nil)
				followsMock.On("CountByChatID", ctx, dbChat.ID).
					Return(1, nil)
				twitchMock.On("GetChannelsByUserIds", []string{"1"}).
					Return([]helix.ChannelInformation{
						{BroadcasterID: "1", BroadcasterName: "Satont"},
					}, nil)
			},
			asserts: func(t *testing.T, keyboard *tg.InlineKeyboardMarkup) {
				assert.Len(t, keyboard.InlineKeyboard, 1)
				assert.Len(t, keyboard.InlineKeyboard[0], 1)
				assert.Equal(t, keyboard.InlineKeyboard[0][0].Text, "Satont")
				assert.Equal(t, keyboard.InlineKeyboard[0][0].CallbackData, "channels_unfollow_1")
			},
		},
		{
			name: "should return keyboard with 2 pages and next buttons",
			setupMocks: func() {
				sessionsMock.On("Get", ctx).Return(session)
				follows := make([]*db_models.Follow, 0, 20)
				for i := 0; i < 20; i++ {
					follows = append(follows, &db_models.Follow{
						ID:        uuid.New(),
						ChannelID: uuid.New(),
						ChatID:    dbChat.ID,
						Channel: &db_models.Channel{
							ID:        uuid.New(),
							ChannelID: strconv.Itoa(i),
						},
					})
				}
				followsMock.On("GetByChatID", ctx, dbChat.ID, 9, 0).
					Return(follows, nil)
				followsMock.On("CountByChatID", ctx, dbChat.ID).
					Return(len(follows), nil)
				channelsIds := lo.Map(follows, func(f *db_models.Follow, _ int) string {
					return f.Channel.ChannelID
				})
				twitchMock.On("GetChannelsByUserIds", channelsIds).
					Return(lo.Map(follows, func(item *db_models.Follow, _ int) helix.ChannelInformation {
						return helix.ChannelInformation{
							BroadcasterID:   item.Channel.ChannelID,
							BroadcasterName: item.Channel.ChannelID,
						}
					}), nil)
			},
			asserts: func(t *testing.T, keyboard *tg.InlineKeyboardMarkup) {
				assert.Greater(t, len(keyboard.InlineKeyboard), 2)
				assert.Contains(
					t,
					keyboard.InlineKeyboard[len(keyboard.InlineKeyboard)-1][0].CallbackData,
					"channels_unfollow_next_page",
				)
			},
		},
		{
			name: "should return keyboard with few pages and next and prev buttons",
			setupMocks: func() {
				session.FollowsMenu.CurrentPage = 3
				sessionsMock.On("Get", ctx).Return(session)
				follows := make([]*db_models.Follow, 0, 15)
				for i := 0; i < 15; i++ {
					follows = append(follows, &db_models.Follow{
						ID:        uuid.New(),
						ChannelID: uuid.New(),
						ChatID:    dbChat.ID,
						Channel: &db_models.Channel{
							ID:        uuid.New(),
							ChannelID: strconv.Itoa(i),
						},
					})
				}
				followsMock.On("GetByChatID", ctx, dbChat.ID, 9, 18).
					Return(follows, nil)
				followsMock.On("CountByChatID", ctx, dbChat.ID).
					Return(100, nil)
				channelsIds := lo.Map(follows, func(f *db_models.Follow, _ int) string {
					return f.Channel.ChannelID
				})
				twitchMock.On("GetChannelsByUserIds", channelsIds).
					Return(lo.Map(follows, func(item *db_models.Follow, _ int) helix.ChannelInformation {
						return helix.ChannelInformation{
							BroadcasterID:   item.Channel.ChannelID,
							BroadcasterName: item.Channel.ChannelID,
						}
					}), nil)
			},
			asserts: func(t *testing.T, keyboard *tg.InlineKeyboardMarkup) {
				assert.Greater(t, len(keyboard.InlineKeyboard), 2)
				latestRow := keyboard.InlineKeyboard[len(keyboard.InlineKeyboard)-1]
				assert.Equal(t, latestRow[0].CallbackData, "channels_unfollow_prev_page")
				assert.Equal(t, latestRow[1].CallbackData, "channels_unfollow_next_page")
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			cmd := &FollowsCommand{
				CommandOpts: &tgtypes.CommandOpts{
					Services: &types.Services{
						Follow: followsMock,
						Twitch: twitchMock,
					},
					SessionManager: sessionsMock,
				},
			}

			keyboard, err := cmd.newKeyboard(ctx, followsMaxRows, followsPerRow)
			assert.NoError(t, err)
			tt.asserts(t, keyboard)

			sessionsMock.AssertExpectations(t)
			followsMock.AssertExpectations(t)
			twitchMock.AssertExpectations(t)

			sessionsMock.ExpectedCalls = nil
			followsMock.ExpectedCalls = nil
			twitchMock.ExpectedCalls = nil
		})
	}
}
