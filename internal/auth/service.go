package auth

import (
	"context"

	authmodels "github.com/t1d333/smartlectures/internal/auth/models"
	"github.com/t1d333/smartlectures/internal/models"
)


type Service interface {
	Login(ctx context.Context, data authmodels.LoginRequest) (string, error)
	Logout(ctx context.Context, session string) error
	Register(ctx context.Context, data authmodels.RegisterRequest) (models.User, error)
	Refresh(ctx context.Context) error
}
