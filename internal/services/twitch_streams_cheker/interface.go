package twitch_streams_cheker

import (
	"context"
	"fmt"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"time"
)

type twitchStreamChecker struct {
	twitch twitch.Interface
}

func NewTwitchStreamChecker(twitch twitch.Interface) *twitchStreamChecker {
	checker := &twitchStreamChecker{
		twitch: twitch,
	}

	return checker
}

func (t *twitchStreamChecker) check(ctx context.Context) {
	return
}

func (t *twitchStreamChecker) StartPolling(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("polled")
				t.check(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (t *twitchStreamChecker) Check() {
	return
}
