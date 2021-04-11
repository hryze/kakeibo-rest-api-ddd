package apperrors

import (
	"fmt"

	"golang.org/x/xerrors"
)

type appError struct {
	// Error for log output
	next    error
	message string
	frame   xerrors.Frame

	// Level for log output
	level level

	// Error for response to client
	status      int
	infoMessage error
}

func (e *appError) Error() string {
	next := AsAppError(e.next)
	if next != nil {
		return next.Error()
	}

	if e.next == nil {
		if e.message != "" {
			return e.message
		}

		return "no message"
	}

	return e.next.Error()
}

func (e *appError) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

func (e *appError) FormatError(p xerrors.Printer) error {
	var message string

	if e.level != "" {
		message += fmt.Sprintf("[%s] ", e.level)
	}

	if e.message != "" {
		message += fmt.Sprintf("%s", e.message)
	}

	p.Print(message)
	e.frame.Format(p)
	return e.next
}

func create(msg string) *appError {
	var e appError
	e.message = msg
	e.frame = xerrors.Caller(2)

	return &e
}

func New(msg string) *appError {
	return create(msg)
}

func Wrap(err error, msg ...string) *appError {
	if err == nil {
		return nil
	}

	var m string
	if len(msg) != 0 {
		m = msg[0]
	}

	e := create(m)
	e.next = err

	return e
}

func AsAppError(err error) *appError {
	if err == nil {
		return nil
	}

	var e *appError
	if xerrors.As(err, &e) {
		return e
	}

	return nil
}
