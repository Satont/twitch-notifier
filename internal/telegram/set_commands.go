package telegram

import (
	"context"
	"github.com/mr-linch/go-tg"
	"go.uber.org/zap"
	"strconv"
)

var defaultCommands = []tg.BotCommand{
	{
		Command:     "follow",
		Description: "Follow to notifications of some streamer",
	},
	{
		Command:     "follows",
		Description: "Show list of followed streamers",
	},
	{
		Command:     "live",
		Description: "Show list of live streamers",
	},
	{
		Command:     "start",
		Description: "Bot settings",
	},
}

func (c *TelegramService) setMyCommands(ctx context.Context) {
	err := c.Client.
		SetMyCommands(defaultCommands).
		Scope(tg.BotCommandScopeDefault{}).
		DoVoid(ctx)
	if err != nil {
		zap.S().Fatalln("Can't set default commands", err)
	}

	for _, admin := range c.services.Config.TelegramBotAdmins {
		newCommands := append(defaultCommands, tg.BotCommand{
			Command:     "broadcast",
			Description: "Send message to all users",
		})

		chatID, err := strconv.Atoi(admin)
		if err != nil {
			zap.S().Errorw("Can't parse chat id", "chatID", admin)
			return
		}

		err = c.Client.
			SetMyCommands(newCommands).
			Scope(tg.BotCommandScopeChat{ChatID: tg.ChatID(chatID)}).
			DoVoid(ctx)
		if err != nil {
			zap.S().Fatalln("Can't set admin commands", err)
		}
	}
}
