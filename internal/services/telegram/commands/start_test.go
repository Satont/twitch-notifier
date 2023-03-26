package commands

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/internal/test_utils"
	"github.com/satont/twitch-notifier/pkg/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestStartCommand_buildKeyboard(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	chat := &db_models.Chat{
		ID:     uuid.New(),
		ChatID: "1",
		Settings: &db_models.ChatSettings{
			ChatLanguage: db_models.ChatLanguageEn,
		},
	}

	i18 := i18n.NewI18nMock()
	i18.
		On("Translate", mock.Anything, mock.Anything, mock.Anything).
		Return("")

	sessionManager := tgtypes.NewMockedSessionManager()
	sessionManager.On("Get", ctx).Return(&tgtypes.Session{
		Chat: chat,
	})

	cmd := &StartCommand{
		CommandOpts: &tgtypes.CommandOpts{
			SessionManager: sessionManager,
			Services: &types.Services{
				I18N: i18,
			},
		},
	}

	keyboard, err := cmd.buildKeyboard(ctx)
	assert.NoError(t, err)

	const buttons = 4
	assert.Equal(t, buttons, len(keyboard.InlineKeyboard))

	assert.Equal(
		t,
		"start_game_change_notification_setting",
		keyboard.InlineKeyboard[0][0].CallbackData,
	)

	assert.Equal(
		t,
		"start_offline_notification",
		keyboard.InlineKeyboard[1][0].CallbackData,
	)

	assert.Equal(
		t,
		"language_picker",
		keyboard.InlineKeyboard[2][0].CallbackData,
	)

	assert.Equal(t, "Github", keyboard.InlineKeyboard[3][0].Text)
	assert.Equal(t, "https://github.com/Satont/twitch-notifier", keyboard.InlineKeyboard[3][0].URL)

	sessionManager.AssertExpectations(t)
	i18.AssertNumberOfCalls(t, "Translate", buttons-1)
}

func TestStartCommand_HandleCommand(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	chat := &db_models.Chat{
		ID: uuid.New(),
		Settings: &db_models.ChatSettings{
			ChatLanguage: db_models.ChatLanguageEn,
		},
	}

	i18 := i18n.NewI18nMock()
	i18.
		On("Translate", mock.Anything, mock.Anything, mock.Anything).
		Return("start command")

	sessionManager := tgtypes.NewMockedSessionManager()
	sessionManager.On("Get", ctx).Return(&tgtypes.Session{
		Chat: chat,
	})

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
		assert.Equal(t, "start command", query.Get("text"))
		assert.NotEmpty(t, query.Get("reply_markup"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
	}))

	cmd := &StartCommand{
		CommandOpts: &tgtypes.CommandOpts{
			SessionManager: sessionManager,
			Services: &types.Services{
				I18N: i18,
			},
		},
	}

	err := cmd.HandleCommand(ctx, &tgb.MessageUpdate{
		Client: test_utils.NewTelegramClient(server),
		Message: &tg.Message{
			Text: "/start",
		},
	})
	assert.NoError(t, err)
}
