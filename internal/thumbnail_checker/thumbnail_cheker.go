package thumbnail_checker

import (
	"context"
)

type ThumbnailChecker interface {
	ValidateThumbnail(ctx context.Context, url string) error
}
