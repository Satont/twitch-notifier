package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/mr-linch/go-tg"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/internal/services/db"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockMessageUpdate struct {
	mock.Mock
}

func (m *MockMessageUpdate) Answer(text string) *tg.SendMessageCall {
	args := m.Called(text)
	return args.Get(0).(*tg.SendMessageCall)
}

func TestFollow(t *testing.T) {
	login := "fukushine"

	mockedTwitch := &twitch.MockedService{}
	channelsMock := &db.ChannelMock{}
	followsMock := &db.FollowMock{}

	user := &helix.User{
		ID:          "1",
		Login:       login,
		DisplayName: "Fukushine",
	}

	ctx := context.Background()

	chat := &ent.Chat{
		ID: uuid.New(),
	}
	chann := &ent.Channel{
		ID:        uuid.New(),
		ChannelID: "1",
	}
	f := &ent.Follow{}

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
		want       *ent.Follow
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
			input:   login,
			want:    f,
			wantErr: false,
			setupMocks: func() {
				mockedTwitch.On("GetUser", "", login).Return(user, nil)
				channelsMock.On("GetByIdOrCreate", ctx, user.ID, channel.ServiceTwitch).Return(chann, nil)
				followsMock.On("Create", ctx, chann.ID, chat.ID).Return(f, nil)
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
			}
			assert.Equal(t, tt.want, got)

			mockedTwitch.AssertExpectations(t)
			channelsMock.AssertExpectations(t)
			followsMock.AssertExpectations(t)
		})
	}

}
