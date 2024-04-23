package service

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/snippets"
	"github.com/t1d333/smartlectures/internal/snippets/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"

	storage "github.com/t1d333/smartlectures/internal/storage"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
	client     storage.StorageClient
}

func (*Service) CreateSnippet(snippet models.Snippet, ctx context.Context) (int, error) {
	panic("unimplemented")
}

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

func (s *Service) SearchSnippet(
	query string,
	ctx context.Context,
) ([]models.Snippet, error) {
	searchResult, err := s.client.SearchSnippet(ctx, &wrapperspb.StringValue{Value: query})
	if err != nil {
		return []models.Snippet{}, fmt.Errorf(
			"failed to get search result from storage service: %w",
			err,
		)
	}

	result := []models.Snippet{}

	for _, item := range searchResult.Items {
		result = append(result,
			models.Snippet{
				SnippetID:   int(item.SnippetId),
				Name:        item.Name,
				Description: item.Description,
				Body:        item.Body,
				UserId:      int(item.UserId),
			},
		)
	}

	return result, nil
}

func (*Service) UpdateSnippet(note models.Snippet, ctx context.Context) error {
	panic("unimplemented")
}

func NewService(
	logger logger.Logger,
	repository repository.Repository,
	client storage.StorageClient,
) snippets.Service {
	return &Service{
		logger:     logger,
		repository: repository,
		client: client,
	}
}
