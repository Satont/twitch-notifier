package commands

import (
	"github.com/mr-linch/go-tg/tgb"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewStartCommand(t *testing.T) {
	t.Parallel()

	router := &tg_types.MockedRouter{}

	router.On("Message", mock.Anything, []tgb.Filter{startCommandFilter}).Return(router)

	NewStartCommand(&tg_types.CommandOpts{Router: router})

	router.AssertExpectations(t)
}
