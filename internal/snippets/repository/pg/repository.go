package pg

import (
	"context"
	"database/sql"
	// "errors"
	"fmt"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	// customErrors "github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/snippets/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Repository struct {
	logger logger.Logger
	pool   *pgxpool.Pool
}

// CreateSnippet implements repository.Repository.
func (*Repository) CreateSnippet(snippet models.Snippet, ctx context.Context) (int, error) {
	panic("unimplemented")
}

// DeleteSnippet implements repository.Repository.
func (*Repository) DeleteSnippet(snippetId int, ctx context.Context) error {
	panic("unimplemented")
}

// GetSnippets implements repository.Repository.
func (r *Repository) GetSnippets(userId int, ctx context.Context) ([]models.Snippet, error) {
	snippets := make([]models.Snippet, 0)

	rows, _ := r.pool.Query(ctx, SelectUserSnippets, userId)
	defer rows.Close()

	userIdTmp := sql.NullInt32{}
	for rows.Next() {
		snippet := models.Snippet{}
		if err := rows.Scan(&snippet.SnippetID, &snippet.Name, &snippet.Description, &snippet.Body, &userIdTmp); err != nil {
			r.logger.Errorf("failed to scan user snippet in snippets repository: %w", err)
			return snippets, fmt.Errorf(
				"failed to scan user snippet in snippets repository: %w",
				err,
			)
		}

		if userIdTmp.Valid {
			snippet.UserId = int(userIdTmp.Int32)
		}

		snippets = append(snippets, snippet)
	}

	return snippets, nil
}

// SearchSnippet implements repository.Repository.
func (*Repository) SearchSnippet(
	req models.SearchRequest,
	ctx context.Context,
) ([]models.NoteSearchItem, error) {
	panic("unimplemented")
}

// UpdateSnippet implements repository.Repository.
func (*Repository) UpdateSnippet(note models.Snippet, ctx context.Context) error {
	panic("unimplemented")
}

func NewRepository(logger logger.Logger, pool *pgxpool.Pool) repository.Repository {
	return &Repository{
		logger: logger,
		pool:   pool,
	}
}
