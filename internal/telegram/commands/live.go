package commands

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/telegram/types"
	"go.uber.org/zap"
	"strings"
	"time"
)

type LiveCommand struct {
	*tgtypes.CommandOpts
}

type liveChannel struct {
	Name      string
	Login     string
	StartedAt time.Time
	Title     string
	Category  string
	Viewers   int
}

func (c *LiveCommand) getList(ctx context.Context) ([]*liveChannel, error) {
	chat := c.SessionManager.Get(ctx).Chat

	follows, err := c.Services.Follow.GetByChatID(ctx, chat.ID, 0, 0)
	if err != nil {
		return nil, err
	}

	if len(follows) == 0 {
		return nil, nil
	}

	channelsIds := lo.Map(follows, func(follow *db_models.Follow, _ int) string {
		return follow.Channel.ChannelID
	})

	streams, err := c.Services.Twitch.GetStreamsByUserIds(channelsIds)
	if err != nil {
		return nil, err
	}

	if len(streams) == 0 {
		return nil, nil
	}

	result := make([]*liveChannel, 0, len(streams))

	for _, stream := range streams {
		result = append(result, &liveChannel{
			Name:      stream.UserName,
			Login:     stream.UserLogin,
			StartedAt: stream.StartedAt,
			Title:     stream.Title,
			Category:  stream.GameName,
			Viewers:   stream.ViewerCount,
		})
	}

	return result, nil
}

func (c *LiveCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	list, err := c.getList(ctx)
	if err != nil {
		zap.S().Error(err)
		return msg.Answer("internal error").DoVoid(ctx)
	}

	if len(list) == 0 {
		return msg.Answer("No one online").DoVoid(ctx)
	}

	message := make([]string, 0, len(list))

	for _, channel := range list {
		channelMessage := make([]string, 0)

		channelMessage = append(
			channelMessage,
			fmt.Sprintf(
				"🟢 %s - %v 👁️️",
				tg.MD.Link(
					channel.Name,
					fmt.Sprintf("https://twitch.tv/%s", channel.Login),
				),
				channel.Viewers,
			),
		)

		if channel.Category != "" {
			channelMessage = append(channelMessage, fmt.Sprintf("🎮 %s", channel.Category))
		}

		if channel.Title != "" {
			channelMessage = append(channelMessage, fmt.Sprintf("📝 %s", channel.Title))
		}

		since := time.Since(channel.StartedAt)
		hour := int(since.Seconds() / 3600)
		minute := int(since.Seconds()/60) % 60
		second := int(since.Seconds()) % 60

		uptime := "⌛ "
		if hour > 0 {
			uptime += fmt.Sprintf("%vh ", hour)
		}

		if minute > 0 {
			uptime += fmt.Sprintf("%vm ", minute)
		}

		if second > 0 {
			uptime += fmt.Sprintf("%vs ", second)
		}

		channelMessage = append(channelMessage, uptime)

		message = append(
			message,
			strings.Join(channelMessage, "\n"),
		)
	}

	return msg.
		Answer(strings.Join(message, "\n\n")).
		ParseMode(tg.MD).
		DisableWebPagePreview(true).
		DoVoid(ctx)
}

var liveCommandFilter = tgb.Command("live")

func NewLiveCommand(opts *tgtypes.CommandOpts) {
	cmd := &LiveCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, channelsAdminFilter, liveCommandFilter)
}
