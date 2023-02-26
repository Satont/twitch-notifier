package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twitch-notifier/ent"
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
	user := &helix.User{
		ID:          "1",
		Login:       login,
		DisplayName: "Fukushine",
	}

	ctx := context.Background()

	sessionManager := session.NewManager[tg_types.Session](tg_types.Session{})

	sessionManager.Get(ctx).Chat = &ent.Chat{
		ID: uuid.New(),
	}

	follow := &FollowCommand{
		&tg_types.CommandOpts{
			Services:       &types.Services{Twitch: mockedTwitch},
			SessionManager: sessionManager,
		},
	}

	mockedTwitch.On("GetUser", "", login).Return(user, nil)

	_, err := follow.createFollow(ctx, login)
	assert.NoError(t, err)

	mockedTwitch.AssertExpectations(t)
}
