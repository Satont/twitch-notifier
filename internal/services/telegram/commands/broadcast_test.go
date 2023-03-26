package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
	tg_types "github.com/satont/twitch-notifier/internal/services/telegram/types"
	"github.com/satont/twitch-notifier/internal/services/types"
	"github.com/satont/twitch-notifier/internal/test_utils"
	"github.com/satont/twitch-notifier/internal/test_utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBroadcastCommand_HandleCommand(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	chatMock := &mocks.DbChatMock{}

	table := []struct {
		name       string
		message    *tgb.MessageUpdate
		serverMock *httptest.Server
		setupMocks func()
	}{
		{
			name: "Should call SendMessage for each chat",
			message: &tgb.MessageUpdate{
				Message: &tg.Message{
					Text: "/broadcast test",
				},
			},
			serverMock: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				query, err := url.ParseQuery(string(body))
				assert.NoError(t, err)
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(
					t,
					fmt.Sprintf("/bot%s/sendMessage", test_utils.TelegramClientToken),
					r.URL.Path,
				)
				assert.Equal(t, "test", query.Get("text"))
				assert.Contains(t, []string{"1", "2"}, query.Get("chat_id"))

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(test_utils.TelegramOkResponse))
			})),
			setupMocks: func() {
				chatMock.
					On("GetAllByService", ctx, db_models.ChatServiceTelegram).
					Return(
						[]*db_models.Chat{{ChatID: "1"}, {ChatID: "2"}},
						nil,
					)
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer tt.serverMock.Close()
			tt.setupMocks()
			client := test_utils.NewTelegramClient(tt.serverMock)
			tt.message.Client = client
			cmd := &BroadcastCommand{
				CommandOpts: &tg_types.CommandOpts{
					Services: &types.Services{
						Chat: chatMock,
					},
				},
			}
			err := cmd.HandleCommand(ctx, tt.message)
			assert.NoError(t, err)

			chatMock.AssertExpectations(t)
		})
	}
}
