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
		translation   string
		lang          string
		data          map[string]string
		expected      string
		expectErr     bool
		localesPath   string
		patchReadFile bool
		patchReadDir  bool
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
		{
			translation:   "expect readFile error",
			expectErr:     true,
			patchReadFile: true,
		},
		{
			translation:  "expect readDir error",
			expectErr:    true,
			patchReadDir: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.translation, func(t *testing.T) {
			if tt.patchReadFile {
				readFile = func(string) ([]byte, error) {
					return nil, os.ErrNotExist
				}
				defer func() { readFile = os.ReadFile }()
			}

			if tt.patchReadDir {
				readDir = func(string) ([]os.DirEntry, error) {
					return nil, os.ErrNotExist
				}
				defer func() { readDir = os.ReadDir }()
			}

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

func TestGetLanguagesCodes(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	i18, err := NewI18n(filepath.Join(wd, "test_locales"))

	assert.Equal(t, []string{"en"}, i18.GetLanguagesCodes())
}
