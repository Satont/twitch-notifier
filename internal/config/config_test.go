package config

import (
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
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

	testCases := []struct {
		name     string
		setupEnv func(t *testing.T) (*Config, error)
		checkEnv func(t *testing.T, config *Config, err error)
	}{
		{
			name: "OK",
			setupEnv: func(t *testing.T) (*Config, error) {
				file, err := os.CreateTemp("", "temp-env")
				assert.NoError(t, err)

				filepath := file.Name()

				_, err = file.Write([]byte(strConfig))
				assert.NoError(t, err)

				defer file.Close()
				defer os.Remove(filepath)

				config, err := NewConfig(&filepath)

				return config, err
			},
			checkEnv: func(t *testing.T, config *Config, err error) {
				assert.NoError(t, err)

				assert.Equal(t, "1", config.TwitchClientId)
				assert.Equal(t, "2", config.TwitchClientSecret)
				assert.Equal(t, "3", config.TelegramToken)
				assert.IsType(t, []string{}, config.TelegramBotAdmins)
				assert.Contains(t, config.TelegramBotAdmins, "4")
				assert.Equal(t, "5", config.DatabaseUrl)
			},
		},
		{
			name: "os.Getwd() provides some error",
			setupEnv: func(t *testing.T) (*Config, error) {
				getWd = func() (string, error) {
					return "", os.ErrNotExist
				}
				defer func() { getWd = os.Getwd }()

				config, err := NewConfig(nil)

				return config, err
			},
			checkEnv: func(t *testing.T, config *Config, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, os.ErrNotExist)
				assert.Nil(t, config)
			},
		},
		{
			name: "envconfig.Process() provides some error",
			setupEnv: func(t *testing.T) (*Config, error) {
				processEnv = func(s string, i interface{}) error {
					return os.ErrNotExist
				}
				defer func() { processEnv = envconfig.Process }()

				config, err := NewConfig(nil)

				return config, err
			},
			checkEnv: func(t *testing.T, config *Config, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, os.ErrNotExist)
				assert.Nil(t, config)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t)

			cfg, err := tt.setupEnv(t)
			tt.checkEnv(t, cfg, err)
		})
	}
}
