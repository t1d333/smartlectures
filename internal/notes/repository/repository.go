package repository

import (
	"context"

	"github.com/t1d333/smartlectures/internal/models"
)

type Repository interface {
	GetNote(ctx context.Context, noteId int) (models.Note, error)
	CreateNote(ctx context.Context, note models.Note) (int, error)
	DeleteNote(ctx context.Context, noteId int) error
	UpdateNote(ctx context.Context, note models.Note) error
	GetNotesOverview(ctx context.Context, userId int) ([]models.NotePreview, error)
}
