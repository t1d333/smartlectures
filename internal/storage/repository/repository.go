package repository

import (
	"context"

	"github.com/t1d333/smartlectures/internal/models"
)

type Repository interface {
	GetNote(ctx context.Context, id int) (models.Note, error)
	CreateNote(ctx context.Context, note models.Note) error
	UpdateNote(ctx context.Context, note models.Note) error
	DeleteNote(ctx context.Context, id int) error
	DeleteDir(ctx context.Context, id int) error
	SearchNote(ctx context.Context, query string, userId int) ([]models.NoteSearchItem, error)
	SearchSnippet(ctx context.Context, query string) ([]models.Snippet, error)
	SearchDir(ctx context.Context, query string) ([]int, error)
}
