package temporal

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

func NewActivity() *Activity {
	cl := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Activity{
		client: cl,
	}
}

type Activity struct {
	client *http.Client
}

var ErrInvalidThumbnail = errors.New("invalid thumbnail")

func (c *Activity) ThumbnailCheckerTemporalActivity(
	ctx context.Context,
	thumbnailUrl string,
) error {
	u, err := url.Parse(thumbnailUrl)
	if err != nil {
		return err
	}

	request := &http.Request{
		URL: u,
	}
	request = request.WithContext(ctx)

	res, err := c.client.Do(request)
	if err != nil {
		return err
	}

	contentType := res.Header.Get("Content-Type")
	isImage := contentType == "image/png" || contentType == "image/jpeg"

	isNotRedirect := res.StatusCode >= 200 && res.StatusCode < 300

	if isImage && isNotRedirect {
		return nil
	}

	return ErrInvalidThumbnail
}
