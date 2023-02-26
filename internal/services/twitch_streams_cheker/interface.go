package twitch_streams_cheker

import (
	"context"
	"fmt"
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

func (t *twitchStreamCheker) StartPolling(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				fmt.Println("polled")
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (t *twitchStreamCheker) Check() {
	return
}
