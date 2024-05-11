package dragonfly

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
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
