package commands

import (
	"fmt"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/pkg/i18n"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLanguagePicker_buildKeyboard(t *testing.T) {
	t.Parallel()

	i18nMock := i18n.NewI18nMock()

	cmd := &LanguagePicker{
		CommandOpts: &tgtypes.CommandOpts{
			Services: &types.Services{
				I18N: i18nMock,
			},
		},
	}

	i18nMock.On("GetLanguagesCodes").Return([]string{"en", "ru"})

	englishFlag := "🇬🇧"
	englishName := "English"

	russianFlag := "🇷🇺"
	russianName := "Русский"

	i18nMock.
		On("Translate", "language.emoji", "en", map[string]string(nil)).
		Return(englishFlag)
	i18nMock.
		On("Translate", "language.name", "en", map[string]string(nil)).
		Return(englishName)

	i18nMock.
		On("Translate", "language.emoji", "ru", map[string]string(nil)).
		Return(russianFlag)
	i18nMock.
		On("Translate", "language.name", "ru", map[string]string(nil)).
		Return(russianName)

	keyboard, err := cmd.buildKeyboard()
	assert.NoError(t, err)

	assert.Equal(t,
		fmt.Sprintf("%s %s", englishFlag, englishName),
		keyboard.InlineKeyboard[0][0].Text,
	)
	assert.Equal(t,
		"language_picker_set_en",
		keyboard.InlineKeyboard[0][0].CallbackData,
	)

	assert.Equal(t,
		fmt.Sprintf("%s %s", russianFlag, russianName),
		keyboard.InlineKeyboard[1][0].Text,
	)
	assert.Equal(t,
		"language_picker_set_ru",
		keyboard.InlineKeyboard[1][0].CallbackData,
	)

	assert.Equal(t,
		"«",
		keyboard.InlineKeyboard[2][0].Text,
	)
	assert.Equal(t, "start_command_menu", keyboard.InlineKeyboard[2][0].CallbackData)

	i18nMock.AssertExpectations(t)
}
