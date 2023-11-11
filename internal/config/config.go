package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TwitchClientId     string   `required:"true"  envconfig:"TWITCH_CLIENTID"`
	TwitchClientSecret string   `required:"true"  envconfig:"TWITCH_CLIENTSECRET"`
	TelegramToken      string   `required:"true"  envconfig:"TELEGRAM_TOKEN"`
	AppEnv             string   `required:"true"  envconfig:"APP_ENV"             default:"development"`
	TelegramBotAdmins  []string `required:"false" envconfig:"TELEGRAM_BOT_ADMINS"`
	DatabaseUrl        string   `required:"true"  envconfig:"DATABASE_URL"`
	SentryDsn          string   `required:"false" envconfig:"SENTRY_DSN"`
	RedisUrl           string   `required:"true"  envconfig:"REDIS_URL"`
}

var getWd = os.Getwd
var processEnv = envconfig.Process

func NewConfig(customPath *string) (*Config, error) {
	var newCfg Config

	var err error

	wd, err := getWd()
	if err != nil {
		return nil, err
	}

	envPath := filepath.Join(wd, ".env")

	if customPath != nil {
		envPath = *customPath
	}

	_ = godotenv.Overload(envPath)
	if err = processEnv("", &newCfg); err != nil {
		return nil, err
	}

	return &newCfg, nil
}
