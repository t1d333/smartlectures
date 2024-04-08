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
	Search(ctx context.Context, query string) error
}
