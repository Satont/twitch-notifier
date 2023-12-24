package temporal

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/announcesender"
	"github.com/satont/twitch-notifier/internal/domain"
	mocklocalizer "github.com/satont/twitch-notifier/internal/i18n/localizer/mocks"
	mockmessagesender "github.com/satont/twitch-notifier/internal/messagesender/mocks"
	mockchat "github.com/satont/twitch-notifier/internal/repository/chat/mocks"
	mockchatsettings "github.com/satont/twitch-notifier/internal/repository/chatsettings/mocks"
	mockfollow "github.com/satont/twitch-notifier/internal/repository/follow/mocks"
	mockthumbnailchecker "github.com/satont/twitch-notifier/internal/thumbnailchecker/mocks"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"go.uber.org/mock/gomock"
)

func Test_OnlineWorkflow(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	localizer := mocklocalizer.NewMockLocalizer(ctrl)
	messageSender := mockmessagesender.NewMockMessageSender(ctrl)
	thumbnailChecker := mockthumbnailchecker.NewMockThumbnailChecker(ctrl)
	followRepository := mockfollow.NewMockRepository(ctrl)
	chatRepository := mockchat.NewMockRepository(ctrl)
	chatSettingsRepository := mockchatsettings.NewMockRepository(ctrl)

	workflow := &Workflow{
		localizer:              localizer,
		messageSender:          messageSender,
		thumbnailChecker:       thumbnailChecker,
		followRepository:       followRepository,
		chatRepository:         chatRepository,
		chatSettingsRepository: chatSettingsRepository,
	}

	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	channel := &domain.Channel{
		ID:        uuid.New(),
		ChannelID: "123",
		Service:   domain.StreamingServiceTwitch,
	}

	followerChat := &domain.Chat{
		ID:      uuid.New(),
		Service: domain.ChatServiceTelegram,
		ChatID:  "11111",
	}
	followers := []domain.Follow{
		{
			ID:        uuid.New(),
			ChatID:    followerChat.ID,
			ChannelID: channel.ID,
			CreatedAt: time.Now(),
		},
	}
	chatSettings := &domain.ChatSettings{
		ID:                            uuid.New(),
		ChatID:                        followerChat.ID,
		Language:                      domain.LanguageEN,
		CategoryChangeNotifications:   true,
		TitleChangeNotifications:      true,
		OfflineNotifications:          true,
		CategoryAndTitleNotifications: true,
		ShowThumbnail:                 true,
	}

	thumbnailChecker.EXPECT().ValidateThumbnail(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	followRepository.EXPECT().GetByChannelID(gomock.Any(), channel.ID).Return(
		followers,
		nil,
	).AnyTimes()
	chatRepository.EXPECT().GetByID(gomock.Any(), followers[0].ChatID).Return(
		followerChat,
		nil,
	).AnyTimes()
	chatSettingsRepository.EXPECT().GetByChatID(gomock.Any(), followerChat.ID).Return(
		chatSettings,
		nil,
	).AnyTimes()
	localizer.EXPECT().MustLocalize(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
		"Hello world",
	).AnyTimes()
	messageSender.EXPECT().SendMessageTelegram(gomock.Any(), gomock.Any()).Return(nil)

	env.ExecuteWorkflow(
		workflow.SendOnline,
		announcesender.ChannelOnlineOpts{
			ChannelID:    channel.ID,
			Category:     "Dota 2",
			Title:        "Hello world",
			ThumbnailURL: "https://twitch.tv/notifier",
		},
	)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	messageSender.EXPECT().SendMessageTelegram(
		gomock.Any(),
		gomock.Any(),
	).Return(errors.New("test"))

	testSuite = &testsuite.WorkflowTestSuite{}
	env = testSuite.NewTestWorkflowEnvironment()

	env.ExecuteWorkflow(
		workflow.SendOnline,
		announcesender.ChannelOnlineOpts{
			ChannelID:    channel.ID,
			Category:     "Dota 2",
			Title:        "Hello world",
			ThumbnailURL: "https://twitch.tv/notifier",
		},
	)

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}
