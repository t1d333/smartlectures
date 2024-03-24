package images

import (
	"context"
	"io"
)

type Service interface {
	UploadImage(img io.Reader, ctx context.Context) (string, error)
}
