package twitch

import (
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"time"
)

func rateLimitCallback(lastResponse *helix.Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	var reset64 int64
	reset64 = int64(lastResponse.GetRateLimitReset())

	currentTime := time.Now().Unix()

	if currentTime < reset64 {
		timeDiff := time.Duration(reset64 - currentTime)
		if timeDiff > 0 {
			fmt.Printf("Waiting on rate limit to pass before sending next request (%d seconds)\n", timeDiff)
			time.Sleep(timeDiff * time.Second)
		}
	}

	return nil
}
