package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	autherrors "github.com/t1d333/smartlectures/internal/auth/errors"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type PostgresRepository struct {
	logger logger.Logger
	pool   *pgxpool.Pool
}

func New(logger logger.Logger, client *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		logger: logger,
		pool:   client,
	}
}

func (r *PostgresRepository) RegisterUser(
	ctx context.Context,
	user models.User,
) (models.User, error) {
	cmd := `INSERT INTO users(username, name, email, surname, password) VALUES ($1, $2, $3, $4, $5) RETURNING user_id;`
	res := r.pool.QueryRow(
		ctx,
		cmd,
		user.Username,
		user.Name,
		user.Email,
		user.Surname,
		user.Password,
	)

	if err := res.Scan(&user.UserId); err != nil {
		return user, fmt.Errorf("PostgresRepository.RegisterUser(user: %v): %w", user, err)
	}

	return user, nil
}

func (r *PostgresRepository) GetOnboardNote(
	ctx context.Context,
	userId int,
) (models.Note, error) {
	cmd := `SELECT note_id, name, body, created_at, last_update, parent_dir, user_id, repeated_num
					FROM notes
					WHERE user_id = $1 AND name=$2;`

	note := models.Note{}
	name := "Введение в приложение"

	row := r.pool.QueryRow(ctx, cmd, userId, name)

	parentDir := sql.NullInt32{}

	if err := row.Scan(&note.NoteId, &note.Name, &note.Body, &note.CreatedAt, &note.LastUpdate, &parentDir, &note.UserId, &note.RepeatedNum); err != nil &&
		!errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("failed to get note in repository: %w", err)
		return note, fmt.Errorf("failed to get note in repository: %w", err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return note, fmt.Errorf("failed to get note in repository: %w", err)
	}

	if parentDir.Valid {
		note.ParentDir = int(parentDir.Int32)
	}

	return note, nil
}

func (r *PostgresRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	user := models.User{}

	cmd := `SELECT user_id, username, name, surname, password FROM users WHERE email = $1`

	row := r.pool.QueryRow(ctx, cmd, email)

	if err := row.Scan(&user.UserId, &user.Username, &user.Name, &user.Surname, &user.Password); err != nil &&
		!errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("PostgresRepository.GetUser(email: %s): %w", email, err)
		return user, fmt.Errorf("PostgresRepository.GetUser(email: %s): %w", email, err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return user, autherrors.ErrUserNotFound
	}

	return user, nil
}

func (r *PostgresRepository) GetUserById(ctx context.Context, id int) (models.User, error) {
	user := models.User{}

	cmd := `SELECT user_id, username, name, surname, password, email FROM users WHERE user_id = $1`

	row := r.pool.QueryRow(ctx, cmd, id)

	if err := row.Scan(&user.UserId, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email); err != nil &&
		!errors.Is(err, pgx.ErrNoRows) {
		r.logger.Errorf("PostgresRepository.GetUserById(id: %d): %w", id, err)
		return user, fmt.Errorf("PostgresRepository.GetUserById(id: %d): %w", id, err)
	} else if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return user, autherrors.ErrUserNotFound
	}

	return user, nil
}
