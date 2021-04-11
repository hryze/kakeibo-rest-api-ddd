package apperrors

import (
	"net/http"
)

type errorString struct {
	Message string `json:"message"`
}

func (e *errorString) Error() string {
	return e.Message
}

func NewErrorString(message string) error {
	return &errorString{
		Message: message,
	}
}

func (e *appError) BadRequest(infoMessage error) *appError {
	e.status = http.StatusBadRequest
	e.infoMessage = infoMessage
	e.LevelInfo()

	return e
}

func (e *appError) InternalServerError(infoMessage error) *appError {
	e.status = http.StatusInternalServerError
	e.infoMessage = infoMessage
	e.LevelError()

	return e
}

func (e *appError) Wrap(err error, msg ...string) *appError {
	if err == nil {
		return e
	}

	var m string
	if len(msg) != 0 {
		m = msg[0]
	}

	er := create(m)
	er.next = err

	e.next = er

	return e
}

func (e *appError) Status() int {
	if e.status != 0 {
		return e.status
	}

	next := AsAppError(e.next)
	if next != nil {
		return next.Status()
	}

	return http.StatusInternalServerError
}

func (e *appError) InfoMessage() error {
	if e.infoMessage != nil {
		return e.infoMessage
	}

	next := AsAppError(e.next)
	if next != nil {
		return next.InfoMessage()
	}

	return &errorString{"Internal Server Error"}
}
