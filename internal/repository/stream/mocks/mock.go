// Code generated by MockGen. DO NOT EDIT.
// Source: stream.go
//
// Generated by this command:
//
//	mockgen -source=stream.go -destination=mocks/mock.go
//

// Package mock_stream is a generated GoMock package.
package mock_stream

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	domain "github.com/satont/twitch-notifier/internal/domain"
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
func (m *MockRepository) Create(ctx context.Context, stream domain.Stream) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, stream)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, stream any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, stream)
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

// GetByChannelId mocks base method.
func (m *MockRepository) GetByChannelId(ctx context.Context, channelId uuid.UUID) ([]domain.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByChannelId", ctx, channelId)
	ret0, _ := ret[0].([]domain.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByChannelId indicates an expected call of GetByChannelId.
func (mr *MockRepositoryMockRecorder) GetByChannelId(ctx, channelId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByChannelId", reflect.TypeOf((*MockRepository)(nil).GetByChannelId), ctx, channelId)
}

// GetById mocks base method.
func (m *MockRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id)
	ret0, _ := ret[0].(*domain.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockRepositoryMockRecorder) GetById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockRepository)(nil).GetById), ctx, id)
}

// GetLatestByChannelId mocks base method.
func (m *MockRepository) GetLatestByChannelId(ctx context.Context, channelId uuid.UUID) (*domain.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestByChannelId", ctx, channelId)
	ret0, _ := ret[0].(*domain.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestByChannelId indicates an expected call of GetLatestByChannelId.
func (mr *MockRepositoryMockRecorder) GetLatestByChannelId(ctx, channelId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestByChannelId", reflect.TypeOf((*MockRepository)(nil).GetLatestByChannelId), ctx, channelId)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, stream domain.Stream) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, stream)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, stream any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, stream)
}