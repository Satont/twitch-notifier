package commands

import (
	"context"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

var channelsAdminFilter = tgb.FilterFunc(func(ctx context.Context, update *tgb.Update) (bool, error) {
	if update.Chat().Type == tg.ChatTypePrivate || update.Chat().Type == tg.ChatTypeSender {
		return true, nil
	}

	admins, err := update.Client.GetChatAdministrators(update.Chat().ID).Do(ctx)
	if err != nil {
		return false, err
	}

	if update.CallbackQuery != nil {
		for _, admin := range admins {
			if admin.User.ID == update.CallbackQuery.From.ID {
				return true, nil
			}
		}
	} else if update.Message != nil && update.Message.From != nil {
		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {
				return true, nil
			}
		}
	} else {
		return true, nil
	}

	return false, nil
})
