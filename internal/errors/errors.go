package errors

type Error struct {
	code        int
	msg         string
	responseMsg string
}

func (e *Error) HttpCode() int {
	return e.code
}

func (e *Error) ResponseMsg() string {
	return e.responseMsg
}

func (e *Error) Error() string {
	return e.responseMsg
}

func New(code int, responseMsg, msg string) error {
	return &Error{
		code:        code,
		responseMsg: responseMsg,
		msg:         msg,
	}
}
