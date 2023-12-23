package thumbnailchecker

import (
	"context"
)

type ThumbnailChecker interface {
	ValidateThumbnail(ctx context.Context, url string) error
	TransformSizes(url string, width int, height int) string
}
