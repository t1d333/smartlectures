package service

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/notes"
	"github.com/t1d333/smartlectures/internal/notes/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"

	storage "github.com/t1d333/smartlectures/internal/storage"
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
	query models.SearchRequest,
	ctx context.Context,
) ([]models.NoteSearchItem, error) {
	searchResult, err := s.client.SearchNote(ctx, &wrapperspb.StringValue{Value: query.Query})
	if err != nil {
		return []models.NoteSearchItem{}, fmt.Errorf(
			"failed to get search result from storage service: %w",
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

func (s *Service) CreateNote(note models.Note, ctx context.Context) (int, error) {
	if note.Name == "" {
		note.Name = newNoteName
	}

	noteId, err := s.repository.CreateNote(note, ctx)
	if err != nil {
		err = fmt.Errorf("failed to create note in notes service: %w", err)
	}

	status, err := s.client.CreateNote(ctx, &storage.Note{
		Id:   int32(noteId),
		Name: note.Name,
		Body: note.Body,
	})

	if status.GetCode() != 204 {
		return 0, fmt.Errorf("failed to indexing note data: %w", err)
	}

	return noteId, err
}

func (s *Service) DeleteNote(noteId int, ctx context.Context) error {
	err := s.repository.DeleteNote(noteId, ctx)
	if err != nil {
		err = fmt.Errorf("failed to delete note in notes service: %w", err)
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

	return err
}

func (s *Service) GetNote(noteId int, ctx context.Context) (models.Note, error) {
	note, err := s.repository.GetNote(noteId, ctx)
	if err != nil {
		return note, fmt.Errorf("failed to get note in notes service: %w", err)
	}
	//
	// data, err := s.client.GetNote(ctx, &wrapperspb.Int32Value{
	// 	Value: int32(noteId),
	// })
	// if err != nil {
	// 	return note, fmt.Errorf("failed to get note data in service: %w", err)
	// }
	//
	// note.Name = data.Name
	// note.Body = data.Body

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
	if err := s.repository.UpdateNote(note, ctx); err != nil {
		return fmt.Errorf("failed to update note in notes service: %w", err)
	}

	status, err := s.client.UpdateNote(ctx, &storage.NoteUpdateRequest{
		Name: note.Name,
		Body: note.Body,
		Id:   int32(note.NoteId),
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
