package errors

import "net/http"

var ErrNoteNotFound = &Error{
	code:        http.StatusNotFound,
	msg:         "note not found",
	responseMsg: "note with this id not found",
}

var InternalError = &Error{
	code:        http.StatusInternalServerError,
	msg:         "internal server error",
	responseMsg: "internal server error",
}

var BadRequestError = &Error {
	code: http.StatusBadRequest,
	msg: "bad request",
	responseMsg: "bad request",
}
