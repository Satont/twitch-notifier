package i18n

import "github.com/stretchr/testify/mock"

type I18nMock struct {
	mock.Mock
}

func (m *I18nMock) Translate(key, language string, data map[string]string) string {
	args := m.Called(key, language, data)
	return args.String(0)
}

func (m *I18nMock) GetLanguagesCodes() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func NewI18nMock() *I18nMock {
	return &I18nMock{}
}
