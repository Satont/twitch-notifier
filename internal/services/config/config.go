package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"path/filepath"
)

type Config struct {
	TwitchClientId     string `required:"true" envconfig:"TWITCH_CLIENTID"`
	TwitchClientSecret string `required:"true" envconfig:"TWITCH_CLIENTSECRET"`
	TelegramToken      string `required:"true" envconfig:"TELEGRAM_TOKEN"`
}

func NewConfig() (*Config, error) {
	var newCfg Config

	var err error

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	envPath := filepath.Join(wd, ".env")
	_ = godotenv.Load(envPath)

	if err = envconfig.Process("", &newCfg); err != nil {
		return nil, err
	}

	return &newCfg, nil
}
