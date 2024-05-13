package dragonfly

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	autherrors "github.com/t1d333/smartlectures/internal/auth/errors"
	"github.com/t1d333/smartlectures/internal/auth/models"
	authmodels "github.com/t1d333/smartlectures/internal/auth/models"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type DragonflyRepository struct {
	logger logger.Logger
	client *redis.Client
}

func New(logger logger.Logger, client *redis.Client) *DragonflyRepository {
	return &DragonflyRepository{
		logger: logger,
		client: client,
	}
}

func (r *DragonflyRepository) DeleteSession(ctx context.Context, userId int, token string) error {
	if res := r.client.HDel(ctx, strconv.Itoa(userId), token); res.Err() != nil {
		return fmt.Errorf(
			"DragonflyRepository.DeleteSession(userId: %d, token: %s): session: %w",
			userId,
			token,
			res.Err(),
		)
	}

	return nil
}

func (r *DragonflyRepository) GetSession(
	ctx context.Context,
	userId int,
	token string,
) (authmodels.SessionInfo, error) {
	cmd := r.client.HExists(ctx, strconv.Itoa(userId), token)

	exists, err := cmd.Result()
	if err != nil {
		return authmodels.SessionInfo{}, fmt.Errorf(
			"DragonflyRepository.GetSession(userId: %d, token: %s): session: %w",
			userId,
			token,
			err,
		)
	}

	if !exists {
		return models.SessionInfo{}, autherrors.ErrSessionDoesNotExists
	}

	if res := r.client.HGet(ctx, strconv.Itoa(userId), token); res.Err() != nil {
		return models.SessionInfo{}, fmt.Errorf(
			"DragonflyRepository.GetSession(userId: %d, token: %s): session: %w",
			userId,
			token,
			res.Err(),
		)
	} else {
		info := models.SessionInfo{}
		if err := res.Scan(&info); err != nil {
			return models.SessionInfo{}, fmt.Errorf(
				"DragonflyRepository.GetSession(userId: %d, token: %s): session: %w",
				userId,
				token,
				err,
			)
		}

		if info.Expire.Compare(time.Now()) <= 0 {
			err := r.DeleteSession(ctx, userId, token)
			if err != nil {
				return models.SessionInfo{}, fmt.Errorf(
					"DragonflyRepository.GetSession(userId: %d, token: %s): session: %w",
					userId,
					token,
					err,
				)
			}

			return models.SessionInfo{}, autherrors.ErrSessionDoesNotExists

		}

		return info, nil
	}
}

func (r *DragonflyRepository) AddSession(
	ctx context.Context,
	sessionInfo authmodels.SessionInfo,
) error {
	if res := r.client.HSet(ctx, strconv.Itoa(sessionInfo.UserId), sessionInfo.Token, sessionInfo); res.Err() != nil {
		return fmt.Errorf(
			"DragonflyRepository.AddSession(session: %v): %w",
			sessionInfo,
			res.Err(),
		)
	}

	return nil
}
