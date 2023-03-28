package twitch

import (
	"github.com/nicklaw5/helix/v2"
)

type Interface interface {
	GetUser(id, login string) (*helix.User, error)
	GetUsers(ids, logins []string) ([]helix.User, error)

	GetStreamByUserId(id string) (*helix.Stream, error)
	GetStreamsByUserIds(ids []string) ([]helix.Stream, error)

	GetChannelByUserId(id string) (*helix.ChannelInformation, error)
	GetChannelsByUserIds(ids []string) ([]helix.ChannelInformation, error)
}
