package commands

import (
	"context"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"go.uber.org/zap"
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

	return nil
}

func NewLiveCommand(opts *tgtypes.CommandOpts) {
	cmd := &LiveCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, tgb.Command("follow"))
}
