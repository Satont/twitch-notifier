package temporal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivity_ThumbnailCheckerTemporalActivityCorrect(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "image/png")
			},
		),
	)
	defer ts.Close()

	activity := NewActivity()
	err := activity.ThumbnailCheckerTemporalActivity(
		context.TODO(),
		ts.URL,
	)

	assert.NoError(t, err)
}

func TestActivity_ThumbnailCheckerTemporalActivityRedirect(t *testing.T) {
	t.Parallel()
	
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://google.com", http.StatusFound)
			},
		),
	)
	defer ts.Close()

	activity := NewActivity()
	err := activity.ThumbnailCheckerTemporalActivity(
		context.TODO(),
		ts.URL,
	)

	assert.ErrorIs(t, err, ErrInvalidThumbnail)
}

func TestActivity_ThumbnailCheckerTemporalActivityNotImage(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
			},
		),
	)
	defer ts.Close()

	activity := NewActivity()
	err := activity.ThumbnailCheckerTemporalActivity(
		context.TODO(),
		ts.URL,
	)

	assert.ErrorIs(t, err, ErrInvalidThumbnail)
}
