package snippets

import (
	"context"

	"github.com/t1d333/smartlectures/internal/models"
)

type Service interface {
	GetSnippets(userId int, ctx context.Context) ([]models.Snippet, error)
	CreateSnippet(snippet models.Snippet, ctx context.Context) (int, error)
	DeleteSnippet(snippetId int, ctx context.Context) error
	UpdateSnippet(note models.Snippet, ctx context.Context) error
	SearchSnippet(query string, ctx context.Context) ([]models.Snippet, error)
}
