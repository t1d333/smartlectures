package grpc

import (
	"context"
	"errors"

	"github.com/t1d333/smartlectures/internal/auth"
	autherrors "github.com/t1d333/smartlectures/internal/auth/errors"
	"github.com/t1d333/smartlectures/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Delivery struct {
	logger  logger.Logger
	service auth.Service
	auth.UnsafeAuthServiceServer
}

func (d *Delivery) CheckAuth(
	ctx context.Context,
	token *wrapperspb.StringValue,
) (*auth.AuthCheckResult, error) {
	session := string(token.GetValue())
	user, err := d.service.GetMe(ctx, session)

	if err != nil && errors.Is(err, autherrors.ErrSessionDoesNotExists) {
		return &auth.AuthCheckResult{
			UserId: 0,
			Status: auth.AuthStatus_Unauthorized,
		}, nil
	} else if err != nil {
		return &auth.AuthCheckResult{}, err
	}

	return &auth.AuthCheckResult{
		UserId: int32(user.UserId),
		Status: auth.AuthStatus_Authorized,
	}, nil
}

func NewDelivery(logger logger.Logger, service auth.Service) auth.AuthServiceServer {
	return &Delivery{
		logger:  logger,
		service: service,
	}
}
