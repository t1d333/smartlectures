package service

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/dirs"
	"github.com/t1d333/smartlectures/internal/dirs/repository"
	"github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"

	storage "github.com/t1d333/smartlectures/internal/storage"
)

const (
	newDirName = "Новая папка"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
	client     storage.StorageClient
}

func (s *Service) CreateDir(ctx context.Context, dir models.Dir) (int, error) {
	if dir.Name == "" {
		dir.Name = newDirName
	}

	dirId, err := s.repository.CreateDir(dir, ctx)
	if err != nil {
		err = fmt.Errorf("failed to create dir in dirs service: %w", err)
	}

	return dirId, err
}

func (s *Service) DeleteDir(ctx context.Context, dirId int) error {
	dir, err := s.GetDir(ctx, dirId)

	if err != nil {
		return fmt.Errorf("DirsService.UpdateDir(dir: %v): %w", dir, err)
	} else if dir.UserId != ctx.Value("userId") {
		return errors.ErrPermissionDenied
	}

	
	err = s.repository.DeleteDir(dirId, ctx)
	if err != nil {
		return fmt.Errorf("failed to delete dir in dirs service: %w", err)
	}

	status, _ := s.client.DeleteDir(ctx, &wrapperspb.Int32Value{Value: int32(dirId)})

	if status.Code != 204 {
		return fmt.Errorf("failed to delete dir from index: %s", status.Message)
	}

	return nil
}

func (s *Service) GetDir(ctx context.Context, dirId int) (models.Dir, error) {
	dir, err := s.repository.GetDir(dirId, ctx)
	if dir.UserId != ctx.Value("userId") {
		return models.Dir{}, errors.ErrPermissionDenied
	}

	if err != nil {
		err = fmt.Errorf("failed to get dir in dirs service: %w", err)
	}

	return dir, err
}

func (s *Service) GetDirsOverview(ctx context.Context, userId int) (models.DirsOverview, error) {
	overview, err := s.repository.GetDirsOverview(userId, ctx)
	if err != nil {
		err = fmt.Errorf("failed to get dirs overview in dirs service: %w", err)
	}

	return overview, err
}

func (s *Service) UpdateDir(ctx context.Context, dir models.Dir) error {
	dir, err := s.GetDir(ctx, dir.DirId)

	if err != nil {
		return fmt.Errorf("DirsService.UpdateDir(dir: %v): %w", dir, err)
	} else if dir.UserId != ctx.Value("userId") {
		return errors.ErrPermissionDenied
	}

	err = s.repository.UpdateDir(dir, ctx)
	if err != nil {
		err = fmt.Errorf("failed to update dir in dirs service: %w", err)
	}

	return err
}

func NewService(
	logger logger.Logger,
	repository repository.Repository,
	client storage.StorageClient,
) dirs.Service {
	return &Service{
		logger:     logger,
		repository: repository,

		client: client,
	}
}
