package twitch_streams_cheker

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/services/types"
	"go.uber.org/zap"
	"time"
)

type TwitchStreamChecker struct {
	services *types.Services
	ticks    int
	tickTime *time.Duration
}

func NewTwitchStreamChecker(services *types.Services, tickTime *time.Duration) *TwitchStreamChecker {
	checker := &TwitchStreamChecker{
		services: services,
		tickTime: tickTime,
	}

	return checker
}

func (t *TwitchStreamChecker) check(ctx context.Context) {
	channels, err := t.services.Channel.GetAll(ctx)
	if err != nil {
		zap.S().Error(err)
		return
	}

	channelsIDs := make([]string, 0, len(channels))
	for _, channel := range channels {
		channelsIDs = append(channelsIDs, channel.ChannelID)
	}

	streams, err := t.services.Twitch.GetStreamsByUserIds(channelsIDs)
	if err != nil {
		zap.S().Error(err)
		return
	}

	fmt.Println(streams)

	return
}

func (t *TwitchStreamChecker) StartPolling(ctx context.Context) {
	tickTime := lo.
		IfF(t.tickTime != nil, func() time.Duration {
			return *t.tickTime
		}).
		Else(lo.
			If(t.services.Config.AppEnv == "development", 10*time.Second).
			Else(1 * time.Minute),
		)
	ticker := time.NewTicker(tickTime)

	go func() {
		for {
			select {
			case <-ticker.C:
				t.ticks++
				t.check(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
