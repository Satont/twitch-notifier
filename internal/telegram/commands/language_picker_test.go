package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	tg_types "github.com/satont/twitch-notifier/internal/telegram/types"
	"github.com/satont/twitch-notifier/internal/types"

	"github.com/google/uuid"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/test_utils"
	"github.com/satont/twitch-notifier/internal/test_utils/mocks"
	i18nmocks "github.com/satont/twitch-notifier/pkg/i18n/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//func TestLanguagePicker_buildKeyboard(t *testing.T) {
//	t.Parallel()
//
//	i18nMock := i18n.NewI18nMock()
//
//	cmd := &LanguagePicker{
//		CommandOpts: &tgtypes.CommandOpts{
//			Services: &types.Services{
//				I18N: i18nMock,
//			},
//		},
//	}
//
//	i18nMock.On("GetLanguagesCodes").Return([]string{"en", "ru"})
//
//	englishFlag := "🇬🇧"
//	englishName := "English"
//
//	russianFlag := "🇷🇺"
//	russianName := "Русский"
//
//	i18nMock.
//		On("Translate", "language.emoji", "en", map[string]string(nil)).
//		Return(englishFlag)
//	i18nMock.
//		On("Translate", "language.name", "en", map[string]string(nil)).
//		Return(englishName)
//
//	i18nMock.
//		On("Translate", "language.emoji", "ru", map[string]string(nil)).
//		Return(russianFlag)
//	i18nMock.
//		On("Translate", "language.name", "ru", map[string]string(nil)).
//		Return(russianName)
//
//	keyboard, err := cmd.buildKeyboard()
//	assert.NoError(t, err)
//
//	assert.Equal(t,
//		fmt.Sprintf("%s %s", englishFlag, englishName),
//		keyboard.InlineKeyboard[0][0].Text,
//	)
//	assert.Equal(t,
//		"language_picker_set_en",
//		keyboard.InlineKeyboard[0][0].CallbackData,
//	)
//
//	assert.Equal(t,
//		fmt.Sprintf("%s %s", russianFlag, russianName),
//		keyboard.InlineKeyboard[1][0].Text,
//	)
//	assert.Equal(t,
//		"language_picker_set_ru",
//		keyboard.InlineKeyboard[1][0].CallbackData,
//	)
//
//	assert.Equal(t,
//		"«",
//		keyboard.InlineKeyboard[2][0].Text,
//	)
//	assert.Equal(t, "start_command_menu", keyboard.InlineKeyboard[2][0].CallbackData)
//
//	i18nMock.AssertExpectations(t)
//}

func TestLanguagePicker_HandleCallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	i18nMock := i18nmocks.NewI18nMock()
	i18nMock.On("GetLanguagesCodes").Return([]string{"en"})

	englishFlag := "🇬🇧"
	englishName := "English"

	i18nMock.
		On("Translate", "language.emoji", "en", map[string]string(nil)).
		Return(englishFlag)
	i18nMock.
		On("Translate", "language.name", "en", map[string]string(nil)).
		Return(englishName)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		query, err := url.ParseQuery(string(body))
		assert.NoError(t, err)

		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(
			t,
			fmt.Sprintf("/bot%s/editMessageReplyMarkup", test_utils.TelegramClientToken),
			r.URL.Path,
		)
		assert.NotEmpty(t, query.Get("reply_markup"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
	}))

	cmd := &LanguagePicker{
		CommandOpts: &tg_types.CommandOpts{
			Services: &types.Services{
				I18N: i18nMock,
			},
		},
	}

	err := cmd.HandleCallback(ctx, &tgb.CallbackQueryUpdate{
		Client: test_utils.NewTelegramClient(server),
		CallbackQuery: &tg.CallbackQuery{
			Message: &tg.MaybeInaccessibleMessage{
				InaccessibleMessage: &tg.InaccessibleMessage{MessageID: 1, Chat: tg.Chat{ID: tg.ChatID(1)}},
			},
		},
	})
	assert.NoError(t, err)
}

func TestLanguagePicker_handleSetLanguage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	chat := &db_models.Chat{
		ID:     uuid.New(),
		ChatID: "1",
		Settings: &db_models.ChatSettings{
			ChatLanguage: db_models.ChatLanguageEn,
		},
	}

	sessionMock := tg_types.NewMockedSessionManager()
	sessionMock.On("Get", ctx).Return(&tg_types.Session{
		Chat: chat,
	})

	i18nMock := i18nmocks.NewI18nMock()
	i18nMock.
		On("Translate", "language.changed", "ru", map[string]string(nil)).
		Return("Now russian")

	chatService := &mocks.DbChatMock{}
	chatService.
		On("Update", ctx, "1", db_models.ChatServiceTelegram, mock.IsType(&db.ChatUpdateQuery{})).
		Return((*db_models.Chat)(nil), nil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(
			t,
			fmt.Sprintf("/bot%s/answerCallbackQuery", test_utils.TelegramClientToken),
			r.URL.Path,
		)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
	}))

	cmd := &LanguagePicker{
		CommandOpts: &tg_types.CommandOpts{
			SessionManager: sessionMock,
			Services: &types.Services{
				I18N: i18nMock,
				Chat: chatService,
			},
		},
	}

	err := cmd.handleSetLanguage(ctx, &tgb.CallbackQueryUpdate{
		Client: test_utils.NewTelegramClient(server),
		CallbackQuery: &tg.CallbackQuery{
			Message: &tg.MaybeInaccessibleMessage{
				Message: &tg.Message{
					ID: 1,
					Chat: tg.Chat{
						ID: tg.ChatID(1),
					},
				},
			},
			Data: "language_picker_set_ru",
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, db_models.ChatLanguageRu, chat.Settings.ChatLanguage)

	sessionMock.AssertExpectations(t)
	chatService.AssertExpectations(t)
}
