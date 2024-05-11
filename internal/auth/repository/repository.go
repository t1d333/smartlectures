package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	authmodels "github.com/t1d333/smartlectures/internal/auth/models"
	"github.com/t1d333/smartlectures/internal/auth/repository/dragonfly"
	"github.com/t1d333/smartlectures/internal/auth/repository/pg"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Repository interface {
	// CheckUserPassword(ctx context.Context, email string, password []byte) error
	RegisterNewUser(ctx context.Context, user models.User) (models.User, error)
	DeleteSession(ctx context.Context, userId int, session string) error
	AddNewSession(ctx context.Context, session authmodels.SessionInfo) error
	CheckExistingUsers(ctx context.Context, user models.User) ([]models.User, error)
	GetUser(ctx context.Context, email string) (models.User, error)
}

type RepositoryImpl struct {
	dragonflyRepo *dragonfly.DragonflyRepository
	postgresRepo  *pg.PostgresRepository
	logger        logger.Logger
}

func (r *RepositoryImpl) GetUser(ctx context.Context, email string) (models.User, error) {
	return r.postgresRepo.GetUser(ctx, email)
}

func New(logger logger.Logger, rdc *redis.Client, pgc *pgxpool.Pool) Repository {
	return &RepositoryImpl{
		logger:        logger,
		dragonflyRepo: dragonfly.New(logger, rdc),
		postgresRepo:  pg.New(logger, pgc),
	}
}

func (r *RepositoryImpl) AddNewSession(ctx context.Context, session authmodels.SessionInfo) error {
	err := r.dragonflyRepo.AddSession(ctx, session)
	if err != nil {
		return fmt.Errorf("authRepository.AddNewSession(session: %v): %w", session, err)
	}

	return nil
}

func (r *RepositoryImpl) CheckExistingUsers(
	tx context.Context,
	user models.User,
) ([]models.User, error) {
	panic("unimplemented")
}

func (r *RepositoryImpl) DeleteSession(ctx context.Context, userId int, session string) error {
	err := r.dragonflyRepo.DeleteSession(ctx, userId, session)
	if err != nil {
		return fmt.Errorf(
			"authRepository.AddDeleteSession(userId: %d, token: %s): %w",
			userId,
			session,
			err,
		)
	}

	return nil
}

func (r *RepositoryImpl) RegisterNewUser(
	ctx context.Context,
	user models.User,
) (models.User, error) {
	panic("unimplemented")
}
