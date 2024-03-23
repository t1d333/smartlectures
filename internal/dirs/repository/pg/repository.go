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
		r.logger.Errorf("failed to create dir in repository", err)
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
		r.logger.Errorf("failed to delete dir in repository", err)

		// TODO: исправить на custom err
		return fmt.Errorf("failed to delete dir in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return customErrors.ErrDirNotFound
	}

	return nil
}

func (r *Repository) GetDir(dirId int, ctx context.Context) (models.Dir, error) {
	result := models.Dir{}
	row := r.pool.QueryRow(ctx, SelectDirByIDCMD, dirId)

	parentDir := sql.NullInt32{}

	if err := row.Scan(&result.DirId, &result.Name, &result.UserId, &parentDir); err != nil &&
		!errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("failed to get dir in repository", err)
		return result, fmt.Errorf("failed to get dir in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return result, customErrors.ErrDirNotFound
	}

	if parentDir.Valid {
		result.ParentDir = int(parentDir.Int32)
	}

	rows, _ := r.pool.Query(ctx, SelectSubdirs, dirId)
	tmp := 0

	for rows.Next() {
		if err := rows.Scan(&tmp); err != nil {
			r.logger.Errorf("failed to get subdirs in repository", err)
			return result, fmt.Errorf("failed to get subdirs in repository: %w", err)
		}

		dir, err := r.GetDir(tmp, ctx)
		if err != nil {
			r.logger.Errorf("failed to get subdir in repository", err)
			return dir, fmt.Errorf("failed to get subdir in repository: %w", err)
		}

		result.Subdirs = append(result.Subdirs, dir)
	}

	return result, nil
}

func (r *Repository) GetDirsOverview(userId int, ctx context.Context) (models.DirsOverview, error) {
	result := models.DirsOverview{}

	rows, _ := r.pool.Query(ctx, SelectUserDirsOverview, userId)
	tmp := 0

	for rows.Next() {
		if err := rows.Scan(&tmp); err != nil {
			r.logger.Errorf("failed to get user dirs overview in repository", err)
			return result, fmt.Errorf("failed to get user dirs overview in repository: %w", err)
		}

		dir, err := r.GetDir(tmp, ctx)
		if err != nil {
			r.logger.Errorf("failed to get dir in repository", err)
			return result, fmt.Errorf("failed to get dir in repository: %w", err)
		}

		result.Dirs = append(result.Dirs, dir)
	}

	return result, nil
}

func (r *Repository) UpdateDir(dir models.Dir, ctx context.Context) error {
	r.logger.Info(dir)
	row := r.pool.QueryRow(
		ctx,
		UpdateDirCMD,
		dir.DirId,
		dir.Name,
		dir.ParentDir,
	)

	tmp := 0

	if err := row.Scan(&tmp); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("failed to update dir in repository", err)
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
