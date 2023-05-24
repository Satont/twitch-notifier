package twitch_streams_cheker

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type thumbNailBuilder struct {
}

func newThumbNailBuilder() *thumbNailBuilder {
	return &thumbNailBuilder{}
}

func (c *thumbNailBuilder) checkValidity(url string, n int) (bool, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := client.Get(url)
	if err != nil {
		return false, err
	}

	if req.StatusCode != 200 && n == 5 {
		return false, fmt.Errorf("url %s is not valid", url)
	} else if req.StatusCode != 200 {
		time.Sleep(5 * time.Second)
		return c.checkValidity(url, n+1)
	} else {
		return true, nil
	}
}

func (c *thumbNailBuilder) Build(thumbNailUrl string, checkValidity bool) (string, error) {
	thumbNail := thumbNailUrl
	thumbNail = strings.Replace(thumbNail, "{width}", "1920", 1)
	thumbNail = strings.Replace(thumbNail, "{height}", "1080", 1)

	if !checkValidity {
		return thumbNail, nil
	}

	valid, err := c.checkValidity(thumbNail, 0)

	if !valid || err != nil {
		thumbNail = strings.Replace(thumbNail, "1920", "1280", 1)
		thumbNail = strings.Replace(thumbNail, "1080", "720", 1)
	}

	if err != nil {
		return thumbNail, err
	}

	return thumbNail, nil
}
