package notes

import (
	"context"

	"github.com/t1d333/smartlectures/internal/models"
)

type Service interface {
	GetNote(noteId int, ctx context.Context) (models.Note, error)
	CreateNote(note models.Note, ctx context.Context) (int, error)
	DeleteNote(noteId int, ctx context.Context) error
	UpdateNote(note models.Note, ctx context.Context) error
	GetNotesOverview(userId int, ctx context.Context) (models.NotesOverview, error)
	SearchNote(req models.SearchRequest, ctx context.Context) ([]models.NoteSearchItem, error)
}
