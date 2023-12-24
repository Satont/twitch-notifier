package domain

import (
	"github.com/google/uuid"
)

type ChatSettings struct {
	ID                            uuid.UUID
	ChatID                        uuid.UUID
	Language                      Language
	CategoryChangeNotifications   bool
	TitleChangeNotifications      bool
	OfflineNotifications          bool
	CategoryAndTitleNotifications bool
	ShowThumbnail                 bool
}
