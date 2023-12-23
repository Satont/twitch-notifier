// Code generated by MockGen. DO NOT EDIT.
// Source: chat.go
//
// Generated by this command:
//
//	mockgen -source=chat.go -destination=mocks/mock.go
//

// Package mock_chat is a generated GoMock package.
package mock_chat

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	domain "github.com/satont/twitch-notifier/internal/domain"
	chat "github.com/satont/twitch-notifier/internal/repository/chat"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, user domain.Chat) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, user)
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, id)
}

// GetAll mocks base method.
func (m *MockRepository) GetAll(ctx context.Context) ([]domain.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]domain.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockRepositoryMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRepository)(nil).GetAll), ctx)
}

// GetByChatServiceAndChatID mocks base method.
func (m *MockRepository) GetByChatServiceAndChatID(ctx context.Context, service chat.ChatService, chatID string) (*domain.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByChatServiceAndChatID", ctx, service, chatID)
	ret0, _ := ret[0].(*domain.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByChatServiceAndChatID indicates an expected call of GetByChatServiceAndChatID.
func (mr *MockRepositoryMockRecorder) GetByChatServiceAndChatID(ctx, service, chatID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByChatServiceAndChatID", reflect.TypeOf((*MockRepository)(nil).GetByChatServiceAndChatID), ctx, service, chatID)
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, id)
}