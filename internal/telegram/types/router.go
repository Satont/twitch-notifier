package tg_types

import (
	"context"

	"github.com/mr-linch/go-tg/tgb"
	"github.com/stretchr/testify/mock"
)

type Router interface {
	Use(mws ...tgb.Middleware) *tgb.Router
	Message(handler tgb.MessageHandler, filters ...tgb.Filter) *tgb.Router
	EditedMessage(handler tgb.MessageHandler, filters ...tgb.Filter) *tgb.Router
	ChannelPost(handler tgb.MessageHandler, filters ...tgb.Filter) *tgb.Router
	EditedChannelPost(handler tgb.MessageHandler, filters ...tgb.Filter) *tgb.Router
	InlineQuery(handler tgb.InlineQueryHandler, filters ...tgb.Filter) *tgb.Router
	ChosenInlineResult(handler tgb.ChosenInlineResultHandler, filters ...tgb.Filter) *tgb.Router
	CallbackQuery(handler tgb.CallbackQueryHandler, filters ...tgb.Filter) *tgb.Router
	ShippingQuery(handler tgb.ShippingQueryHandler, filters ...tgb.Filter) *tgb.Router
	PreCheckoutQuery(handler tgb.PreCheckoutQueryHandler, filters ...tgb.Filter) *tgb.Router
	Poll(handler tgb.PollHandler, filters ...tgb.Filter) *tgb.Router
	PollAnswer(handler tgb.PollAnswerHandler, filters ...tgb.Filter) *tgb.Router
	MyChatMember(handler tgb.ChatMemberUpdatedHandler, filters ...tgb.Filter) *tgb.Router
	ChatMember(handler tgb.ChatMemberUpdatedHandler, filters ...tgb.Filter) *tgb.Router
	ChatJoinRequest(handler tgb.ChatJoinRequestHandler, filters ...tgb.Filter) *tgb.Router
	Error(handler tgb.ErrorHandler) *tgb.Router
	Update(handler tgb.HandlerFunc, filters ...tgb.Filter) *tgb.Router
	Handle(ctx context.Context, update *tgb.Update) error
}

type MockedRouter struct {
	mock.Mock
}

func (m *MockedRouter) Use(mws ...tgb.Middleware) *tgb.Router {
	args := m.Called(mws)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) Message(handler tgb.MessageHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) EditedMessage(handler tgb.MessageHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) ChannelPost(handler tgb.MessageHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) EditedChannelPost(handler tgb.MessageHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) InlineQuery(handler tgb.InlineQueryHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) ChosenInlineResult(handler tgb.ChosenInlineResultHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) CallbackQuery(handler tgb.CallbackQueryHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) ShippingQuery(handler tgb.ShippingQueryHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) PreCheckoutQuery(handler tgb.PreCheckoutQueryHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) Poll(handler tgb.PollHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) PollAnswer(handler tgb.PollAnswerHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) MyChatMember(handler tgb.ChatMemberUpdatedHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) ChatMember(handler tgb.ChatMemberUpdatedHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) ChatJoinRequest(handler tgb.ChatJoinRequestHandler, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) Error(handler tgb.ErrorHandler) *tgb.Router {
	args := m.Called(handler)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) Update(handler tgb.HandlerFunc, filters ...tgb.Filter) *tgb.Router {
	args := m.Called(handler, filters)

	return args.Get(0).(*tgb.Router)
}

func (m *MockedRouter) Handle(ctx context.Context, update *tgb.Update) error {
	args := m.Called(ctx, update)

	return args.Error(0)
}
