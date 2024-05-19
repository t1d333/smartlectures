package grpc

import (
	"context"
	"fmt"

	"github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
	service "github.com/t1d333/smartlectures/internal/storage"
	storage "github.com/t1d333/smartlectures/internal/storage"
	"github.com/t1d333/smartlectures/pkg/logger"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Delivery struct {
	logger  logger.Logger
	service storage.Service
	storage.UnimplementedStorageServer
}

func NewDelivery(logger logger.Logger, service storage.Service) storage.StorageServer {
	return &Delivery{
		logger:  logger,
		service: service,
	}
}

func (*Delivery) GetDir(context.Context, *wrapperspb.Int32Value) (*storage.Dir, error) {
	panic("unimplemented")
}

func (d *Delivery) GetNote(ctx context.Context, id *wrapperspb.Int32Value) (*storage.Note, error) {
	note, err := d.service.GetNote(ctx, int(id.GetValue()))
	if err != nil {
		return nil, errors.ErrNoteNotFound
	}

	return &storage.Note{
		Id:   int32(note.NoteId),
		Name: note.Name,
		Body: note.Body,
	}, nil
}

// SearchDir implements service.StorageServer.
func (*Delivery) SearchDir(
	context.Context,
	*wrapperspb.StringValue,
) (*storage.SearchResponse, error) {
	panic("unimplemented")
}

func (d *Delivery) SearchSnippet(
	ctx context.Context,
	query *wrapperspb.StringValue,
) (*storage.SearchSnippetResponse, error) {
	searchItems, err := d.service.SearchSnippet(ctx, query.GetValue())
	if err != nil {
		return nil, fmt.Errorf("failed to search note: %w", err)
	}

	result := &storage.SearchSnippetResponse{
		Items: []*storage.Snippet{},
	}

	for _, item := range searchItems {
		result.Items = append(result.Items, &storage.Snippet{
			SnippetId:   int32(item.SnippetID),
			Name:        item.Name,
			Description: item.Description,
			Body:        item.Body,
			UserId:      int32(item.UserId),
		})
	}

	return result, nil
}

func (*Delivery) CreateDir(context.Context, *storage.Dir) (*status.Status, error) {
	panic("unimplemented")
}

func (d *Delivery) CreateNote(ctx context.Context, note *storage.Note) (*status.Status, error) {
	err := d.service.CreateNote(
		ctx,
		models.Note{
			UserId:    int(note.GetUserId()),
			NoteId:    int(note.GetId()),
			Body:      note.GetBody(),
			Name:      note.GetName(),
			ParentDir: int(note.GetParentDir()),
		},
	)
	if err != nil {
		return &status.Status{
			Code:    500,
			Message: "",
		}, nil
	}

	return &status.Status{Code: 204}, nil
}

func (d *Delivery) DeleteDir(
	ctx context.Context,
	v *wrapperspb.Int32Value,
) (*status.Status, error) {
	err := d.service.DeleteDir(ctx, int(v.GetValue()))
	if err != nil {

		d.logger.Error(err)
		return &status.Status{
			Code: 500,
		}, nil
	}

	return &status.Status{
		Code: 204,
	}, nil
}

func (d *Delivery) DeleteNote(
	ctx context.Context,
	id *wrapperspb.Int32Value,
) (*status.Status, error) {
	err := d.service.DeleteNote(ctx, int(id.GetValue()))
	if err != nil {
		d.logger.Error(err)
		return &status.Status{
			Code: 500,
		}, nil
	}

	return &status.Status{
		Code: 204,
	}, nil
}

func (d *Delivery) SearchNote(
	ctx context.Context,
	req *service.SearchRequest,
) (*storage.SearchResponse, error) {
	searchItems, err := d.service.SearchNote(
		context.WithValue(ctx, "userId", int(req.UserId)),
		req.GetQuery(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search note: %w", err)
	}

	result := &storage.SearchResponse{
		Items: []*storage.NoteSearchItem{},
	}

	for _, item := range searchItems {
		result.Items = append(result.Items, &storage.NoteSearchItem{
			Id:            int32(item.NoteID),
			Name:          item.Name,
			NameHighlight: item.NameHighlight,
			BodyHighlight: item.BodyHighlight,
		})
	}

	return result, nil
}

func (d *Delivery) UpdateNote(
	ctx context.Context,
	req *storage.NoteUpdateRequest,
) (*status.Status, error) {
	err := d.service.UpdateNote(
		ctx,
		models.Note{
			NoteId:    int(req.GetId()),
			Name:      req.GetName(),
			Body:      req.GetBody(),
			ParentDir: int(req.GetParentDir()),
		},
	)
	if err != nil {
		d.logger.Error(err)
		return &status.Status{Code: 500}, nil
	}

	return &status.Status{Code: 204}, nil
}
