package message_sender

import (
	"context"
)

type TgParseMode string

const (
	TgParseModeMD TgParseMode = "markdown"
)

type MessageOpts struct {
	ChatID      string
	ChatService string
	Text        string
	ImageURL    string

	TgParseMode TgParseMode
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
	SendMessage(ctx context.Context, opts *MessageOpts) error
}
