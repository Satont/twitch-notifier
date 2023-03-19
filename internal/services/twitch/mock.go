package twitch

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/stretchr/testify/mock"
	"strings"
)

type MockedService struct {
	mock.Mock
}

func (m *MockedService) GetUser(id, login string) (*helix.User, error) {
	args := m.Called(id, login)
	return args.Get(0).(*helix.User), args.Error(1)
}

func (m *MockedService) GetUsers(ids, logins []string) ([]helix.User, error) {
	args := m.Called(ids, logins)
	return args.Get(0).([]helix.User), args.Error(1)
}

func (m *MockedService) GetStreamByUserId(id string) (*helix.Stream, error) {
	args := m.Called(id)
	return args.Get(0).(*helix.Stream), args.Error(1)
}

func (m *MockedService) GetStreamsByUserIds(ids []string) ([]helix.Stream, error) {
	args := m.Called(ids)
	return args.Get(0).([]helix.Stream), args.Error(1)
}

func (m *MockedService) GetChannelByUserId(id string) (*helix.ChannelInformation, error) {
	strings.ReplaceAll(id, " ", "")
	args := m.Called(id)
	return args.Get(0).(*helix.ChannelInformation), args.Error(1)
}

func (m *MockedService) GetChannelsByUserIds(ids []string) ([]helix.ChannelInformation, error) {
	args := m.Called(ids)
	return args.Get(0).([]helix.ChannelInformation), args.Error(1)
}