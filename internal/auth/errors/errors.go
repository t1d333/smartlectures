package errors

import (
	"net/http"

	customerrors "github.com/t1d333/smartlectures/internal/errors"
)


var ErrUserNotFound = customerrors.New(
	http.StatusNotFound,
	"user not found",
	"user with this email not found",
)

var ErrWrongPassword = customerrors.New(
	http.StatusBadRequest,
	"wrong email or password",
	"wrong email or password",
)

var ErrBadToken = customerrors.New(
	http.StatusBadRequest,
	"bad auth cookie",
	"bad auth cookie",
)
