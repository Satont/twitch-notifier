package twitch_streams_cheker

import (
	"context"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/message_sender"
	"github.com/satont/twitch-notifier/internal/services/types"
	"go.uber.org/zap"
	"time"
)

type TwitchStreamChecker struct {
	services *types.Services
	ticks    int
	tickTime *time.Duration
	sender   message_sender.MessageSenderInterface
}

func NewTwitchStreamChecker(
	services *types.Services,
	sender message_sender.MessageSenderInterface,
	tickTime *time.Duration,
) *TwitchStreamChecker {
	checker := &TwitchStreamChecker{
		services: services,
		tickTime: tickTime,
		sender:   sender,
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

	currentTwitchStreams, err := t.services.Twitch.GetStreamsByUserIds(channelsIDs)
	if err != nil {
		zap.S().Error(err)
		return
	}

	for _, channel := range channels {
		currentDBStream, err := t.services.Stream.GetLatestByChannelID(ctx, channel.ID)
		if err != nil {
			zap.S().Error(err)
			continue
		}

		twitchCurrentStream, twitchCurrentStreamOk := lo.Find(currentTwitchStreams, func(stream helix.Stream) bool {
			return stream.UserID == channel.ChannelID
		})

		// if stream becomes offline
		if !twitchCurrentStreamOk && currentDBStream != nil {
			_, err = t.services.Stream.UpdateOneByStreamID(ctx, currentDBStream.ID, &db.StreamUpdateQuery{
				IsLive: lo.ToPtr(false),
			})
			if err != nil {
				zap.S().Error(err)
				continue
			}
		}

		// if stream becomes online
		if twitchCurrentStreamOk && currentDBStream == nil {
			_, err = t.services.Stream.CreateOneByChannelID(ctx, channel.ID, &db.StreamUpdateQuery{
				StreamID: twitchCurrentStream.ID,
				IsLive:   lo.ToPtr(true),
				Category: lo.ToPtr(twitchCurrentStream.GameName),
				Title:    lo.ToPtr(twitchCurrentStream.Title),
			})
			if err != nil {
				zap.S().Error(err)
				continue
			}
		}

		// stream is still online
		if twitchCurrentStreamOk && currentDBStream != nil {
			latestTitle := currentDBStream.Titles[len(currentDBStream.Titles)-1]
			latestCategory := currentDBStream.Categories[len(currentDBStream.Categories)-1]

			if twitchCurrentStream.GameName != latestCategory {
				_, err = t.services.Stream.UpdateOneByStreamID(ctx, currentDBStream.ID, &db.StreamUpdateQuery{
					Category: lo.ToPtr(twitchCurrentStream.GameName),
				})
				if err != nil {
					zap.S().Error(err)
					continue
				}
			}

			if twitchCurrentStream.Title != latestTitle {
				_, err = t.services.Stream.UpdateOneByStreamID(ctx, currentDBStream.ID, &db.StreamUpdateQuery{
					Title: lo.ToPtr(twitchCurrentStream.Title),
				})
				if err != nil {
					zap.S().Error(err)
					continue
				}
			}
		}
	}

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
