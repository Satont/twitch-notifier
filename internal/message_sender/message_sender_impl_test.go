package message_sender

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/satont/twitch-notifier/internal/db/db_models"

	"github.com/satont/twitch-notifier/internal/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestMessageSender_SendMessage(t *testing.T) {
	t.Parallel()

	chat := &db_models.Chat{
		ChatID:  "-123",
		Service: db_models.ChatServiceTelegram,
	}

	table := []struct {
		name         string
		chat         *db_models.Chat
		opts         *MessageOpts
		createServer func(*testing.T) *httptest.Server
	}{
		{
			name: "should call send message method",
			chat: chat,
			opts: &MessageOpts{
				Text: "test",
			},
			createServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							body, err := io.ReadAll(r.Body)
							assert.NoError(t, err)
							query, err := url.ParseQuery(string(body))
							assert.NoError(t, err)

							assert.Equal(t, "test", query.Get("text"))
							assert.Equal(t, "-123", query.Get("chat_id"))

							assert.Equal(t, http.MethodPost, r.Method)
							assert.Equal(
								t,
								fmt.Sprintf("/bot%s/sendMessage", test_utils.TelegramClientToken),
								r.URL.Path,
							)

							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
						},
					),
				)
			},
		},
		{
			name: "should call send photo method",
			chat: chat,
			opts: &MessageOpts{
				Text:     "test photo",
				ImageURL: "https://example.com/image.jpg",
			},
			createServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							body, err := io.ReadAll(r.Body)
							assert.NoError(t, err)
							query, err := url.ParseQuery(string(body))
							assert.NoError(t, err)

							assert.Equal(t, "test photo", query.Get("caption"))
							assert.Equal(t, "https://example.com/image.jpg", query.Get("photo"))
							assert.Equal(t, "-123", query.Get("chat_id"))

							assert.Equal(t, http.MethodPost, r.Method)
							assert.Equal(
								t,
								fmt.Sprintf("/bot%s/sendPhoto", test_utils.TelegramClientToken),
								r.URL.Path,
							)

							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
						},
					),
				)
			},
		},
		{
			name: "should call send message method with parse mode",
			chat: chat,
			opts: &MessageOpts{
				Text:        "test md",
				TgParseMode: TgParseModeMD,
			},
			createServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							body, err := io.ReadAll(r.Body)
							assert.NoError(t, err)
							query, err := url.ParseQuery(string(body))
							assert.NoError(t, err)

							assert.Equal(t, "test md", query.Get("text"))
							assert.Equal(t, "-123", query.Get("chat_id"))
							assert.Equal(t, "Markdown", query.Get("parse_mode"))

							assert.Equal(t, http.MethodPost, r.Method)
							assert.Equal(
								t,
								fmt.Sprintf("/bot%s/sendMessage", test_utils.TelegramClientToken),
								r.URL.Path,
							)

							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
						},
					),
				)
			},
		},
		{
			name: "should send keyboard buttons",
			chat: chat,
			opts: &MessageOpts{
				Text: "test buttons",
				Buttons: [][]KeyboardButton{
					{
						KeyboardButton{Text: "click me", CallbackData: "click"},
					},
				},
			},
			createServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							body, err := io.ReadAll(r.Body)
							assert.NoError(t, err)
							query, err := url.ParseQuery(string(body))
							assert.NoError(t, err)

							assert.Equal(t, "test buttons", query.Get("text"))
							assert.Equal(t, "-123", query.Get("chat_id"))

							keyboard := map[string]any{}

							err = json.Unmarshal([]byte(query.Get("reply_markup")), &keyboard)
							assert.NoError(t, err)

							assert.Equal(
								t,
								"click me",
								keyboard["inline_keyboard"].([]interface{})[0].([]interface{})[0].(map[string]any)["text"],
							)
							assert.Equal(
								t,
								"click",
								keyboard["inline_keyboard"].([]interface{})[0].([]interface{})[0].(map[string]any)["callback_data"],
							)

							assert.Equal(t, http.MethodPost, r.Method)
							assert.Equal(
								t,
								fmt.Sprintf("/bot%s/sendMessage", test_utils.TelegramClientToken),
								r.URL.Path,
							)

							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
						},
					),
				)
			},
		},
		{
			name: "should skip button",
			chat: chat,
			opts: &MessageOpts{
				Text: "test buttons",
				Buttons: [][]KeyboardButton{
					{
						KeyboardButton{Text: "click me", CallbackData: "click", SkipInGroup: true},
					},
				},
			},
			createServer: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							body, err := io.ReadAll(r.Body)
							assert.NoError(t, err)
							query, err := url.ParseQuery(string(body))
							assert.NoError(t, err)

							assert.Equal(t, "test buttons", query.Get("text"))
							assert.Equal(t, "-123", query.Get("chat_id"))
							assert.Empty(t, query.Get("reply_markup"))

							assert.Equal(t, http.MethodPost, r.Method)
							assert.Equal(
								t,
								fmt.Sprintf("/bot%s/sendMessage", test_utils.TelegramClientToken),
								r.URL.Path,
							)

							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
						},
					),
				)
			},
		},
	}

	for _, tt := range table {
		t.Run(
			tt.name, func(c *testing.T) {
				server := tt.createServer(c)
				tgClient := test_utils.NewTelegramClient(server)
				sender := NewMessageSender(tgClient)

				err := sender.SendMessage(context.Background(), tt.opts)
				assert.NoError(c, err)
				assert.Nil(c, err)
			},
		)
	}
}
