package app

import "errors"

const (
	ExitSuccess  = 0
	ExitUsage    = 1
	ExitAuth     = 2
	ExitNetwork  = 3
	ExitTimeout  = 4
	ExitBackend  = 5
	ExitInternal = 10
)

type Error struct {
	Code int
	Err  error
}

func (e *Error) Error() string {
	if e == nil || e.Err == nil {
		return "application error"
	}

	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func WithExitCode(code int, err error) error {
	if err == nil {
		return nil
	}

	return &Error{
		Code: code,
		Err:  err,
	}
}

func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	var coded *Error
	if errors.As(err, &coded) && coded != nil && coded.Code != 0 {
		return coded.Code
	}

	return ExitInternal
}
