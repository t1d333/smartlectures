package service

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/notes"
	"github.com/t1d333/smartlectures/internal/notes/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

const (
	newNoteName = "Новый файл"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
}

func (s *Service) CreateNote(note models.Note, ctx context.Context) (int, error) {
	if note.Name == "" {
		note.Name = newNoteName
	}

	noteId, err := s.repository.CreateNote(note, ctx)
	if err != nil {
		err = fmt.Errorf("failed to create note in notes service: %w", err)
	}

	return noteId, err
}

func (s *Service) DeleteNote(noteId int, ctx context.Context) error {
	err := s.repository.DeleteNote(noteId, ctx)
	if err != nil {
		err = fmt.Errorf("failed to delete note in notes service: %w", err)
	}

	return err
}

func (s *Service) GetNote(noteId int, ctx context.Context) (models.Note, error) {
	note, err := s.repository.GetNote(noteId, ctx)
	if err != nil {
		err = fmt.Errorf("failed to get note in notes service: %w", err)
	} 

	return note, err
}

func (s *Service) GetNotesOverview(userId int, ctx context.Context) (models.NotesOverview, error) {
	notes, err := s.repository.GetNotesOverview(userId, ctx)
	if err != nil {
		return models.NotesOverview{}, fmt.Errorf(
			"failed to get notes overview in notes service: %w",
			err,
		)
	}

	return models.NotesOverview{Notes: notes}, nil
}

func (s *Service) UpdateNote(note models.Note, ctx context.Context) error {
	err := s.repository.UpdateNote(note, ctx)
	if err != nil {
		err = fmt.Errorf("failed to update note in notes service: %w", err)
	}

	return err
}

func NewService(logger logger.Logger, repository repository.Repository) notes.Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}
