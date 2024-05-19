package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1d333/smartlectures/internal/auth"
	autherrors "github.com/t1d333/smartlectures/internal/auth/errors"
	customerrors "github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func NewAuthHandler(client auth.AuthServiceClient, logger logger.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := c.Cookie("session")
		if err != nil {
			_ = c.AbortWithError(
				customerrors.BadRequestError.HttpCode(),
				customerrors.BadRequestError,
			)
			return
		}

		res, err := client.CheckAuth(c.Request.Context(), &wrapperspb.StringValue{
			Value: token,
		})
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, customerrors.InternalError)
			return
		}

		if res.Status == auth.AuthStatus_Unauthorized {
			_ = c.AbortWithError(http.StatusUnauthorized, autherrors.ErrUserUnauthorized)
		}

		c.Set("userId", int(res.GetUserId()))

		c.Next()
	}
}
