package tg_types

import (
	"context"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/stretchr/testify/mock"
)

type Router interface {
	Use(mws ...tgb.Middleware) Router
	Message(handler tgb.MessageHandler, filters ...tgb.Filter) Router
	EditedMessage(handler tgb.MessageHandler, filters ...tgb.Filter) Router
	ChannelPost(handler tgb.MessageHandler, filters ...tgb.Filter) Router
	EditedChannelPost(handler tgb.MessageHandler, filters ...tgb.Filter) Router
	InlineQuery(handler tgb.InlineQueryHandler, filters ...tgb.Filter) Router
	ChosenInlineResult(handler tgb.ChosenInlineResultHandler, filters ...tgb.Filter) Router
	CallbackQuery(handler tgb.CallbackQueryHandler, filters ...tgb.Filter) Router
	ShippingQuery(handler tgb.ShippingQueryHandler, filters ...tgb.Filter) Router
	PreCheckoutQuery(handler tgb.PreCheckoutQueryHandler, filters ...tgb.Filter) Router
	Poll(handler tgb.PollHandler, filters ...tgb.Filter) Router
	PollAnswer(handler tgb.PollAnswerHandler, filters ...tgb.Filter) Router
	MyChatMember(handler tgb.ChatMemberUpdatedHandler, filters ...tgb.Filter) Router
	ChatMember(handler tgb.ChatMemberUpdatedHandler, filters ...tgb.Filter) Router
	ChatJoinRequest(handler tgb.ChatJoinRequestHandler, filters ...tgb.Filter) Router
	Error(handler tgb.ErrorHandler) Router
	Update(handler tgb.HandlerFunc, filters ...tgb.Filter) Router
	Handle(ctx context.Context, update *tgb.Update) error
}

type MockedRouter struct {
	mock.Mock
}

func (m *MockedRouter) Use(mws ...tgb.Middleware) Router {
	args := m.Called(mws)

	return args.Get(0).(Router)
}

func (m *MockedRouter) Message(handler tgb.MessageHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) EditedMessage(handler tgb.MessageHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) ChannelPost(handler tgb.MessageHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) EditedChannelPost(handler tgb.MessageHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) InlineQuery(handler tgb.InlineQueryHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) ChosenInlineResult(handler tgb.ChosenInlineResultHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) CallbackQuery(handler tgb.CallbackQueryHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) ShippingQuery(handler tgb.ShippingQueryHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) PreCheckoutQuery(handler tgb.PreCheckoutQueryHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) Poll(handler tgb.PollHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) PollAnswer(handler tgb.PollAnswerHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) MyChatMember(handler tgb.ChatMemberUpdatedHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) ChatMember(handler tgb.ChatMemberUpdatedHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) ChatJoinRequest(handler tgb.ChatJoinRequestHandler, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) Error(handler tgb.ErrorHandler) Router {
	args := m.Called(handler)

	return args.Get(0).(Router)
}

func (m *MockedRouter) Update(handler tgb.HandlerFunc, filters ...tgb.Filter) Router {
	args := m.Called(handler, filters)

	return args.Get(0).(Router)
}

func (m *MockedRouter) Handle(ctx context.Context, update *tgb.Update) error {
	args := m.Called(ctx, update)

	return args.Error(0)
}
