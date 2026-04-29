package domain

import "fmt"

type ErrCode string

const (
	CodeNotFound     ErrCode = "NOT_FOUND"
	CodeUnauthorized ErrCode = "UNAUTHORIZED"
	CodeInvalid      ErrCode = "INVALID"
)

type AppError struct {
	Code ErrCode
	Op   string
	Err  error
}

func (e *AppError) Error() string {
	if e.Op == "" {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Op, e.Err)
}

func (e *AppError) Unwrap() error { return e.Err }

func NotFound(op string, err error) error {
	return &AppError{Code: CodeNotFound, Op: op, Err: err}
}

func Unauthorized(op string, err error) error {
	return &AppError{Code: CodeUnauthorized, Op: op, Err: err}
}

func Invalid(op string, err error) error {
	return &AppError{Code: CodeInvalid, Op: op, Err: err}
}
