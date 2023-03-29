package types

import (
	"github.com/satont/twitch-notifier/internal/config"
	db2 "github.com/satont/twitch-notifier/internal/db"
	"github.com/satont/twitch-notifier/internal/message_sender"
	"github.com/satont/twitch-notifier/internal/twitch"
	"github.com/satont/twitch-notifier/pkg/i18n"
)

type Services struct {
	Config        *config.Config
	Twitch        twitch.Interface
	Chat          db2.ChatInterface
	Channel       db2.ChannelInterface
	Follow        db2.FollowInterface
	Stream        db2.StreamInterface
	I18N          i18n.Interface
	MessageSender message_sender.MessageSenderInterface
}
