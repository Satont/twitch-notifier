package fx

import (
	"github.com/satont/twitch-notifier/internal/repository/channel"
	"github.com/satont/twitch-notifier/internal/repository/chat"
	"github.com/satont/twitch-notifier/internal/repository/chatsettings"
	"github.com/satont/twitch-notifier/internal/repository/follow"
	"github.com/satont/twitch-notifier/internal/repository/stream"
	"go.uber.org/fx"
)

var Module = fx.Options(
	channel.Module,
	chat.Module,
	chatsettings.Module,
	follow.Module,
	stream.Module,
)
