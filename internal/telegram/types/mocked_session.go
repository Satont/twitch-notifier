package tg_types

import (
	"context"

	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/stretchr/testify/mock"
)

type MockedSessionManager[T Session] struct {
	mock.Mock
}

func (m *MockedSessionManager[T]) SetEqualFunc(fn func(t T, t2 T) bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MockedSessionManager[T]) Setup(opt session.ManagerOption, opts ...session.ManagerOption) {
	//TODO implement me
	panic("implement me")
}

func (m *MockedSessionManager[T]) Get(ctx context.Context) *T {
	args := m.Called(ctx)

	return args.Get(0).(*T)
}

func (m *MockedSessionManager[T]) Reset(session *T) {
	//TODO implement me
	panic("implement me")
}

func (m *MockedSessionManager[T]) Filter(fn func(t *T) bool) tgb.Filter {
	//TODO implement me
	panic("implement me")
}

func (m *MockedSessionManager[T]) Wrap(next tgb.Handler) tgb.Handler {
	//TODO implement me
	panic("implement me")
}

func NewMockedSessionManager() *MockedSessionManager[Session] {
	return &MockedSessionManager[Session]{}
}
