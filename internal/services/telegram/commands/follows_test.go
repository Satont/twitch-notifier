package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/mr-linch/go-tg"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFollowsCommand_newKeyboard(t *testing.T) {
	t.Parallel()

	type fields struct {
		CommandOpts *tgtypes.CommandOpts
	}
	type args struct {
		follows []*db_models.Follow
	}

	mockedTwitch := &twitch.MockedService{}

	mockedTwitch.On("GetChannelsByUserIds", []string{"1", "2", "3", "4"}).Return([]helix.ChannelInformation{
		{
			BroadcasterID:   "1",
			BroadcasterName: "first",
		},
		{
			BroadcasterID:   "2",
			BroadcasterName: "second",
		},
		{
			BroadcasterID:   "3",
			BroadcasterName: "third",
		},
		{
			BroadcasterID:   "4",
			BroadcasterName: "fourth",
		},
	}, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *tg.InlineKeyboardMarkup
		wantErr bool
	}{
		{
			name: "should return keyboard with 4 buttons",
			fields: fields{
				CommandOpts: &tgtypes.CommandOpts{
					Services: &types.Services{
						Twitch: mockedTwitch,
					},
				},
			},
			args: args{
				follows: []*db_models.Follow{
					{
						Channel: &db_models.Channel{
							ChannelID: "1",
						},
					},
					{
						Channel: &db_models.Channel{
							ChannelID: "2",
						},
					},
					{
						Channel: &db_models.Channel{
							ChannelID: "3",
						},
					},
					{
						Channel: &db_models.Channel{
							ChannelID: "4",
						},
					},
				},
			},
			want: &tg.InlineKeyboardMarkup{
				InlineKeyboard: [][]tg.InlineKeyboardButton{
					{
						{
							Text:         "first",
							CallbackData: "channels_unfollow_1",
						},
						{
							Text:         "second",
							CallbackData: "channels_unfollow_2",
						},
						{
							Text:         "third",
							CallbackData: "channels_unfollow_3",
						},
					},
					{
						{
							Text:         "fourth",
							CallbackData: "channels_unfollow_4",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &FollowsCommand{
				CommandOpts: tt.fields.CommandOpts,
			}
			got, err := c.newKeyboard(tt.args.follows)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Len(t, got.InlineKeyboard, 2)

			for rowI, row := range got.InlineKeyboard {
				assert.LessOrEqual(t, len(row), 3)

				for btnI, btn := range row {
					assert.Equal(t, tt.want.InlineKeyboard[rowI][btnI].Text, btn.Text)
					assert.Equal(t, tt.want.InlineKeyboard[rowI][btnI].CallbackData, btn.CallbackData)
				}
			}
		})
	}
}

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
	channelsMock := &db.ChannelMock{}
	//followsMock := &db.FollowMock{}

	ctx := context.Background()
	chat := &db_models.Chat{
		ID:     uuid.New(),
		ChatID: "1",
	}

	commandOpts := &tgtypes.CommandOpts{
		Services: &types.Services{
			Channel: channelsMock,
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
