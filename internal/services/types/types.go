package types

import (
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/twitch"
)

type Services struct {
	Twitch  twitch.Interface
	Chat    db.ChatInterface
	Channel db.ChannelInterface
	Follow  db.FollowInterface
}
