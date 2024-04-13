package service

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/snippets"
	"github.com/t1d333/smartlectures/internal/snippets/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
	// client     storage.StorageClient
}

// CreateSnippet implements snippets.Service.
func (*Service) CreateSnippet(snippet models.Snippet, ctx context.Context) (int, error) {
	panic("unimplemented")
}

// DeleteSnippet implements snippets.Service.
func (*Service) DeleteSnippet(snippetId int, ctx context.Context) error {
	panic("unimplemented")
}

func (s *Service) GetSnippets(userId int, ctx context.Context) ([]models.Snippet, error) {
	snippets, err := s.repository.GetSnippets(userId, ctx)
	if err != nil {
		return []models.Snippet{}, fmt.Errorf(
			"failed to get snippets in snippets service: %w",
			err,
		)
	}

	return snippets, nil
}

func (*Service) SearchSnippet(
	req models.SearchRequest,
	ctx context.Context,
) ([]models.NoteSearchItem, error) {
	panic("unimplemented")
}

func (*Service) UpdateSnippet(note models.Snippet, ctx context.Context) error {
	panic("unimplemented")
}

func NewService(
	logger logger.Logger,
	repository repository.Repository,
) snippets.Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}
