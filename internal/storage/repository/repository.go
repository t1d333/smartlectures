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
	SearchNoteByName(ctx context.Context, query string) ([]int, error)
	SearchNoteByBody(ctx context.Context, query string) ([]int, error)
	SearchDir(ctx context.Context, query string) ([]int, error)
}
