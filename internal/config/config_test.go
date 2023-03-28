package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var strConfig = `
TWITCH_CLIENTID=1
TWITCH_CLIENTSECRET=2
TELEGRAM_TOKEN=3
TELEGRAM_BOT_ADMINS=4
DATABASE_URL=5
`

func Test_NewConfig(t *testing.T) {
	t.Parallel()

	table := []struct {
		name            string
		wantErr         bool
		patchWd         bool
		patchProcessenv bool
	}{
		{
			name:    "wd error",
			wantErr: true,
			patchWd: true,
		},
		{
			name:            "process env error",
			wantErr:         true,
			patchProcessenv: true,
		},
		{
			name: "success",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			if tt.patchWd {
				getWd = func() (string, error) {
					return "", os.ErrNotExist
				}
				defer func() { getWd = os.Getwd }()
			}

			if tt.patchProcessenv {
				processEnv = func(s string, i interface{}) error {
					return os.ErrNotExist
				}
				defer func() { processEnv = envconfig.Process }()
			}

			file, err := os.CreateTemp("", "notifier-temp-env")
			assert.NoError(t, err)

			filePath := file.Name()

			_, err = file.Write([]byte(strConfig))
			assert.NoError(t, err)
			defer os.Remove(filePath)
			defer file.Close()

			config, err := NewConfig(&filePath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, "1", config.TwitchClientId)
				assert.Equal(t, "2", config.TwitchClientSecret)
				assert.Equal(t, "3", config.TelegramToken)
				assert.IsType(t, []string{}, config.TelegramBotAdmins)
				assert.Contains(t, config.TelegramBotAdmins, "4")
				assert.Equal(t, "5", config.DatabaseUrl)
			}
		})
	}
}
