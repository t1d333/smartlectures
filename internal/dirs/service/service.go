package service

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/dirs"
	"github.com/t1d333/smartlectures/internal/dirs/repository"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
)

const (
	newDirName = "Новая папка"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
}

func (s *Service) CreateDir(dir models.Dir, ctx context.Context) (int, error) {
	if dir.Name == "" {
		dir.Name = newDirName
	}

	dirId, err := s.repository.CreateDir(dir, ctx)
	if err != nil {
		err = fmt.Errorf("failed to create dir in dirs service: %w", err)
	}

	return dirId, err
}

func (s *Service) DeleteDir(dirId int, ctx context.Context) error {
	err := s.repository.DeleteDir(dirId, ctx)
	if err != nil {
		err = fmt.Errorf("failed to delete dir in dirs service: %w", err)
	}

	return err
}

func (s *Service) GetDir(dirId int, ctx context.Context) (models.Dir, error) {
	dir, err := s.repository.GetDir(dirId, ctx)
	if err != nil {
		err = fmt.Errorf("failed to get dir in dirs service: %w", err)
	}

	if dir.RepeatedNum != 0 {
		dir.Name = fmt.Sprintf("%s(%d)", dir.Name, dir.RepeatedNum)
	}

	return dir, err
}

func (s *Service) GetDirsOverview(userId int, ctx context.Context) (models.DirsOverview, error) {
	overview, err := s.repository.GetDirsOverview(userId, ctx)
	if err != nil {
		err = fmt.Errorf("failed to get dirs overview in dirs service: %w", err)
	}

	return overview, err
}

func (s *Service) UpdateDir(dir models.Dir, ctx context.Context) error {
	err := s.repository.UpdateDir(dir, ctx)
	if err != nil {
		err = fmt.Errorf("failed to update dir in dirs service: %w", err)
	}

	return err
}

func NewService(logger logger.Logger, repository repository.Repository) dirs.Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}
