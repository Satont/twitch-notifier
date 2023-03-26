package i18n

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestNewI18n(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	localesPath := filepath.Join(wd, "test_locales")

	table := []struct {
		translation string
		lang        string
		data        map[string]string
		expected    string
		expectErr   bool
		localesPath string
	}{
		{
			translation: "hello",
			lang:        "en",
			data:        nil,
			expected:    "world",
			localesPath: localesPath,
		},
		{
			translation: "nested.world.hello",
			lang:        "en",
			data:        nil,
			expected:    "nested world",
			localesPath: localesPath,
		},
		{
			translation: "templated",
			lang:        "en",
			data: map[string]string{
				"hello": "templated",
			},
			expected:    "hello templated",
			localesPath: localesPath,
		},
		{
			translation: "expectEmptyString",
			lang:        "en",
			data:        nil,
			expected:    "",
			localesPath: localesPath,
		},
		{
			translation: "expect error",
			expectErr:   true,
			localesPath: "/tmp/somefreakingstupidnotifierlocalespath",
		},
	}

	for _, tt := range table {
		t.Run(tt.translation, func(t *testing.T) {
			i18, err := NewI18n(tt.localesPath)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t,
				tt.expected,
				i18.Translate(tt.translation, tt.lang, tt.data),
			)
		})
	}
}
