package service

import (
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/notes"
	"github.com/t1d333/smartlectures/internal/notes/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
}

func (Service) CreateNote(note models.Note) error {
	panic("unimplemented")
}

func (Service) DeleteNote(noteId int) error {
	panic("unimplemented")
}

func (Service) GetNote(noteId int) (models.Note, error) {
	panic("unimplemented")
}

func (Service) GetNotesOverview(userId int) error {
	panic("unimplemented")
}

func (Service) UpdateNote(noteId int) error {
	panic("unimplemented")
}

func NewService(logger logger.Logger, repository repository.Repository) notes.Service {
	return Service{
		logger:     logger,
		repository: repository,
	}
}
