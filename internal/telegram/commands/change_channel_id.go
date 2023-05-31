package commands

import (
	"context"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/telegram/types"
	"github.com/satont/twitch-notifier/internal/types"
	"go.uber.org/zap"
	"strings"
)

type ChangeChannelId struct {
	*tgtypes.CommandOpts
}

func (c *ChangeChannelId) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	text := strings.ReplaceAll(msg.Message.Text, "/change_channel_id ", "")
	splittedText := strings.Split(strings.TrimSpace(text), " ")

	if len(splittedText) != 2 {
		return nil
	}

	sourceChannelID := splittedText[0]
	targetChannelID := splittedText[1]

	_, err := c.Services.Channel.Update(
		ctx,
		sourceChannelID,
		db_models.ChannelServiceTwitch,
		&db.ChannelUpdateQuery{
			DangerNewChannelId: &targetChannelID,
		},
	)

	if err != nil {
		zap.S().Error(err)
	}

	return msg.Answer("done").DoVoid(ctx)
}

var (
	changeChannelIdFilter            = tgb.Command("change_channel_id")
	changeChannelIdFilterAdminFilter = func(services *types.Services) tgb.Filter {
		return tgb.FilterFunc(func(ctx context.Context, update *tgb.Update) (bool, error) {
			return lo.Contains(services.Config.TelegramBotAdmins, update.Message.Chat.ID.PeerID()), nil
		})
	}
)

func NewChangeChannelId(opts *tgtypes.CommandOpts) {
	cmd := &ChangeChannelId{
		CommandOpts: opts,
	}
	opts.Router.Message(
		cmd.HandleCommand,
		changeChannelIdFilter,
		changeChannelIdFilterAdminFilter(opts.Services),
	)
}
