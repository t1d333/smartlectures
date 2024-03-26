package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	customErrors "github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/internal/notes/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Repository struct {
	logger logger.Logger
	pool   *pgxpool.Pool
}

func (r *Repository) CreateNote(note models.Note, ctx context.Context) (int, error) {
	row := r.pool.QueryRow(
		ctx,
		InsertNewNoteCMD,
		note.Name,
		note.Body,
		note.ParentDir,
		note.UserId,
	)

	noteId := 0

	if err := row.Scan(&noteId); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("failed to create note in repository: %w", err)
		return noteId, fmt.Errorf("failed to create note in repository: %w", err)
	}

	return noteId, nil
}

func (r *Repository) DeleteNote(noteId int, ctx context.Context) error {
	rows, _ := r.pool.Query(
		ctx,
		DeleteNodeCMD,
		noteId,
	)

	defer rows.Close()

	if err := rows.Scan(); err != nil && !errors.Is(err, pgx.ErrNoRows) {

		// TODO: добавить центральное логирование
		r.logger.Errorf("failed to delete note in repository: %w", err)

		// TODO: исправить на custom err
		return fmt.Errorf("failed to delete note in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return customErrors.ErrNoteNotFound
	}

	return nil
}

func (r *Repository) GetNote(noteId int, ctx context.Context) (models.Note, error) {
	note := models.Note{}

	row := r.pool.QueryRow(ctx, SelectNoteByIDCMD, noteId)

	parentDir := sql.NullInt32{}

	if err := row.Scan(&note.NoteId, &note.Name, &note.Body, &note.CreatedAt, &note.LastUpdate, &parentDir, &note.UserId); err != nil &&
		!errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("failed to get note in repository: %w", err)
		return note, fmt.Errorf("failed to get note in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return note, customErrors.ErrNoteNotFound
	}

	if parentDir.Valid {
		note.ParentDir = int(parentDir.Int32)
	}

	return note, nil
}

func (r *Repository) GetNotesOverview(
	userId int,
	ctx context.Context,
) ([]models.NotePreview, error) {
	overview := make([]models.NotePreview, 0)

	rows, _ := r.pool.Query(ctx, SelectUserNotesOverview, userId)
	defer rows.Close()

	parentDir := sql.NullInt32{}
	for rows.Next() {
		note := models.NotePreview{}
		if err := rows.Scan(&note.NoteId, &note.Name, &parentDir); err != nil {
			r.logger.Errorf("failed to scan user note preview in notes repository: %w", err)
			return overview, fmt.Errorf(
				"failed to scan user note preview in notes repository: %w",
				err,
			)
		}

		if parentDir.Valid {
			note.ParentDir = int(parentDir.Int32)
		}

		overview = append(overview, note)
	}

	return overview, nil
}

func (r *Repository) UpdateNote(note models.Note, ctx context.Context) error {
	rows := r.pool.QueryRow(
		ctx,
		UpdateNoteCMD,
		note.NoteId,
		note.Name,
		note.Body,
		note.ParentDir,
	)

	tmp := 0

	if err := rows.Scan(&tmp); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("failed to update note in repository: %w", err)
		return fmt.Errorf("failed to update note in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return customErrors.ErrNoteNotFound
	}

	return nil
}

func NewRepository(logger logger.Logger, pool *pgxpool.Pool) repository.Repository {
	return &Repository{
		logger: logger,
		pool:   pool,
	}
}
