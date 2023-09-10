package types

import (
	"github.com/satont/twitch-notifier/internal/config"
	"github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/message_sender"
	"github.com/satont/twitch-notifier/internal/twitch"
	"github.com/satont/twitch-notifier/pkg/i18n"
)

type Services struct {
	Config        *config.Config
	Twitch        twitch.Interface
	Chat          db.ChatInterface
	Channel       db.ChannelInterface
	Follow        db.FollowInterface
	Stream        db.StreamInterface
	QueueJob      db.QueueJobInterface
	I18N          i18n.Interface
	MessageSender message_sender.MessageSenderInterface
}
