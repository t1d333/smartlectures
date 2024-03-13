package notes

import "github.com/t1d333/smartlectures/internal/models"

type Service interface {
	GetNote(noteId int) (models.Note, error)
	CreateNote(note models.Note) error
	DeleteNote(noteId int) error
	UpdateNote(noteId int) error
	GetNotesOverview(userId int) error
}
