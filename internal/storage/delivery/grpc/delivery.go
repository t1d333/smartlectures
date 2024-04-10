package grpc

import (
	"context"

	"github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
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

// SearchNoteByBody implements service.StorageServer.
func (*Delivery) SearchNoteByBody(
	context.Context,
	*wrapperspb.StringValue,
) (*storage.SearchResponse, error) {
	panic("unimplemented")
}

// SearchNoteByName implements service.StorageServer.
func (*Delivery) SearchNoteByName(
	context.Context,
	*wrapperspb.StringValue,
) (*storage.SearchResponse, error) {
	panic("unimplemented")
}

func (*Delivery) CreateDir(context.Context, *storage.Dir) (*status.Status, error) {
	panic("unimplemented")
}

func (d *Delivery) CreateNote(ctx context.Context, note *storage.Note) (*status.Status, error) {
	err := d.service.CreateNote(
		ctx,
		models.Note{NoteId: int(note.Id), Body: note.Body, Name: note.Name},
	)
	if err != nil {
		return &status.Status{
			Code:    500,
			Message: "",
		}, nil
	}

	return &status.Status{Code: 204}, nil
}

func (*Delivery) DeleteDir(context.Context, *wrapperspb.Int32Value) (*status.Status, error) {
	panic("unimplemented")
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

func (*Delivery) Search(context.Context, *wrapperspb.StringValue) (*storage.SearchResponse, error) {
	panic("unimplemented")
}

func (d *Delivery) UpdateNote(
	ctx context.Context,
	req *storage.NoteUpdateRequest,
) (*status.Status, error) {
	err := d.service.UpdateNote(
		ctx,
		models.Note{NoteId: int(req.GetId()), Name: req.GetName(), Body: req.GetBody()},
	)
	if err != nil {
		d.logger.Error(err)
		return &status.Status{Code: 500}, nil
	}

	return &status.Status{Code: 204}, nil
}
