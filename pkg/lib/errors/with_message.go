package errors

import (
	"fmt"
)

// copy from "github.com/pkg/errors"

type withMessage struct {
	cause error
	msg   string
}

func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

func (w *withMessage) Error() string {
	return w.msg + "\n->" + w.cause.Error()
}

func (w *withMessage) Unwrap() error {
	return w.cause
}
