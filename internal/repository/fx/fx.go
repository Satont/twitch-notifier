package fx

import (
	"github.com/satont/twitch-notifier/internal/repository/channel"
	"github.com/satont/twitch-notifier/internal/repository/chat"
	"github.com/satont/twitch-notifier/internal/repository/chat_settings"
	"github.com/satont/twitch-notifier/internal/repository/follow"
	"github.com/satont/twitch-notifier/internal/repository/stream"
	"go.uber.org/fx"
)

var Module = fx.Options(
	channel.Module,
	chat.Module,
	chat_settings.Module,
	follow.Module,
	stream.Module,
)
