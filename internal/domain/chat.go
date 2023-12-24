package domain

import (
	"github.com/google/uuid"
)

type Chat struct {
	ID      uuid.UUID
	Service ChatService
	ChatID  string
}
