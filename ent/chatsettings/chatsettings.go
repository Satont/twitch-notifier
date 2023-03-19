// Code generated by ent, DO NOT EDIT.

package chatsettings

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the chatsettings type in the database.
	Label = "chat_settings"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldGameChangeNotification holds the string denoting the game_change_notification field in the database.
	FieldGameChangeNotification = "game_change_notification"
	// FieldOfflineNotification holds the string denoting the offline_notification field in the database.
	FieldOfflineNotification = "offline_notification"
	// FieldChatLanguage holds the string denoting the chat_language field in the database.
	FieldChatLanguage = "chat_language"
	// FieldChatID holds the string denoting the chat_id field in the database.
	FieldChatID = "chat_id"
	// EdgeChat holds the string denoting the chat edge name in mutations.
	EdgeChat = "chat"
	// Table holds the table name of the chatsettings in the database.
	Table = "chat_settings"
	// ChatTable is the table that holds the chat relation/edge.
	ChatTable = "chat_settings"
	// ChatInverseTable is the table name for the Chat entity.
	// It exists in this package in order to avoid circular dependency with the "chat" package.
	ChatInverseTable = "chats"
	// ChatColumn is the table column denoting the chat relation/edge.
	ChatColumn = "chat_id"
)

// Columns holds all SQL columns for chatsettings fields.
var Columns = []string{
	FieldID,
	FieldGameChangeNotification,
	FieldOfflineNotification,
	FieldChatLanguage,
	FieldChatID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultGameChangeNotification holds the default value on creation for the "game_change_notification" field.
	DefaultGameChangeNotification bool
	// DefaultOfflineNotification holds the default value on creation for the "offline_notification" field.
	DefaultOfflineNotification bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// ChatLanguage defines the type for the "chat_language" enum field.
type ChatLanguage string

// ChatLanguageEn is the default value of the ChatLanguage enum.
const DefaultChatLanguage = ChatLanguageEn

// ChatLanguage values.
const (
	ChatLanguageRu ChatLanguage = "ru"
	ChatLanguageEn ChatLanguage = "en"
)

func (cl ChatLanguage) String() string {
	return string(cl)
}

// ChatLanguageValidator is a validator for the "chat_language" field enum values. It is called by the builders before save.
func ChatLanguageValidator(cl ChatLanguage) error {
	switch cl {
	case ChatLanguageRu, ChatLanguageEn:
		return nil
	default:
		return fmt.Errorf("chatsettings: invalid enum value for chat_language field: %q", cl)
	}
}