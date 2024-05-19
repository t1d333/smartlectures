package service

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/notes"
	"github.com/t1d333/smartlectures/internal/notes/repository"
	storage "github.com/t1d333/smartlectures/internal/storage"
	"github.com/t1d333/smartlectures/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	newNoteName = "Новый файл"
)

type Service struct {
	logger     logger.Logger
	repository repository.Repository
	client     storage.StorageClient
}

func (s *Service) SearchNote(
	ctx context.Context,
	query models.SearchRequest,
) ([]models.NoteSearchItem, error) {
	userId := ctx.Value("userId").(int)

	searchResult, err := s.client.SearchNote(ctx, &storage.SearchRequest{
		Query:  query.Query,
		UserId: int32(userId),
	})
	if err != nil {
		return []models.NoteSearchItem{}, fmt.Errorf(
			"NotesService.SearchNote(query: %v): %w",
			query,
			err,
		)
	}

	result := []models.NoteSearchItem{}

	for _, item := range searchResult.Items {
		result = append(result, models.NoteSearchItem{
			NoteID:        int(item.GetId()),
			Name:          item.GetName(),
			BodyHighlight: item.GetBodyHighlight(),
			NameHighlight: item.GetNameHighlight(),
		})
	}

	return result, nil
}

func (s *Service) CreateNote(ctx context.Context, note models.Note) (int, error) {
	if note.Name == "" {
		note.Name = newNoteName
	}

	noteId, err := s.repository.CreateNote(ctx, note)
	if err != nil {
		return noteId, fmt.Errorf("failed to create note in notes service: %w", err)
	}

	s.logger.Error(note.ParentDir)

	status, err := s.client.CreateNote(ctx, &storage.Note{
		Id:        int32(noteId),
		UserId:    int32(note.UserId),
		Name:      note.Name,
		Body:      note.Body,
		ParentDir: int32(note.ParentDir),
	})

	if status.GetCode() != 204 {
		return 0, fmt.Errorf("failed to indexing note data: %w", err)
	}

	return noteId, err
}

func (s *Service) DeleteNote(ctx context.Context, noteId int) error {
	userId := ctx.Value("userId")

	if note, err := s.repository.GetNote(ctx, noteId); err != nil {
		return fmt.Errorf("NotesService.UpdateNote(note: %v)", note)
	} else if note.UserId != userId {
		return errors.ErrPermissionDenied
	}

	err := s.repository.DeleteNote(ctx, noteId)
	if err != nil {
		return fmt.Errorf("failed to delete note in notes service: %w", err)
	}

	status, err := s.client.DeleteNote(ctx, &wrapperspb.Int32Value{Value: int32(noteId)})
	if err != nil || status.GetCode() != 204 {
		s.logger.Errorw(
			"failed to delete note data index from storage",
			"id",
			noteId,
			"err",
			err,
			"response",
			status.GetMessage(),
		)
	}

	return nil
}

func (s *Service) GetNote(ctx context.Context, noteId int) (models.Note, error) {
	note, err := s.repository.GetNote(ctx, noteId)
	userId := ctx.Value("userId").(int)

	if err != nil {
		return models.Note{}, fmt.Errorf("failed to get note in notes service: %w", err)
	}

	if note.UserId != userId {
		fmt.Println(note.UserId, userId)
		return models.Note{}, errors.ErrPermissionDenied
	}
	return note, err
}

func (s *Service) GetNotesOverview(ctx context.Context, userId int) (models.NotesOverview, error) {
	notes, err := s.repository.GetNotesOverview(ctx, userId)
	if err != nil {
		return models.NotesOverview{}, fmt.Errorf(
			"failed to get notes overview in notes service: %w",
			err,
		)
	}

	return models.NotesOverview{Notes: notes}, nil
}

func (s *Service) UpdateNote(ctx context.Context, note models.Note) error {
	userId := ctx.Value("userId")

	if note, err := s.repository.GetNote(ctx, note.NoteId); err != nil {
		return fmt.Errorf("NotesService.UpdateNote(note: %v)", note)
	} else if note.UserId != userId {
		return errors.ErrPermissionDenied
	}

	if err := s.repository.UpdateNote(ctx, note); err != nil {
		return fmt.Errorf("failed to update note in notes service: %w", err)
	}

	status, err := s.client.UpdateNote(ctx, &storage.NoteUpdateRequest{
		Name:      note.Name,
		Body:      note.Body,
		Id:        int32(note.NoteId),
		ParentDir: int32(note.ParentDir),
	})
	if status.Code != 204 || err != nil {
		return fmt.Errorf("failed to update note data in storage: %w", err)
	}

	return nil
}

func NewService(
	logger logger.Logger,
	repository repository.Repository,
	client storage.StorageClient,
) notes.Service {
	return &Service{
		logger:     logger,
		repository: repository,
		client:     client,
	}
}
