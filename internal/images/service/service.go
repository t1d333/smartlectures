package service

import (
	"context"
	"fmt"
	"io"

	"github.com/t1d333/smartlectures/internal/images"
	"github.com/t1d333/smartlectures/internal/images/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
}

func (s *Service) UploadImage(img io.Reader, ctx context.Context) (string, error) {
	src, err := s.repository.UploadImage(img, ctx)
	if err != nil {
		err = fmt.Errorf("failed upload image in service: %w", err)
	}

	return src, err
}

func NewService(logger logger.Logger, repository repository.Repository) images.Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}
