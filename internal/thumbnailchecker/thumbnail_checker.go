package thumbnailchecker

import (
	"context"
)

//go:generate go run go.uber.org/mock/mockgen -source=thumbnail_checker.go -destination=mocks/mock.go

type ThumbnailChecker interface {
	ValidateThumbnail(ctx context.Context, url string) error
	TransformSizes(url string, width int, height int) string
}
