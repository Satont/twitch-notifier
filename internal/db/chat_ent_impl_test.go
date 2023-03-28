package db

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestChatService_GetByID(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	chatService := NewChatEntRepository(entClient)

	_, err = chatService.Create(
		context.Background(),
		"123",
		db_models.ChatServiceTelegram,
	)
	assert.NoError(t, err)

	table := []struct {
		name    string
		chatID  string
		wantNil bool
		expects struct {
			chatID                  string
			service                 db_models.ChatService
			language                db_models.ChatLanguage
			gameChangeNotification  bool
			streamStartNotification bool
		}
	}{
		{
			name:    "Get chat by id",
			chatID:  "123",
			wantNil: false,
			expects: struct {
				chatID                  string
				service                 db_models.ChatService
				language                db_models.ChatLanguage
				gameChangeNotification  bool
				streamStartNotification bool
			}{
				chatID:                  "123",
				service:                 db_models.ChatServiceTelegram,
				language:                db_models.ChatLanguageEn,
				gameChangeNotification:  true,
				streamStartNotification: true,
			},
		},
		{
			name:    "Should fail if chat not found",
			chatID:  "321",
			wantNil: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.chatID, func(t *testing.T) {
			chat, err := chatService.GetByID(
				context.Background(),
				tt.chatID,
				db_models.ChatServiceTelegram,
			)

			if tt.wantNil {
				assert.Nil(t, chat)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expects.chatID, chat.ChatID)
				assert.Equal(t, tt.expects.service, chat.Service)
				assert.Equal(t, tt.expects.language, chat.Settings.ChatLanguage)
				assert.Equal(t, tt.expects.gameChangeNotification, chat.Settings.GameChangeNotification)
				assert.Equal(t, tt.expects.streamStartNotification, chat.Settings.OfflineNotification)
				assert.Equal(t, chat.ID, chat.Settings.ChatID)
			}
		})
	}
}

func TestChatService_Create(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	chatService := NewChatEntRepository(entClient)

	table := []struct {
		name    string
		chatID  string
		wantErr bool
	}{
		{
			name:    "Create chat",
			chatID:  "123",
			wantErr: false,
		},
		{
			name:    "Should fail if chat already exists",
			chatID:  "123",
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.chatID, func(t *testing.T) {
			chat, err := chatService.Create(
				context.Background(),
				tt.chatID,
				db_models.ChatServiceTelegram,
			)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.chatID, chat.ChatID)
				assert.Equal(t, db_models.ChatServiceTelegram, chat.Service)
				assert.NotEmpty(t, chat.Settings.ID)
				assert.Equal(t, db_models.ChatLanguageEn, chat.Settings.ChatLanguage)
				assert.Equal(t, true, chat.Settings.GameChangeNotification)
				assert.Equal(t, true, chat.Settings.OfflineNotification)
				assert.Equal(t, chat.ID, chat.Settings.ChatID)
			}
		})
	}
}

func TestChatService_Update(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	chatService := NewChatEntRepository(entClient)

	table := []struct {
		name         string
		chatID       string
		wantErr      bool
		shouldCreate bool
		newValues    struct {
			language                db_models.ChatLanguage
			gameChangeNotification  bool
			streamStartNotification bool
		}
	}{
		{
			name:         "Update chat",
			chatID:       "123",
			wantErr:      false,
			shouldCreate: true,
			newValues: struct {
				language                db_models.ChatLanguage
				gameChangeNotification  bool
				streamStartNotification bool
			}{
				language:                db_models.ChatLanguageRu,
				gameChangeNotification:  false,
				streamStartNotification: false,
			},
		},
		{
			name:         "Should fail if chat not found",
			chatID:       "321",
			wantErr:      true,
			shouldCreate: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.chatID, func(t *testing.T) {
			if tt.shouldCreate {
				_, err = chatService.Create(
					context.Background(),
					tt.chatID,
					db_models.ChatServiceTelegram,
				)
				assert.NoError(t, err)
			}

			newChat, err := chatService.Update(
				context.Background(),
				tt.chatID,
				db_models.ChatServiceTelegram,
				&ChatUpdateQuery{
					Settings: &ChatUpdateSettingsQuery{
						GameChangeNotification: lo.ToPtr(false),
						OfflineNotification:    lo.ToPtr(false),
						ChatLanguage:           lo.ToPtr(db_models.ChatLanguageRu),
					},
				},
			)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.chatID, newChat.ChatID)
				assert.Equal(t, tt.newValues.language, newChat.Settings.ChatLanguage)
				assert.Equal(t, tt.newValues.gameChangeNotification, newChat.Settings.GameChangeNotification)
				assert.Equal(t, tt.newValues.streamStartNotification, newChat.Settings.OfflineNotification)
			}
		})
	}
}

func TestChatService_GetAllByService(t *testing.T) {
	t.Parallel()

	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	ctx := context.Background()

	chatService := NewChatEntRepository(entClient)

	var created []*db_models.Chat

	for i := 0; i < 10; i++ {
		newChat, err := chatService.Create(
			ctx,
			strconv.Itoa(i),
			db_models.ChatServiceTelegram,
		)
		assert.NoError(t, err)
		created = append(created, newChat)
	}

	chats, err := chatService.GetAllByService(ctx, db_models.ChatServiceTelegram)
	assert.NoError(t, err)
	assert.Len(t, chats, 10)

	for _, chat := range chats {
		assert.Contains(t, created, chat)
	}
}
