package middl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	customErrors "github.com/t1d333/smartlectures/internal/errors"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	var customErr *customErrors.Error
	for _, err := range c.Errors {
		if errors.As(err, &customErr) {
			break
		}
	}

	if customErr != nil {
		c.JSON(customErr.HttpCode(), gin.H{"msg": customErr.ResponseMsg()})
	} else if c.Errors.Last() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "internal server error"})
	}
}
