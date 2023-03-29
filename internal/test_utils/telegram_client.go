package test_utils

import (
	"github.com/mr-linch/go-tg"
	"net/http"
	"net/http/httptest"
)

const (
	TelegramClientToken = "1234:secret"
	TelegramOkResponse  = `{"ok":true}`
)

func NewTelegramClient(server *httptest.Server) *tg.Client {
	client := tg.New(TelegramClientToken,
		tg.WithClientServerURL(server.URL),
		tg.WithClientDoer(&http.Client{}),
	)

	return client
}
