package db_models

import (
	"github.com/google/uuid"
)

type ChatService string

const (
	ChatServiceTelegram ChatService = "telegram"
)

func (s ChatService) String() string {
	return string(s)
}

func LanguageExists(l ChatLanguage) bool {
	switch l {
	case ChatLanguageRu, ChatLanguageEn:
		return true
	default:
		return false
	}
}

type Chat struct {
	ID      uuid.UUID   `json:"id,omitempty"`
	ChatID  string      `json:"chat_id,omitempty"`
	Service ChatService `json:"service,omitempty"`

	Follows  []*Follow     `json:"follows,omitempty"`
	Settings *ChatSettings `json:"settings,omitempty"`
}

type ChatLanguage string

var DefaultChatLanguage = ChatLanguageEn

var (
	ChatLanguageRu ChatLanguage = "ru"
	ChatLanguageEn ChatLanguage = "en"
)

func (cl ChatLanguage) String() string {
	return string(cl)
}

type ChatSettings struct {
	ID                      uuid.UUID    `json:"id,omitempty"`
	GameChangeNotification  bool         `json:"game_change_notification,omitempty"`
	TitleChangeNotification bool         `json:"title_change_notification,omitempty"`
	OfflineNotification     bool         `json:"offline_notification,omitempty"`
	ChatLanguage            ChatLanguage `json:"chat_language,omitempty"`
	ChatID                  uuid.UUID    `json:"chat_id,omitempty"`
}
