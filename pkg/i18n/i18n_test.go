package i18n

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestNewI18n(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tt, err := NewI18n(filepath.Join(wd, "test_locales"))
	if err != nil {
		t.Fatal(err)
	}

	table := []struct {
		translation string
		lang        string
		data        map[string]string
		expected    string
	}{
		{
			translation: "hello",
			lang:        "en",
			data:        nil,
			expected:    "world",
		},
		{
			translation: "nested.world.hello",
			lang:        "en",
			data:        nil,
			expected:    "nested world",
		},
		{
			translation: "templated",
			lang:        "en",
			data: map[string]string{
				"hello": "templated",
			},
			expected: "hello templated",
		},
		{
			translation: "expectEmptyString",
			lang:        "en",
			data:        nil,
			expected:    "",
		},
	}

	for _, test := range table {
		assert.Equal(t,
			test.expected,
			tt.Translate(test.translation, test.lang, test.data),
		)
	}
}
