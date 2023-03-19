package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent"
	"github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/satont/twitch-notifier/ent/follow"
	"github.com/satont/twitch-notifier/internal/services/db/db_models"
)

type followService struct {
	entClient *ent.Client
}

func (f *followService) convertEntity(follow *ent.Follow) *db_models.Follow {
	convertedFollow := &db_models.Follow{
		ID: follow.ID,
	}

	if follow.Edges.Channel != nil {
		convertedFollow.ChannelID = follow.Edges.Channel.ID
	}

	if follow.Edges.Chat != nil {
		convertedFollow.ChatID = follow.Edges.Chat.ID
	}

	return convertedFollow
}

func (f *followService) Create(
	ctx context.Context,
	channelID uuid.UUID,
	chatID uuid.UUID,
) (*db_models.Follow, error) {
	_, err := f.entClient.Follow.
		Create().
		SetChatID(chatID).
		SetChannelID(channelID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return f.GetByChatAndChannel(ctx, channelID, chatID)
}

func (f *followService) Delete(ctx context.Context, followID uuid.UUID) error {
	err := f.entClient.Follow.
		DeleteOneID(followID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (f *followService) GetByChatAndChannel(
	ctx context.Context,
	channelID uuid.UUID,
	chatID uuid.UUID,
) (*db_models.Follow, error) {
	fol, err := f.entClient.Follow.
		Query().
		Where(follow.ChannelID(channelID), follow.ChatID(chatID)).
		WithChannel().
		WithChat().
		First(ctx)

	if err != nil && ent.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if fol == nil {
		return nil, nil
	}

	return f.convertEntity(fol), err
}

func (f *followService) GetByChannelID(ctx context.Context, channelID uuid.UUID) ([]*db_models.Follow, error) {
	follows, err := f.entClient.Follow.
		Query().
		Where(follow.HasChannelWith(channel.IDEQ(channelID))).
		WithChannel().
		WithChat().
		All(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]*db_models.Follow, len(follows))
	for i, foll := range follows {
		result[i] = f.convertEntity(foll)
	}

	return result, nil
}

func (f *followService) GetByChatID(ctx context.Context, chatID uuid.UUID) ([]*db_models.Follow, error) {
	follows, err := f.entClient.Follow.
		Query().
		Where(follow.HasChatWith(chat.IDEQ(chatID))).
		WithChat().
		WithChannel().
		All(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]*db_models.Follow, len(follows))
	for i, foll := range follows {
		result[i] = f.convertEntity(foll)
	}

	return result, nil
}

func NewFollowService(entClient *ent.Client) FollowInterface {
	return &followService{entClient: entClient}
}
