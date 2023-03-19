package i18n

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNested(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"foo": "bar",
	}

	res, ok := GetNested[string](data, "foo")
	assert.True(t, ok, "expected to get a value")
	assert.Equal(t, "bar", res, "expected to get a value")

	res, ok = GetNested[string](data, "bar")
	assert.False(t, ok, "expected to be false")
	assert.Equal(t, "", res, "expected to not get a value")
}
