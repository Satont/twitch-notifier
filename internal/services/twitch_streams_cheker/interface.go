package twitch_streams_cheker

import (
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"time"
)

type twitchStreamCheker struct {
	twitch twitch.Interface
}

func NewTwitchStreamCheker(twitch twitch.Interface) *twitchStreamCheker {
	checker := &twitchStreamCheker{
		twitch: twitch,
	}

	return checker
}

func (t *twitchStreamCheker) StartPolling() {
	go func() {
		for {
			t.Check()

			time.Sleep(1 * time.Minute)
		}
	}()

	return
}

func (t *twitchStreamCheker) Check() {
	return
}
