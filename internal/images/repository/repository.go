package repository

import (
	"context"
	"io"
)


type Repository interface {
	UploadImage(image io.Reader, ctx context.Context) (string, error)
}
