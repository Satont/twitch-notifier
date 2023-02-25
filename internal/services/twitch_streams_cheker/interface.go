package twitch_streams_cheker

import (
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"time"
)

type twitchStreamCheker struct {
	twitch twitch.TwitchService
}

func NewTwitchStreamCheker(twitch twitch.TwitchService) *twitchStreamCheker {
	checker := &twitchStreamCheker{
		twitch: twitch,
	}

	go func() {
		for {
			checker.Check()

			time.Sleep(1 * time.Minute)
		}
	}()

	return checker
}

func (t *twitchStreamCheker) Check() {
	return
}
