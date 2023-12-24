package announcesender

import (
	"context"

	"github.com/google/uuid"
)

//go:generate go run go.uber.org/mock/mockgen -source=announce_sender.go -destination=mocks/mock.go

type AnnounceSender interface {
	SendOnline(ctx context.Context, opts ChannelOnlineOpts) error
	SendOffline(ctx context.Context, opts ChannelOfflineOpts) error
	SendTitleChange(ctx context.Context, opts ChannelTitleChangeOpts) error
	SendCategoryChange(ctx context.Context, opts ChannelCategoryChangeOpts) error
	SendTitleAndCategoryChange(ctx context.Context, opts ChannelTitleAndCategoryChangeOpts) error
}

type ChannelOnlineOpts struct {
	ChannelID    uuid.UUID
	Category     string
	Title        string
	ThumbnailURL string
}

type ChannelOfflineOpts struct {
	ChannelID uuid.UUID
}

type ChannelTitleChangeOpts struct {
	ChannelID uuid.UUID
	OldTitle  string
	NewTitle  string
}

type ChannelCategoryChangeOpts struct {
	ChannelID   uuid.UUID
	OldCategory string
	NewCategory string
}

type ChannelTitleAndCategoryChangeOpts struct {
	ChannelID   uuid.UUID
	OldCategory string
	NewCategory string
	OldTitle    string
	NewTitle    string
}
