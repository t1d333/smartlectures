package pg

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/notes/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Repository struct {
	logger logger.Logger
	pool   *pgxpool.Pool
}

// CreateNote implements repository.Repository.
func (*Repository) CreateNote(note models.Note) error {
	panic("unimplemented")
}

// DeleteNote implements repository.Repository.
func (*Repository) DeleteNote(noteId int) error {
	panic("unimplemented")
}

// GetNote implements repository.Repository.
func (*Repository) GetNote(noteId int) (models.Note, error) {
	panic("unimplemented")
}

// GetNotesOverview implements repository.Repository.
func (*Repository) GetNotesOverview(userId int) error {
	panic("unimplemented")
}

// UpdateNote implements repository.Repository.
func (*Repository) UpdateNote(noteId int) error {
	panic("unimplemented")
}

func NewRepository(logger logger.Logger, pool *pgxpool.Pool) repository.Repository {
	return &Repository{
		logger: logger,
		pool:   pool,
	}
}
