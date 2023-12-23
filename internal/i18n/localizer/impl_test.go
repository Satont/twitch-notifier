package localizer

import (
	"testing"

	"github.com/satont/twitch-notifier/internal/domain"
	"github.com/satont/twitch-notifier/internal/i18n/store"
	"github.com/satont/twitch-notifier/internal/i18n/store/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestImpl_Localize(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	store := mock_store.NewMockI18nStore(ctrl)

	impl := NewLocalizer(store)

	table := []struct {
		name        string
		language    domain.Language
		key         string
		attrs       map[string]any
		expected    string
		expectedErr error
		mock        func()
	}{
		{
			name:        "empty key",
			language:    domain.LanguageEN,
			key:         "",
			attrs:       nil,
			expected:    "",
			expectedErr: ErrKeyIsEmpty,
		},
		{
			name:        "correct replace channel name",
			language:    domain.LanguageEN,
			key:         "test",
			attrs:       map[string]any{"channel": "notifier"},
			expected:    "test notifier",
			expectedErr: nil,
			mock: func() {
				store.EXPECT().GetKey(domain.LanguageEN, "test").Return("test {{ channel }}", nil)
			},
		},
		{
			name:        "should not replace unknown attribute",
			language:    domain.LanguageEN,
			key:         "test",
			attrs:       nil,
			expected:    "test {{ channel }}",
			expectedErr: nil,
			mock: func() {
				store.EXPECT().GetKey(domain.LanguageEN, "test").Return("test {{ channel }}", nil)
			},
		},
	}

	for _, tt := range table {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				if tt.mock != nil {
					tt.mock()
				}

				opts := []Option{
					WithLanguage(tt.language),
					WithKey(tt.key),
				}
				for k, v := range tt.attrs {
					opts = append(opts, WithAttribute(k, v))
				}

				res, err := impl.Localize(opts...)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, res)
				}
			},
		)
	}
}

func TestImpl_MustLocalize(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockedStore := mock_store.NewMockI18nStore(ctrl)

	impl := NewLocalizer(mockedStore)

	table := []struct {
		name        string
		language    domain.Language
		key         string
		attrs       map[string]any
		expected    string
		expectedErr error
		mock        func()
		shouldPanic bool
	}{
		{
			name:        "should take default language string",
			language:    domain.LanguageRU,
			key:         "test",
			attrs:       nil,
			expected:    "test {{ channel }}",
			expectedErr: nil,
			mock: func() {
				mockedStore.EXPECT().GetKey(domain.LanguageRU, "test").Return("", store.ErrKeyNotFound)
				mockedStore.EXPECT().GetKey(domain.LanguageEN, "test").Return("test {{ channel }}", nil)
			},
		},
		{
			name:     "should panic if default language string not found",
			language: domain.LanguageRU,
			key:      "test",
			mock: func() {
				mockedStore.EXPECT().GetKey(domain.LanguageRU, "test").Return("", store.ErrKeyNotFound)
				mockedStore.EXPECT().GetKey(domain.LanguageEN, "test").Return("", store.ErrKeyNotFound)
			},
			shouldPanic: true,
		},
	}

	for _, tt := range table {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				if tt.mock != nil {
					tt.mock()
				}

				if tt.shouldPanic {
					assert.Panics(
						t,
						func() {
							opts := []Option{
								WithLanguage(tt.language),
								WithKey(tt.key),
							}
							for k, v := range tt.attrs {
								opts = append(opts, WithAttribute(k, v))
							}

							impl.MustLocalize(opts...)
						},
					)
					return
				}

				opts := []Option{
					WithLanguage(tt.language),
					WithKey(tt.key),
				}
				for k, v := range tt.attrs {
					opts = append(opts, WithAttribute(k, v))
				}

				res := impl.MustLocalize(opts...)
				assert.Equal(t, tt.expected, res)
			},
		)
	}
}
