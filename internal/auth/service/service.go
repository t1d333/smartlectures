package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/t1d333/smartlectures/internal/auth"
	autherrors "github.com/t1d333/smartlectures/internal/auth/errors"
	authmodels "github.com/t1d333/smartlectures/internal/auth/models"
	"github.com/t1d333/smartlectures/internal/auth/repository"
	"github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

const Expire = 24 * time.Hour

type Service struct {
	logger     logger.Logger
	repository repository.Repository
}

func New(logger logger.Logger, repo repository.Repository) auth.Service {
	return &Service{
		logger:     logger,
		repository: repo,
	}
}

func (s *Service) Login(ctx context.Context, data authmodels.LoginRequest) (string, error) {
	token := uuid.NewString()

	user, err := s.repository.GetUser(ctx, data.Email)
	if err != nil {
		return "", fmt.Errorf("authService.Login(data: %v): %w", data, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", autherrors.ErrWrongPassword
		}
	}


	err = s.repository.AddNewSession(ctx, authmodels.SessionInfo{
		UserId:    user.UserId,
		Token:     token,
		IPAddress: "",
		Expire:    time.Now().Add(Expire),
	})
	
	token = fmt.Sprintf("%s$%d", token, user.UserId)

	return token, err
}

func (s *Service) Logout(ctx context.Context, session string) error {
	tmp := strings.Split(session, "$")

	if len(tmp) != 2 {
		return autherrors.ErrBadToken
	}

	token := tmp[0]
	userId, err := strconv.Atoi(tmp[1])
	if err != nil {
		return autherrors.ErrBadToken
	}

	return s.repository.DeleteSession(ctx, userId, token)
}

func (*Service) Refresh(ctx context.Context) error {
	panic("unimplemented")
}

func (*Service) Register(
	ctx context.Context,
	data authmodels.RegisterRequest,
) (models.User, error) {
	panic("unimplemented")
}
