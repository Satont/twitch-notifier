package twitch_streams_cheker

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mr-linch/go-tg"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"github.com/satont/twitch-notifier/internal/message_sender"
	"github.com/satont/twitch-notifier/internal/types"
	"go.uber.org/zap"
)

type TwitchStreamChecker struct {
	services         *types.Services
	ticks            int
	tickTime         *time.Duration
	sender           message_sender.MessageSenderInterface
	thumbNailBuilder *thumbNailBuilder
}

func NewTwitchStreamChecker(
	services *types.Services,
	sender message_sender.MessageSenderInterface,
	tickTime *time.Duration,
) *TwitchStreamChecker {
	checker := &TwitchStreamChecker{
		services:         services,
		tickTime:         tickTime,
		sender:           sender,
		thumbNailBuilder: newThumbNailBuilder(),
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

	twitchChannels, err := t.services.Twitch.GetChannelsByUserIds(channelsIDs)
	if err != nil {
		zap.S().Error(err)
		return
	}

	currentTwitchStreams, err := t.services.Twitch.GetStreamsByUserIds(channelsIDs)
	if err != nil {
		zap.S().Error(err)
		return
	}

	wg := &sync.WaitGroup{}
	for _, channel := range channels {
		wg.Add(1)

		go func(channel *db_models.Channel) {
			defer wg.Done()
			twitchChannel, twitchChannelOk := lo.Find(
				twitchChannels,
				func(item helix.ChannelInformation) bool {
					return item.BroadcasterID == channel.ChannelID
				},
			)
			if !twitchChannelOk {
				return
			}

			currentDBStream, err := t.services.Stream.GetLatestByChannelID(ctx, channel.ID)
			if err != nil {
				zap.S().Error(err)
				return
			}

			followers, err := t.services.Follow.GetByChannelID(ctx, channel.ID)
			if err != nil {
				zap.S().Error(err)
				return
			}

			twitchCurrentStream, twitchCurrentStreamOk := lo.Find(
				currentTwitchStreams,
				func(stream helix.Stream) bool {
					return stream.UserID == channel.ChannelID
				},
			)

			if twitchCurrentStreamOk && twitchCurrentStream.Type != "live" {
				return
			}

			// if stream becomes offline
			if !twitchCurrentStreamOk && currentDBStream != nil && currentDBStream.EndedAt == nil {
				_, err = t.services.Stream.UpdateOneByStreamID(
					ctx,
					currentDBStream.ID,
					&db.StreamUpdateQuery{
						IsLive: lo.ToPtr(false),
					},
				)
				if err != nil {
					zap.S().Error(err)
					return
				}

				// send message to all followers
				for _, follower := range followers {
					if !follower.Chat.Settings.OfflineNotification {
						continue
					}

					message := t.services.I18N.Translate(
						"notifications.streams.nowOffline",
						follower.Chat.Settings.ChatLanguage.String(),
						map[string]string{
							"channelLink": tg.MD.Link(
								twitchChannel.BroadcasterName,
								fmt.Sprintf("https://twitch.tv/%s", twitchChannel.BroadcasterName),
							),
							"categories": strings.Join(currentDBStream.Categories, " -> "),
							"duration": time.Now().UTC().Sub(currentDBStream.StartedAt).
								Truncate(1 * time.Second).
								String(),
						},
					)
					unfollowButton := message_sender.KeyboardButton{
						Text: t.services.I18N.Translate(
							"commands.unfollow.callbackButton",
							follower.Chat.Settings.ChatLanguage.String(),
							map[string]string{
								"streamer": twitchChannel.BroadcasterName,
							},
						),
						CallbackData: fmt.Sprintf("channels_unfollow_%v", channel.ID),
						SkipInGroup:  true,
					}

					err = t.sender.SendMessage(
						ctx, follower.Chat, &message_sender.MessageOpts{
							Text:      message,
							ParseMode: &tg.MD,
							Buttons:   [][]message_sender.KeyboardButton{{unfollowButton}},
						},
					)
					if err != nil {
						zap.S().Error(err)
						continue
					}
				}
			}

			// if stream becomes online
			if twitchCurrentStreamOk && currentDBStream == nil {
				_, err = t.services.Stream.CreateOneByChannelID(
					ctx,
					channel.ID,
					&db.StreamUpdateQuery{
						StreamID: twitchCurrentStream.ID,
						IsLive:   lo.ToPtr(true),
						Category: lo.ToPtr(twitchCurrentStream.GameName),
						Title:    lo.ToPtr(twitchCurrentStream.Title),
					},
				)
				if err != nil {
					zap.S().Error(err)
					return
				}

				for _, follower := range followers {
					message := t.services.I18N.Translate(
						"notifications.streams.nowOnline",
						follower.Chat.Settings.ChatLanguage.String(),
						map[string]string{
							"channelLink": tg.MD.Link(
								twitchChannel.BroadcasterName,
								fmt.Sprintf("https://twitch.tv/%s", twitchChannel.BroadcasterName),
							),
							"category": twitchCurrentStream.GameName,
							"title":    twitchCurrentStream.Title,
						},
					)

					unfollowButton := message_sender.KeyboardButton{
						Text: t.services.I18N.Translate(
							"commands.unfollow.callbackButton",
							follower.Chat.Settings.ChatLanguage.String(),
							map[string]string{
								"streamer": twitchChannel.BroadcasterName,
							},
						),
						CallbackData: fmt.Sprintf("channels_unfollow_%v", channel.ID),
						SkipInGroup:  true,
					}

					thumbNail, err := t.thumbNailBuilder.Build(twitchCurrentStream.ThumbnailURL, true)
					if err != nil {
						zap.S().Error(err)
					}

					err = t.sender.SendMessage(
						ctx, follower.Chat, &message_sender.MessageOpts{
							Text: message,
							ImageURL: lo.If(
								follower.Chat.Settings.ImageInNotification,
								fmt.Sprintf("%s?%d", thumbNail, time.Now().Unix()),
							).Else(""),
							ParseMode: &tg.MD,
							Buttons:   [][]message_sender.KeyboardButton{{unfollowButton}},
						},
					)
					if err != nil {
						zap.S().Error(err)
						continue
					}
				}
			}

			// stream is still online, need to check do we need to update title or category
			if twitchCurrentStreamOk && currentDBStream != nil &&
				currentDBStream.ID == twitchCurrentStream.ID {
				latestTitle := ""
				if len(currentDBStream.Titles) > 0 {
					latestTitle = currentDBStream.Titles[len(currentDBStream.Titles)-1]
				}
				latestCategory := ""
				if len(currentDBStream.Categories) > 0 {
					latestCategory = currentDBStream.Categories[len(currentDBStream.Categories)-1]
				}

				// stream is online, and both title and category changed, so we need to send a complex notification
				if twitchCurrentStream.GameName != latestCategory &&
					twitchCurrentStream.Title != latestTitle {
					_, err = t.services.Stream.UpdateOneByStreamID(
						ctx,
						currentDBStream.ID,
						&db.StreamUpdateQuery{
							Category: lo.ToPtr(twitchCurrentStream.GameName),
							Title:    lo.ToPtr(twitchCurrentStream.Title),
						},
					)
					if err != nil {
						zap.S().Error(err)
						return
					}

					thumbNail, err := t.thumbNailBuilder.Build(twitchCurrentStream.ThumbnailURL, false)
					if err != nil {
						zap.S().Error(err)
					}

					for _, follower := range followers {
						if !follower.Chat.Settings.GameAndTitleChangeNotification {
							continue
						}

						unfollowButton := message_sender.KeyboardButton{
							Text: t.services.I18N.Translate(
								"commands.unfollow.callbackButton",
								follower.Chat.Settings.ChatLanguage.String(),
								map[string]string{
									"streamer": twitchChannel.BroadcasterName,
								},
							),
							CallbackData: fmt.Sprintf("channels_unfollow_%v", channel.ID),
							SkipInGroup:  true,
						}

						err = t.sender.SendMessage(
							ctx, follower.Chat, &message_sender.MessageOpts{
								Text: t.services.I18N.Translate(
									"notifications.streams.titleAndCategoryChanged",
									follower.Chat.Settings.ChatLanguage.String(),
									map[string]string{
										"channelLink": tg.MD.Link(
											twitchChannel.BroadcasterName,
											fmt.Sprintf(
												"https://twitch.tv/%s",
												twitchChannel.BroadcasterName,
											),
										),
										"category":    tg.MD.Bold(twitchCurrentStream.GameName),
										"oldCategory": tg.MD.Bold(latestCategory),
										"title":       tg.MD.Bold(twitchCurrentStream.Title),
										"oldTitle":    tg.MD.Bold(latestTitle),
									},
								),
								ParseMode: &tg.MD,
								ImageURL: lo.If(
									follower.Chat.Settings.ImageInNotification,
									fmt.Sprintf("%s?%d", thumbNail, time.Now().Unix()),
								).Else(""),
								Buttons: [][]message_sender.KeyboardButton{{unfollowButton}},
							},
						)
						if err != nil {
							zap.S().Error(err)
							continue
						}
					}
					return
				}

				if twitchCurrentStream.GameName != latestCategory {
					_, err = t.services.Stream.UpdateOneByStreamID(
						ctx,
						currentDBStream.ID,
						&db.StreamUpdateQuery{
							Category: lo.ToPtr(twitchCurrentStream.GameName),
						},
					)
					if err != nil {
						zap.S().Error(err)
						return
					}

					for _, follower := range followers {
						if !follower.Chat.Settings.GameChangeNotification {
							continue
						}

						thumbNail, err := t.thumbNailBuilder.Build(twitchCurrentStream.ThumbnailURL, true)
						if err != nil {
							zap.S().Error(err)
						}

						err = t.sender.SendMessage(
							ctx, follower.Chat, &message_sender.MessageOpts{
								Text: t.services.I18N.Translate(
									"notifications.streams.newCategory",
									follower.Chat.Settings.ChatLanguage.String(),
									map[string]string{
										"channelLink": tg.MD.Link(
											twitchChannel.BroadcasterName,
											fmt.Sprintf(
												"https://twitch.tv/%s",
												twitchChannel.BroadcasterName,
											),
										),
										"category":    tg.MD.Bold(twitchCurrentStream.GameName),
										"oldCategory": tg.MD.Bold(latestCategory),
									},
								),
								ParseMode: &tg.MD,
								ImageURL: lo.If(
									follower.Chat.Settings.ImageInNotification,
									fmt.Sprintf("%s?%d", thumbNail, time.Now().Unix()),
								).Else(""),
							},
						)
						if err != nil {
							zap.S().Error(err)
							continue
						}
					}
					return
				}

				if twitchCurrentStream.Title != latestTitle {
					_, err = t.services.Stream.UpdateOneByStreamID(
						ctx,
						currentDBStream.ID,
						&db.StreamUpdateQuery{
							Title: lo.ToPtr(twitchCurrentStream.Title),
						},
					)
					if err != nil {
						zap.S().Error(err)
						return
					}

					for _, follower := range followers {
						if !follower.Chat.Settings.TitleChangeNotification {
							continue
						}

						thumbNail, err := t.thumbNailBuilder.Build(twitchCurrentStream.ThumbnailURL, true)
						if err != nil {
							zap.S().Error(err)
						}

						err = t.sender.SendMessage(
							ctx, follower.Chat, &message_sender.MessageOpts{
								Text: t.services.I18N.Translate(
									"notifications.streams.titleChanged",
									follower.Chat.Settings.ChatLanguage.String(),
									map[string]string{
										"channelLink": tg.MD.Link(
											twitchChannel.BroadcasterName,
											fmt.Sprintf(
												"https://twitch.tv/%s",
												twitchChannel.BroadcasterName,
											),
										),
										"category": twitchCurrentStream.GameName,
										"title":    tg.MD.Bold(twitchCurrentStream.Title),
										"oldTitle": tg.MD.Bold(latestTitle),
									},
								),
								ParseMode: &tg.MD,
								ImageURL: lo.If(
									follower.Chat.Settings.ImageInNotification,
									fmt.Sprintf("%s?%d", thumbNail, time.Now().Unix()),
								).Else(""),
							},
						)
						if err != nil {
							zap.S().Error(err)
							continue
						}
					}
				}

			}
		}(channel)
	}
	wg.Wait()
}

func (t *TwitchStreamChecker) StartPolling(ctx context.Context) {
	tickTime := lo.
		IfF(
			t.tickTime != nil, func() time.Duration {
				return *t.tickTime
			},
		).
		Else(
			lo.
				If(t.services.Config.AppEnv == "development", 10*time.Second).
				Else(1 * time.Minute),
		)
	ticker := time.NewTicker(tickTime)

	t.check(ctx)

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
