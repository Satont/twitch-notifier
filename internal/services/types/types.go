package types

import (
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/internal/services/twitch"
	"github.com/satont/twitch-notifier/pkg/i18n"
)

type Services struct {
	Twitch  twitch.Interface
	Chat    db.ChatInterface
	Channel db.ChannelInterface
	Follow  db.FollowInterface
	I18N    *i18n.I18n
}
