package thumbnailchecker

import (
	"context"
)

type ThumbnailChecker interface {
	ValidateThumbnail(ctx context.Context, url string) error
}
