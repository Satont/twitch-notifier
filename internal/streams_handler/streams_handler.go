package streams_handler

import (
	"context"
)

type StreamsHandler interface {
	Online(ctx context.Context, opts ChannelOnlineOpts) error
	Offline(ctx context.Context, opts ChannelOfflineOpts) error
	MetadataChange(ctx context.Context, opts ChannelMetaDataChangedOpts) error
}

type ChannelOnlineOpts struct {
	ChannelID    string
	ThumbnailURL string
	Category     Category
	Title        string
}

type ChannelOfflineOpts struct {
	ChannelID string
}

type ChannelMetaDataChangedOpts struct {
	ChannelID   string
	OldTitle    string
	NewTitle    string
	OldCategory Category
	NewCategory Category
}

type Category struct {
	ID   string
	Name string
}
