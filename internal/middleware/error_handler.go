package middl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	customErrors "github.com/t1d333/smartlectures/internal/errors"
	"github.com/t1d333/smartlectures/pkg/logger"
)

func NewErrorHandler(logger logger.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var customErr *customErrors.Error

		err := c.Errors[0]

		logger.Error(err)

		if errors.As(err, &customErr) {
			c.JSON(customErr.HttpCode(), gin.H{"msg": customErr.ResponseMsg()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "internal server error"})
		}
	}
}
