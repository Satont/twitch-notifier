package commands

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
)

type StartCommand struct {
	*tg_types.CommandOpts
}

func (c *StartCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](3).Row(
		tg.NewInlineKeyboardButtonCallback("+", "+"),
		tg.NewInlineKeyboardButtonCallback("+", "+"),
		tg.NewInlineKeyboardButtonCallback("+", "+"),
		tg.NewInlineKeyboardButtonCallback("+", "+"),
		tg.NewInlineKeyboardButtonCallback("+", "+"),
	)
	keyBoard := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	spew.Dump(c.SessionManager.Get(ctx).Chat)

	description := c.Services.I18N.Translate("bot.description", "en", nil)

	return msg.Answer(description).ReplyMarkup(keyBoard).DoVoid(ctx)
}

func NewStartCommand(opts *tg_types.CommandOpts) {
	cmd := &StartCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.HandleCommand, tgb.Command("start",
		tgb.WithCommandAlias("help"),
		tgb.WithCommandAlias("info"),
	))
}
