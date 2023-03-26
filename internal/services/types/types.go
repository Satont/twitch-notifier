package types

import (
	"github.com/satont/twitch-notifier/internal/services/config"
	"github.com/satont/twitch-notifier/internal/services/db"
	"github.com/satont/twitch-notifier/pkg/i18n"

	"github.com/satont/twitch-notifier/internal/services/message_sender"
	"github.com/satont/twitch-notifier/internal/services/twitch"
)

type Services struct {
	Config        *config.Config
	Twitch        twitch.Interface
	Chat          db.ChatInterface
	Channel       db.ChannelInterface
	Follow        db.FollowInterface
	Stream        db.StreamInterface
	I18N          i18n.Interface
	MessageSender message_sender.MessageSenderInterface
}
