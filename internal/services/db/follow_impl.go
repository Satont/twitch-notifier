package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
)

type followService struct {
	entClient *ent.Client
}

func (f followService) Create(
	ctx context.Context,
	channelID uuid.UUID,
	chatID uuid.UUID,
) (*ent.Follow, error) {
	follow, err := f.entClient.Follow.
		Create().
		SetChatID(chatID).
		SetChannelID(channelID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return follow, nil
}

func (f followService) Delete(ctx context.Context, followID uuid.UUID) error {
	err := f.entClient.Follow.
		DeleteOneID(followID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewFollowService(entClient *ent.Client) FollowInterface {
	return &followService{entClient: entClient}
}
