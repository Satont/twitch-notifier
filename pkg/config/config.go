package config

type Config struct {
	PostgresUrl        string
	AppEnv             string
	TwitchClientID     string
	TwitchClientSecret string
	TemporalUrl        string
}

func New() (*Config, error) {
	return &Config{
		PostgresUrl:        "",
		TwitchClientID:     "your-client-id",
		TwitchClientSecret: "your-client-secret",
		TemporalUrl:        "",
	}, nil
}
