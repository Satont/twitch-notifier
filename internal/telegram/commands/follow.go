package commands

import (
	"context"
	"errors"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/satont/twitch-notifier/internal/db/db_models"
	tgtypes "github.com/satont/twitch-notifier/internal/telegram/types"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

type FollowCommand struct {
	*tgtypes.CommandOpts
}

var (
	twitchInvalidNamesString = "Invalid login names, emails or IDs in request"
	channelNotFoundError     = errors.New("channel not found")
	invalidNameError         = errors.New(twitchInvalidNamesString)
	TwitchLinkRegular        = regexp.MustCompile(`(?:https?://)?(?:www\.)?twitch\.tv/(\w+)`)
)

func (c *FollowCommand) createFollow(ctx context.Context, chat *db_models.Chat, input string) (*db_models.Follow, error) {
	twitchChannel, err := c.Services.Twitch.GetUser("", input)
	if err != nil {
		if err.Error() == twitchInvalidNamesString {
			return nil, invalidNameError
		}

		return nil, err
	}

	if twitchChannel == nil {
		return nil, channelNotFoundError
	}

	dbChannel, err := c.Services.Channel.GetByIdOrCreate(ctx, twitchChannel.ID, db_models.ChannelServiceTwitch)
	if err != nil {
		return nil, err
	}

	follow, err := c.Services.Follow.Create(ctx, dbChannel.ID, chat.ID)
	if err != nil {
		return nil, err
	}

	return follow, nil
}

func (c *FollowCommand) handleScene(ctx context.Context, msg *tgb.MessageUpdate) error {
	chat := c.SessionManager.Get(ctx).Chat

	nicknames := make([]string, 0)

	regularMatches := TwitchLinkRegular.FindAllStringSubmatch(msg.Text, -1)

	if len(regularMatches) > 0 {
		for _, match := range regularMatches {
			nicknames = append(nicknames, match[1])
		}
	} else {
		nicknames = append(nicknames, msg.Text)
	}

	succeeded := make([]string, 0)
	failed := make([]string, 0)

	for _, nickname := range nicknames {
		_, err := c.createFollow(ctx, chat, nickname)

		if errors.Is(err, channelNotFoundError) {
			message := c.Services.I18N.Translate(
				"commands.follow.errors.streamerNotFound",
				chat.Settings.ChatLanguage.String(),
				map[string]string{
					"streamer": nickname,
				},
			)
			failed = append(failed, message)
		} else if errors.Is(err, db_models.FollowAlreadyExistsError) {
			message := c.Services.I18N.Translate(
				"commands.follow.errors.alreadyFollowed",
				chat.Settings.ChatLanguage.String(),
				map[string]string{
					"streamer": nickname,
				},
			)
			failed = append(failed, message)
		} else if errors.Is(err, invalidNameError) {
			message := c.Services.I18N.Translate(
				"commands.follow.errors.badUsername",
				chat.Settings.ChatLanguage.String(),
				map[string]string{
					"streamer": nickname,
				},
			)
			failed = append(failed, message)
		} else if err != nil {
			zap.S().Error(err)
			failed = append(failed, "internal error")
		} else {
			message := c.Services.I18N.Translate(
				"commands.follow.success",
				chat.Settings.ChatLanguage.String(),
				map[string]string{
					"streamer": nickname,
				},
			)
			succeeded = append(succeeded, message)
		}
	}

	c.SessionManager.Get(ctx).Scene = ""

	message := strings.Join(succeeded, "\n")
	message += "\n\n"
	message += strings.Join(failed, "\n")

	return msg.Answer(message).DoVoid(ctx)
}

func (c *FollowCommand) HandleCommand(ctx context.Context, msg *tgb.MessageUpdate) error {
	session := c.SessionManager.Get(ctx)

	text := strings.ReplaceAll(msg.Text, "/follow", "")
	text = strings.TrimSpace(text)

	if text != "" {
		msg.Text = text
		return c.handleScene(ctx, msg)
	} else {
		c.SessionManager.Get(ctx).Scene = "follow"
		return msg.
			Answer(c.Services.I18N.Translate(
				"commands.follow.enter",
				session.Chat.Settings.ChatLanguage.String(),
				nil,
			)).
			DoVoid(ctx)
	}
}

var (
	followCommandQuery = tgb.Command("follow")
)

func NewFollowCommand(opts *tgtypes.CommandOpts) {
	cmd := &FollowCommand{
		CommandOpts: opts,
	}

	opts.Router.Message(cmd.handleScene, channelsAdminFilter, tgb.FilterFunc(func(ctx context.Context, update *tgb.Update) (bool, error) {
		session := opts.SessionManager.Get(ctx)
		return session.Scene == "follow", nil
	}))
	opts.Router.Message(cmd.HandleCommand, channelsAdminFilter, followCommandQuery)
}
