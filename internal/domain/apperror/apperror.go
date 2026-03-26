// Package apperror defines the standard domain-level error type.
// It is the single source of truth for structured errors across the
// usecase layer. HTTP status codes are included so the delivery layer
// (handler) can map them to HTTP responses without knowing about
// business logic — and without the domain knowing about HTTP.
package apperror

import (
	"fmt"
	"net/http"
)

// DomainError is the standard structured error for the domain and
// usecase layers. Delivery layers (HTTP handlers, queue jobs) check
// for this type and convert it to an appropriate response format.
type DomainError struct {
	// Code is a machine-readable error identifier (e.g. "ERR_NOT_FOUND").
	Code string
	// Message is a human-readable description returned to the client.
	Message string
	// Field identifies which request field caused the error (optional).
	Field string
	// Status is the suggested HTTP status code.
	Status int
}

// Error implements the built-in error interface.
func (e *DomainError) Error() string { return e.Message }

// GetStatus returns the HTTP status code associated with this error.
// Used by pkg/http response helpers via a local interface — so
// pkg/http never needs to import internal/domain/apperror.
func (e *DomainError) GetStatus() int { return e.Status }

// GetCode returns the machine-readable error code.
func (e *DomainError) GetCode() string { return e.Code }

// GetField returns the field name that caused the error (may be empty).
func (e *DomainError) GetField() string { return e.Field }

// ── Constructors ──────────────────────────────────────────────────

// New creates a fully-specified DomainError.
func New(code, field, message string, status int) *DomainError {
	return &DomainError{Code: code, Field: field, Message: message, Status: status}
}

// NotFound creates a 404 error.
func NotFound(format string, a ...any) *DomainError {
	return New("ERR_NOT_FOUND", "", fmt.Sprintf(format, a...), http.StatusNotFound)
}

// BadRequest creates a 400 error.
func BadRequest(format string, a ...any) *DomainError {
	return New("ERR_BAD_REQUEST", "", fmt.Sprintf(format, a...), http.StatusBadRequest)
}

// BadRequestField creates a 400 error tied to a specific request field.
func BadRequestField(field, format string, a ...any) *DomainError {
	return New("ERR_BAD_REQUEST", field, fmt.Sprintf(format, a...), http.StatusBadRequest)
}

// Conflict creates a 409 error.
func Conflict(code, field, message string) *DomainError {
	return New(code, field, message, http.StatusConflict)
}

// Forbidden creates a 403 error.
func Forbidden(format string, a ...any) *DomainError {
	return New("ERR_FORBIDDEN", "", fmt.Sprintf(format, a...), http.StatusForbidden)
}

// Unauthorized creates a 401 error.
func Unauthorized(format string, a ...any) *DomainError {
	return New("ERR_UNAUTHORIZED", "", fmt.Sprintf(format, a...), http.StatusUnauthorized)
}
