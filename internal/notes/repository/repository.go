package repository

import (
	"context"

	"github.com/t1d333/smartlectures/internal/models"
)

type Repository interface {
	GetNote(noteId int, ctx context.Context) (models.Note, error)
	CreateNote(note models.Note, ctx context.Context) (int, error)
	DeleteNote(noteId int, ctx context.Context) error
	UpdateNote(note models.Note, ctx context.Context) error
	GetNotesOverview(userId int, ctx context.Context) ([]models.NotePreview, error)
}
