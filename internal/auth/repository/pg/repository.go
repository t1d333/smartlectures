package pg

import (
	"context"
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
	res := r.pool.QueryRow(ctx, cmd, user.Username, user.Name, user.Email, user.Surname, user.Password)

	if err := res.Scan(&user.UserId); err != nil {
		return user, fmt.Errorf("PostgresRepository.RegisterUser(user: %v): %w", user, err)
	}

	return user, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
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
