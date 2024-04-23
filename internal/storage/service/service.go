package service

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/models"
	storage "github.com/t1d333/smartlectures/internal/storage"
	"github.com/t1d333/smartlectures/internal/storage/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
}

func (s *Service) SearchSnippet(ctx context.Context, query string) ([]models.Snippet, error) {
	result, err := s.repository.SearchSnippet(ctx, query)
	if err != nil {
		return []models.Snippet{}, fmt.Errorf("failed to search snippet in service: %w", err)
	}

	return result, nil
}

func (s *Service) DeleteDir(ctx context.Context, id int) error {
	err := s.repository.DeleteDir(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete dir in service: %w", err)
	}

	return nil
}

func (*Service) SearchDir(ctx context.Context, query string) ([]models.Dir, error) {
	panic("unimplemented")
	
}

func (s *Service) SearchNote(ctx context.Context, query string) ([]models.NoteSearchItem, error) {
	result, err := s.repository.SearchNote(ctx, query)
	if err != nil {
		return []models.NoteSearchItem{}, fmt.Errorf("failed to search note in service: %w", err)
	}

	return result, nil
}

func (s *Service) GetNote(ctx context.Context, id int) (models.Note, error) {
	note, err := s.repository.GetNote(ctx, id)
	if err != nil {
		return note, fmt.Errorf("failed to get note in service: %w", err)
	}

	return note, nil
}

func (s *Service) CreateNote(ctx context.Context, note models.Note) error {
	err := s.repository.CreateNote(ctx, note)
	if err != nil {
		return fmt.Errorf("failed to create note in service: %w", err)
	}

	return err
}

func (s *Service) DeleteNote(ctx context.Context, id int) error {
	err := s.repository.DeleteNote(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete note data in service: %w", err)
	}

	return nil
}

func (s *Service) UpdateNote(ctx context.Context, note models.Note) error {
	err := s.repository.UpdateNote(ctx, note)
	if err != nil {
		return fmt.Errorf("failed to update note data in service: %w", err)
	}

	return nil
}

func NewService(logger logger.Logger, rep repository.Repository) storage.Service {
	return &Service{
		logger:     logger,
		repository: rep,
	}
}
