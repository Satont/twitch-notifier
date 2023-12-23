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

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	return errors.New("invalid thumbnail")
}
