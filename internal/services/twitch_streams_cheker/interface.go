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

func (t *twitchStreamChecker) StartPolling(ctx context.Context) {
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

func (t *twitchStreamChecker) Check() {
	return
}
