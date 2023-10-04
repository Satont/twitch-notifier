package message_sender

import (
	"context"

	"github.com/mr-linch/go-tg"
	"github.com/satont/twitch-notifier/internal/db/db_models"
)

type MessageOpts struct {
	Text        string
	ImageURL    string
	ParseMode   *tg.ParseMode
	Buttons     [][]KeyboardButton
	SkipButtons bool
}

type KeyboardButton struct {
	// kostil chto bi skipnut knopki v gruppah
	SkipInGroup bool

	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	// this is not needed currently
	// URL                          string        `json:"url,omitempty"`
	// WebApp                       *WebAppInfo   `json:"web_app,omitempty"`
	// LoginURL                     *LoginURL     `json:"login_url,omitempty"`
	// SwitchInlineQuery            string        `json:"switch_inline_query,omitempty"`
	// SwitchInlineQueryCurrentChat string        `json:"switch_inline_query_current_chat,omitempty"`
	// CallbackGame                 *CallbackGame `json:"callback_game,omitempty"`
	// Pay                          bool          `json:"pay,omitempty"`
}

type MessageSenderInterface interface {
	SendMessage(ctx context.Context, chat *db_models.Chat, opts *MessageOpts) error
}
