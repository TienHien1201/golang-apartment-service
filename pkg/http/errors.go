package xhttp

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, field, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Field:   field,
	}
}

func NotFoundErrorf(format string, a ...any) *AppError {
	return NewAppError("ERR_NOT_FOUND", "", fmt.Sprintf(format, a...), http.StatusNotFound)
}

func BadRequestErrorf(format string, a ...any) *AppError {
	return NewAppError("ERR_BAD_REQUEST", "", fmt.Sprintf(format, a...), http.StatusBadRequest)
}

func ForbiddenErrorf(format string, a ...any) *AppError {
	return NewAppError("ERR_FORBIDDEN", "", fmt.Sprintf(format, a...), http.StatusForbidden)
}

func UnauthorizedErrof(format string, a ...any) *AppError {
	return NewAppError("ERR_UNAUTHORIZED", "", fmt.Sprintf(format, a...), http.StatusUnauthorized)
}
