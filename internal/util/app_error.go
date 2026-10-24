package util

import "net/http"

const (
	SQL_UNIQUE_VIOLATION string = "23505"
)

var (
	ErrInternalServerErr *AppError = New("internal server error", http.StatusInternalServerError, nil)
)

type AppError struct {
	Message    string
	HttpStatus int
	Err        error
}

func New(message string, status int, err error) *AppError {
	return &AppError{
		Message:    message,
		HttpStatus: status,
		Err:        err,
	}
}

func InternalServerError(err error) *AppError {
	return &AppError{
		Err:        err,
		Message:    "internal server error",
		HttpStatus: http.StatusInternalServerError,
	}
}

func (ae *AppError) Error() string {
	return ae.Message
}
