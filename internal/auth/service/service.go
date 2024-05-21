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

	storage "github.com/t1d333/smartlectures/internal/storage"
)

const Expire = 24 * time.Hour

type Service struct {
	logger        logger.Logger
	repository    repository.Repository
	storageClient storage.StorageClient
}

func New(
	logger logger.Logger,
	repo repository.Repository,
	storageClient storage.StorageClient,
) auth.Service {
	return &Service{
		logger:        logger,
		repository:    repo,
		storageClient: storageClient,
	}
}

func (s *Service) GetMe(ctx context.Context, session string) (models.User, error) {
	tmp := strings.Split(session, "$")

	if len(tmp) != 2 {
		return models.User{}, autherrors.ErrBadToken
	}

	token := tmp[0]
	userId, err := strconv.Atoi(tmp[1])
	if err != nil {
		return models.User{}, autherrors.ErrBadToken
	}

	if _, err = s.repository.GetSession(ctx, userId, token); err != nil {
		return models.User{}, fmt.Errorf("AuthService.GetMe(session: %s): %w", session, err)
	}

	user, err := s.repository.GetUserById(ctx, userId)
	if err != nil {
		return models.User{}, fmt.Errorf("AuthService.GetMe(session: %s): %w", session, err)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, data authmodels.LoginRequest) (string, error) {
	token := uuid.NewString()
	clientIp := ctx.Value("client_ip").(string)

	user, err := s.repository.GetUserByEmail(ctx, data.Email)
	if err != nil && errors.Is(err, autherrors.ErrUserNotFound) {
		return "", autherrors.ErrWrongPassword
	} else if err != nil {
		return "", fmt.Errorf("AuthService.Login(data: %v): %w", data, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", autherrors.ErrWrongPassword
		}
	}

	err = s.repository.AddNewSession(ctx, authmodels.SessionInfo{
		UserId:    user.UserId,
		Token:     token,
		IPAddress: clientIp,
		Expire:    time.Now().Add(Expire),
	})

	token = fmt.Sprintf("%s$%d", token, user.UserId)

	return token, err
}

func (s *Service) IsAuthorized(ctx context.Context, session string) error {
	tmp := strings.Split(session, "$")

	if len(tmp) != 2 {
		return autherrors.ErrBadToken
	}

	token := tmp[0]
	userId, err := strconv.Atoi(tmp[1])
	if err != nil {
		return autherrors.ErrBadToken
	}

	if _, err := s.repository.GetSession(ctx, userId, token); err != nil {
		if errors.Is(err, autherrors.ErrSessionDoesNotExists) {
			return autherrors.ErrUserUnauthorized
		}
		return fmt.Errorf("AuthService.IsAuthorized(session: %s): %w", session, err)
	}

	return nil
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

func (s *Service) Register(
	ctx context.Context,
	data authmodels.RegisterRequest,
) (models.User, error) {
	_, err := s.repository.GetUserByEmail(ctx, data.Email)

	if err != nil && !errors.Is(err, autherrors.ErrUserNotFound) {
		return models.User{}, fmt.Errorf("AuthService.Register(data: %v): %w", data, err)
	} else if errors.Is(err, autherrors.ErrUserNotFound) {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, fmt.Errorf("AuthService.Register(data: %v): %w", data, err)
		}

		user, err := s.repository.RegisterNewUser(ctx, models.User{
			Username: data.Username,
			Email:    data.Email,
			Surname:  data.Surname,
			Name:     data.Name,
			Password: hashedPassword,
		})
		if err != nil {
			return models.User{}, fmt.Errorf("AuthService(data: %v): %w", data, err)
		}

		onboardNote, err := s.repository.GetOnboardNote(ctx, user.UserId)
		if err != nil {
			s.logger.Error("failed to get onboard note in AuthService.Register()", err)
		}

		_, err = s.storageClient.CreateNote(ctx, &storage.Note{
			Id:        int32(onboardNote.NoteId),
			Name:      onboardNote.Name,
			Body:      onboardNote.Body,
			ParentDir: int32(onboardNote.ParentDir),
			UserId:    int32(user.UserId),
		})

		if err != nil {
			s.logger.Error("failed to get onboard note in AuthService.Register()", err)
		}
		
		return user, nil
	}

	return models.User{}, autherrors.ErrUserAlreadyExists
}
