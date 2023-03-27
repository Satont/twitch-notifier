package commands

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/test_utils"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/internal/test_utils/mocks"
	"github.com/stretchr/testify/assert"
)

//func TestFollowsCommand_newKeyboard(t *testing.T) {
//	t.Parallel()
//
//	type fields struct {
//		CommandOpts *tgtypes.CommandOpts
//	}
//
//	type args struct {
//		maxRows int
//		perRow  int
//	}
//
//	ctx := context.Background()
//
//	chat := &db_models.Chat{
//		ID:     uuid.New(),
//		ChatID: "1",
//	}
//
//	sessionManager := tgtypes.NewMockedSessionManager()
//	mockedTwitch := &twitch.MockedService{}
//	mockedFollow := &db.FollowMock{}
//
//	tests := []struct {
//		name       string
//		fields     fields
//		args       args
//		want       *tg.InlineKeyboardMarkup
//		wantErr    bool
//		setupMocks func()
//	}{
//		{
//			name: "should return keyboard with 4 buttons",
//			fields: fields{
//				CommandOpts: &tgtypes.CommandOpts{
//					SessionManager: sessionManager,
//					Services: &types.Services{
//						Twitch: mockedTwitch,
//						Follow: mockedFollow,
//					},
//				},
//			},
//			args: args{
//				maxRows: 5,
//				perRow:  3,
//			},
//			want: &tg.InlineKeyboardMarkup{
//				InlineKeyboard: [][]tg.InlineKeyboardButton{
//					{
//						{
//							Text:         "first",
//							CallbackData: "channels_unfollow_1",
//						},
//						{
//							Text:         "second",
//							CallbackData: "channels_unfollow_2",
//						},
//						{
//							Text:         "third",
//							CallbackData: "channels_unfollow_3",
//						},
//					},
//					{
//						{
//							Text:         "fourth",
//							CallbackData: "channels_unfollow_4",
//						},
//					},
//				},
//			},
//			wantErr: false,
//			setupMocks: func() {
//				mockedTwitch.
//					On("GetChannelsByUserIds", []string{"1", "2", "3", "4"}).
//					Return([]helix.ChannelInformation{
//						{
//							BroadcasterID:   "1",
//							BroadcasterName: "first",
//						},
//						{
//							BroadcasterID:   "2",
//							BroadcasterName: "second",
//						},
//						{
//							BroadcasterID:   "3",
//							BroadcasterName: "third",
//						},
//						{
//							BroadcasterID:   "4",
//							BroadcasterName: "fourth",
//						},
//					}, nil)
//
//				sessionManager.On("Get", ctx).Return(&tgtypes.Session{
//					Chat: chat,
//					FollowsMenu: &tgtypes.Menu{
//						CurrentPage: 1,
//						TotalPages:  0,
//					},
//				})
//
//				mockedFollow.On("CountByChatID", ctx, chat.ID).Return(
//					4,
//					nil,
//				)
//
//				mockedFollow.On("GetByChatID", ctx, chat.ID, followsMaxRows*followsPerRow, 0).Return(
//					[]*db_models.Follow{
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "1",
//							},
//						},
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "2",
//							},
//						},
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "3",
//							},
//						},
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "4",
//							},
//						},
//					},
//					nil,
//				)
//			},
//		},
//		{
//			name: "test pagination with next button",
//			fields: fields{
//				CommandOpts: &tgtypes.CommandOpts{
//					SessionManager: sessionManager,
//					Services: &types.Services{
//						Twitch: mockedTwitch,
//						Follow: mockedFollow,
//					},
//				},
//			},
//			args: args{
//				maxRows: 1,
//				perRow:  2,
//			},
//			want: &tg.InlineKeyboardMarkup{
//				InlineKeyboard: [][]tg.InlineKeyboardButton{
//					{
//						{
//							Text:         "third",
//							CallbackData: "channels_unfollow_3",
//						},
//						{
//							Text:         "fourth",
//							CallbackData: "channels_unfollow_4",
//						},
//					},
//					{
//						{
//							Text:         "»",
//							CallbackData: "channels_unfollow_next_page",
//						},
//					},
//				},
//			},
//			wantErr: false,
//			setupMocks: func() {
//				mockedTwitch.
//					On("GetChannelsByUserIds", []string{"3", "4"}).
//					Return([]helix.ChannelInformation{
//						{
//							BroadcasterID:   "3",
//							BroadcasterName: "third",
//						},
//						{
//							BroadcasterID:   "4",
//							BroadcasterName: "fourth",
//						},
//					}, nil)
//
//				sessionManager.On("Get", ctx).Return(&tgtypes.Session{
//					Chat: chat,
//					FollowsMenu: &tgtypes.Menu{
//						CurrentPage: 1,
//						TotalPages:  0,
//					},
//				})
//
//				mockedFollow.On("CountByChatID", ctx, chat.ID).Return(
//					4,
//					nil,
//				)
//
//				mockedFollow.On("GetByChatID", ctx, chat.ID, 2, 0).Return(
//					[]*db_models.Follow{
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "3",
//							},
//						},
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "4",
//							},
//						},
//					},
//					nil,
//				)
//			},
//		},
//		{
//			name: "test pagination with prev and next buttons",
//			fields: fields{
//				CommandOpts: &tgtypes.CommandOpts{
//					SessionManager: sessionManager,
//					Services: &types.Services{
//						Twitch: mockedTwitch,
//						Follow: mockedFollow,
//					},
//				},
//			},
//			args: args{
//				maxRows: 1,
//				perRow:  1,
//			},
//			want: &tg.InlineKeyboardMarkup{
//				InlineKeyboard: [][]tg.InlineKeyboardButton{
//					{
//						{
//							Text:         "third",
//							CallbackData: "channels_unfollow_3",
//						},
//						{
//							Text:         "fourth",
//							CallbackData: "channels_unfollow_4",
//						},
//					},
//					{
//						{
//							Text:         "«",
//							CallbackData: "channels_unfollow_prev_page",
//						},
//						{
//							Text:         "»",
//							CallbackData: "channels_unfollow_next_page",
//						},
//					},
//				},
//			},
//			wantErr: false,
//			setupMocks: func() {
//				mockedTwitch.
//					On("GetChannelsByUserIds", []string{"2"}).
//					Return([]helix.ChannelInformation{
//						{
//							BroadcasterID:   "2",
//							BroadcasterName: "third",
//						},
//					}, nil)
//
//				sessionManager.On("Get", ctx).Return(&tgtypes.Session{
//					Chat: chat,
//					FollowsMenu: &tgtypes.Menu{
//						CurrentPage: 2,
//						TotalPages:  4,
//					},
//				})
//
//				mockedFollow.On("CountByChatID", ctx, chat.ID).Return(
//					8,
//					nil,
//				)
//
//				mockedFollow.On("GetByChatID", ctx, chat.ID, 1, 1).Return(
//					[]*db_models.Follow{
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "2",
//							},
//						},
//					},
//					nil,
//				)
//			},
//		},
//		{
//			name: "test pagination on page 1",
//			fields: fields{
//				CommandOpts: &tgtypes.CommandOpts{
//					SessionManager: sessionManager,
//					Services: &types.Services{
//						Twitch: mockedTwitch,
//						Follow: mockedFollow,
//					},
//				},
//			},
//			args: args{
//				maxRows: 1,
//				perRow:  1,
//			},
//			want: &tg.InlineKeyboardMarkup{
//				InlineKeyboard: [][]tg.InlineKeyboardButton{
//					{
//						{
//							Text:         "third",
//							CallbackData: "channels_unfollow_3",
//						},
//					},
//					{
//						{
//							Text:         "»",
//							CallbackData: "channels_unfollow_next_page",
//						},
//					},
//				},
//			},
//			wantErr: false,
//			setupMocks: func() {
//				sessionManager.On("Get", ctx).Return(&tgtypes.Session{
//					Chat: chat,
//					FollowsMenu: &tgtypes.Menu{
//						CurrentPage: 1,
//						TotalPages:  2,
//					},
//				})
//
//				mockedFollow.On("GetByChatID", ctx, chat.ID, 1, 0).Return(
//					[]*db_models.Follow{
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "3",
//							},
//						},
//					},
//					nil,
//				)
//
//				mockedFollow.On("CountByChatID", ctx, chat.ID).Return(
//					2,
//					nil,
//				)
//
//				mockedTwitch.
//					On("GetChannelsByUserIds", []string{"3"}).
//					Return([]helix.ChannelInformation{
//						{
//							BroadcasterID:   "3",
//							BroadcasterName: "third",
//						},
//					}, nil)
//			},
//		},
//		{
//			name: "test pagination on page 2",
//			fields: fields{
//				CommandOpts: &tgtypes.CommandOpts{
//					SessionManager: sessionManager,
//					Services: &types.Services{
//						Twitch: mockedTwitch,
//						Follow: mockedFollow,
//					},
//				},
//			},
//			args: args{
//				maxRows: 1,
//				perRow:  1,
//			},
//			want: &tg.InlineKeyboardMarkup{
//				InlineKeyboard: [][]tg.InlineKeyboardButton{
//					{
//						{
//							Text:         "fourth",
//							CallbackData: "channels_unfollow_4",
//						},
//					},
//					{
//						{
//							Text:         "«",
//							CallbackData: "channels_unfollow_prev_page",
//						},
//					},
//				},
//			},
//			wantErr: false,
//			setupMocks: func() {
//				sessionManager.On("Get", ctx).Return(&tgtypes.Session{
//					Chat: chat,
//					FollowsMenu: &tgtypes.Menu{
//						CurrentPage: 2,
//						TotalPages:  2,
//					},
//				})
//
//				mockedFollow.On("GetByChatID", ctx, chat.ID, 1, 1).Return(
//					[]*db_models.Follow{
//						{
//							Channel: &db_models.Channel{
//								ChannelID: "4",
//							},
//						},
//					},
//					nil,
//				)
//
//				mockedFollow.On("CountByChatID", ctx, chat.ID).Return(
//					2,
//					nil,
//				)
//
//				mockedTwitch.
//					On("GetChannelsByUserIds", []string{"4"}).
//					Return([]helix.ChannelInformation{
//						{
//							BroadcasterID:   "4",
//							BroadcasterName: "fourth",
//						},
//					}, nil)
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &FollowsCommand{
//				CommandOpts: tt.fields.CommandOpts,
//			}
//			tt.setupMocks()
//
//			got, err := c.newKeyboard(ctx, tt.args.maxRows, tt.args.perRow)
//			if tt.wantErr {
//				assert.Error(t, err)
//				return
//			}
//
//			assert.Len(t, got.InlineKeyboard, len(tt.want.InlineKeyboard))
//
//			for rowI, row := range tt.want.InlineKeyboard {
//				assert.LessOrEqual(t, len(got.InlineKeyboard[rowI]), 3)
//				assert.Greater(t, len(got.InlineKeyboard[rowI]), 0)
//
//				for btnI, btn := range row {
//					assert.Equal(t, btn.Text, got.InlineKeyboard[rowI][btnI].Text)
//					assert.Equal(t, btn.CallbackData, got.InlineKeyboard[rowI][btnI].CallbackData)
//				}
//			}
//
//			mockedFollow.AssertExpectations(t)
//			mockedTwitch.AssertExpectations(t)
//			sessionManager.AssertExpectations(t)
//
//			mockedTwitch.ExpectedCalls = nil
//			mockedFollow.ExpectedCalls = nil
//			sessionManager.ExpectedCalls = nil
//		})
//	}
//}

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
	followsMock.On("GetByChatID", ctx, chat.ID, 15, 0).Return([]*db_models.Follow{}, nil)

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
