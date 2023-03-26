package twitch

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockedApi(server *httptest.Server) (*twitchService, error) {
	apiClient, err := helix.NewClient(&helix.Options{
		ClientID:   "test",
		APIBaseURL: server.URL,
	})
	if err != nil {
		return nil, err
	}

	apiClient.SetAppAccessToken("test")

	return &twitchService{
		apiClient: apiClient,
	}, nil
}

func TestTwitchService_GetUser(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"id":"1","login":"test"}]}`))
	}))
	defer server.Close()

	twitchService, err := newMockedApi(server)
	assert.NoError(t, err)

	user, err := twitchService.GetUser("1", "")
	assert.NoError(t, err)

	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "test", user.Login)
}

func TestTwitchService_GetUsers(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"id":"1","login":"test"},{"id":"2","login":"test2"}]}`))
	}))
	defer server.Close()

	twitchService, err := newMockedApi(server)
	assert.NoError(t, err)

	expectedUsers := []helix.User{
		{ID: "1", Login: "test"},
		{ID: "2", Login: "test2"},
	}

	table := []struct {
		name   string
		ids    []string
		logins []string
	}{
		{
			name:   "ids",
			ids:    []string{"1", "2"},
			logins: []string{},
		},
		{
			name:   "logins",
			ids:    []string{},
			logins: []string{"test", "test2"},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			users, err := twitchService.GetUsers(tt.ids, tt.logins)
			assert.NoError(t, err)

			assert.Equal(t, expectedUsers, users)
		})
	}
}

func TestTwitchService_GetStreamByUserId(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"id":"1","user_name":"test","game_name": "Dota 2"}]}`))
	}))
	defer server.Close()

	twitchService, err := newMockedApi(server)
	assert.NoError(t, err)

	stream, err := twitchService.GetStreamByUserId("1")
	assert.NoError(t, err)

	assert.Equal(t, "1", stream.ID)
	assert.Equal(t, "test", stream.UserName)
	assert.Equal(t, "Dota 2", stream.GameName)
}

func TestTwitchService_GetStreamsByUserId(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"id":"1","user_name":"test","game_name": "Dota 2"}, {"id":"2","user_name":"test2","game_name": "Dota 3"}]}`))
	}))
	defer server.Close()

	twitchService, err := newMockedApi(server)
	assert.NoError(t, err)

	streams, err := twitchService.GetStreamsByUserIds([]string{"1", "2"})
	assert.NoError(t, err)

	assert.Equal(t, []helix.Stream{
		{ID: "1", UserName: "test", GameName: "Dota 2"},
		{ID: "2", UserName: "test2", GameName: "Dota 3"},
	}, streams)
}

func TestTwitchService_GetChannelByUserId(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[
		{"broadcaster_id":"1","broadcaster_name":"test","game_name": "Dota 2", "title": "tiitle"}
		]}`))
	}))
	defer server.Close()

	twitchService, err := newMockedApi(server)
	assert.NoError(t, err)

	channel, err := twitchService.GetChannelByUserId("1")
	assert.NoError(t, err)

	assert.Equal(t, "1", channel.BroadcasterID)
	assert.Equal(t, "test", channel.BroadcasterName)
	assert.Equal(t, "Dota 2", channel.GameName)
	assert.Equal(t, "tiitle", channel.Title)
}

func TestTwitchService_GetChannelsByUserId(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[
		{"broadcaster_id":"1","broadcaster_name":"test","game_name": "Dota 2", "title": "tiitle"},
		{"broadcaster_id":"2","broadcaster_name":"test2","game_name": "Dota 3", "title": "tiitle2"}
		]}`))
	}))
	defer server.Close()

	twitchService, err := newMockedApi(server)
	assert.NoError(t, err)

	channels, err := twitchService.GetChannelsByUserIds([]string{"1", "2"})
	assert.NoError(t, err)

	assert.Equal(t, []helix.ChannelInformation{
		{BroadcasterID: "1", BroadcasterName: "test", GameName: "Dota 2", Title: "tiitle"},
		{BroadcasterID: "2", BroadcasterName: "test2", GameName: "Dota 3", Title: "tiitle2"},
	}, channels)
}
