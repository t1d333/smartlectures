package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/t1d333/smartlectures/internal/dirs/repository"
	customErrors "github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Repository struct {
	logger logger.Logger
	pool   *pgxpool.Pool
}

func (r *Repository) CreateDir(dir models.Dir, ctx context.Context) (int, error) {
	row := r.pool.QueryRow(
		ctx,
		InsertNewDirCMD,
		dir.Name,
		dir.UserId,
		dir.ParentDir,
	)

	dirId := 0

	if err := row.Scan(&dirId); err != nil {
		r.logger.Errorf("failed to create dir in repository: %w", err)
		return dirId, fmt.Errorf("failed to create dir in repository: %w", err)
	}

	return dirId, nil
}

func (r *Repository) DeleteDir(dirId int, ctx context.Context) error {
	rows, _ := r.pool.Query(
		ctx,
		DeleteDirCMD,
		dirId,
	)

	defer rows.Close()

	if err := rows.Scan(); err != nil && !errors.Is(err, pgx.ErrNoRows) {

		// TODO: добавить центральное логирование
		r.logger.Errorf("failed to delete dir in repository: %w", err)

		// TODO: исправить на custom err
		return fmt.Errorf("failed to delete dir in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return customErrors.ErrDirNotFound
	}

	return nil
}

func (r *Repository) GetDir(dirId int, ctx context.Context) (models.Dir, error) {
	dir := models.Dir{}
	row := r.pool.QueryRow(ctx, SelectDirByIDCMD, dirId)

	parentDir := sql.NullInt32{}

	if err := row.Scan(&dir.DirId, &dir.Name, &dir.UserId, &parentDir); err != nil &&
		!errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("failed to get dir in repository: %w", err)
		return dir, fmt.Errorf("failed to get note in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return dir, customErrors.ErrDirNotFound
	}

	if parentDir.Valid {
		dir.ParentDir = int(parentDir.Int32)
	}

	// TODO: else

	return dir, nil
}

func (*Repository) GetDirsOverview(userId int, ctx context.Context) (models.DirsOverview, error) {
	panic("unimplemented")
}

func (r *Repository) UpdateDir(dir models.Dir, ctx context.Context) error {
	row := r.pool.QueryRow(
		ctx,
		UpdateDirCMD,
		dir.DirId,
		dir.Name,
		dir.ParentDir,
	)

	tmp := 0

	if err := row.Scan(&tmp); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("failed to update dir in repository: %w", err)
		return fmt.Errorf("failed to update dir in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		r.logger.Info(1235)
		return customErrors.ErrDirNotFound
	}

	return nil
}

func NewRepository(logger logger.Logger, pool *pgxpool.Pool) repository.Repository {
	return &Repository{
		logger: logger,
		pool:   pool,
	}
}
