package commands

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/telegram/types"
	"github.com/satont/twitch-notifier/internal/types"
	"go.uber.org/zap"
)

type BroadcastCommand struct {
	*tgtypes.CommandOpts
}

func (c *BroadcastCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	chats, err := c.Services.Chat.GetAllByService(ctx, db_models.ChatServiceTelegram)
	if err != nil {
		zap.S().Error(err)
		return msg.Answer("Error").DoVoid(ctx)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(chats))

	for _, chat := range chats {
		go func(chat *db_models.Chat) {
			defer wg.Done()

			chatId, _ := strconv.Atoi(chat.ChatID)

			// filter channels, thay have negative id.
			if chatId <= 0 {
				return
			}

			err := msg.Client.
				SendMessage(
					tg.ChatID(chatId),
					strings.Replace(msg.Message.Text, "/broadcast ", "", 1),
				).DoVoid(ctx)
			if err != nil {
				zap.S().Error(err)
			}
		}(chat)
	}

	wg.Wait()

	return nil
}

var (
	broadcastCommandFilter      = tgb.Command("broadcast")
	broadcastCommandAdminFilter = func(services *types.Services) tgb.Filter {
		return tgb.FilterFunc(func(ctx context.Context, update *tgb.Update) (bool, error) {
			return lo.Contains(services.Config.TelegramBotAdmins, update.Message.Chat.ID.PeerID()), nil
		})
	}
)

func NewBroadcastCommand(opts *tgtypes.CommandOpts) {
	cmd := &BroadcastCommand{
		CommandOpts: opts,
	}
	opts.Router.Message(
		cmd.HandleCommand,
		broadcastCommandFilter,
		broadcastCommandAdminFilter(opts.Services),
	)
}
