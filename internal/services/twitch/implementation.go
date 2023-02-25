package twitch

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"time"
)

type twitchService struct {
	apiClient *helix.Client
}

func NewTwitchService(clientId string, clientSecret string) (TwitchService, error) {
	apiClient, err := helix.NewClient(&helix.Options{
		ClientID:      clientId,
		ClientSecret:  clientSecret,
		RateLimitFunc: rateLimitCallback,
	})

	if err != nil {
		return nil, err
	}

	token, err := apiClient.RequestAppAccessToken([]string{})
	if err != nil {
		panic(err)
	}
	apiClient.SetAppAccessToken(token.Data.AccessToken)

	go func() {
		for {
			newToken, tokenErr := apiClient.RequestAppAccessToken([]string{})
			if tokenErr != nil {
				panic(tokenErr)
			}
			apiClient.SetAppAccessToken(newToken.Data.AccessToken)
			time.Sleep(1 * time.Hour)
		}
	}()

	return &twitchService{
		apiClient: apiClient,
	}, nil
}

func (t *twitchService) GetUser(id, login string) (*helix.User, error) {
	users, err := t.GetUsers([]string{id}, []string{login})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

func (t *twitchService) GetUsers(ids, logins []string) ([]helix.User, error) {
	data := lo.If(ids[0] == "", logins).Else(ids)

	reqData := &chunkedRequestData[*helix.UsersParams, *helix.UsersResponse]{
		ids:       data,
		requestFn: t.apiClient.GetUsers,
		responseSelectorFn: func(response *helix.UsersResponse) interface{} {
			return response.Data.Users
		},
		paramFn: func(chunk []string) *helix.UsersParams {
			if ids[0] != "" {
				return &helix.UsersParams{
					IDs: chunk,
				}
			} else {
				return &helix.UsersParams{
					Logins: chunk,
				}
			}
		},
	}

	users, err := getDataChunked[helix.User](reqData)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (t *twitchService) GetStreamByUserId(id string) (*helix.Stream, error) {
	streams, err := t.GetStreamsByUserIds([]string{id})
	if err != nil {
		return nil, err
	}

	if len(streams) == 0 {
		return nil, nil
	}

	return &streams[0], nil
}

func (t *twitchService) GetStreamsByUserIds(ids []string) ([]helix.Stream, error) {
	reqData := &chunkedRequestData[*helix.StreamsParams, *helix.StreamsResponse]{
		ids:       ids,
		requestFn: t.apiClient.GetStreams,
		responseSelectorFn: func(response *helix.StreamsResponse) interface{} {
			return response.Data.Streams
		},
		paramFn: func(chunk []string) *helix.StreamsParams {
			return &helix.StreamsParams{
				UserIDs: chunk,
			}
		},
	}

	streams, err := getDataChunked[helix.Stream](reqData)
	if err != nil {
		return nil, err
	}

	return streams, nil
}

func (t *twitchService) GetChannelByUserId(id string) (*helix.ChannelInformation, error) {
	channels, err := t.GetChannelsByUserIds([]string{id})
	if err != nil {
		return nil, err
	}

	if len(channels) == 0 {
		return nil, nil
	}

	return &channels[0], nil
}

func (t *twitchService) GetChannelsByUserIds(ids []string) ([]helix.ChannelInformation, error) {
	reqData := &chunkedRequestData[*helix.GetChannelInformationParams, *helix.GetChannelInformationResponse]{
		ids:       ids,
		requestFn: t.apiClient.GetChannelInformation,
		responseSelectorFn: func(response *helix.GetChannelInformationResponse) interface{} {
			return response.Data.Channels
		},
		paramFn: func(chunk []string) *helix.GetChannelInformationParams {
			return &helix.GetChannelInformationParams{
				BroadcasterIDs: chunk,
			}
		},
	}

	channels, err := getDataChunked[helix.ChannelInformation](reqData)
	if err != nil {
		return nil, err
	}

	return channels, nil
}
