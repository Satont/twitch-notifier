package config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var config = `
TWITCH_CLIENTID=1
TWITCH_CLIENTSECRET=2
TELEGRAM_TOKEN=3
TELEGRAM_BOT_ADMINS=4
`

func Test_NewConfig(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "notifier-temp-env")
	if err != nil {
		log.Fatal(err)
	}

	filePath := file.Name()

	defer os.Remove(filePath)

	_, err = file.Write([]byte(config))
	assert.NoError(t, err)

	config, err := NewConfig(&filePath)
	assert.NoError(t, err)

	assert.Equal(t, "1", config.TwitchClientId)
	assert.Equal(t, "2", config.TwitchClientSecret)
	assert.Equal(t, "3", config.TelegramToken)
	assert.Contains(t, config.TelegramBotAdmins, "4")
}
