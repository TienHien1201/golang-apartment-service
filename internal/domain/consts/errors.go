package consts

import (
	"fmt"

	"thomas.vn/apartment_service/internal/domain/apperror"
)

const (
	ErrUserEmailAlreadyExists = "ERR_USER_EMAIL_ALREADY_EXISTS"
	ErrJobCodeAlreadyExists   = "ERR_JOB_CODE_ALREADY_EXISTS"
	ErrJobCodeNotFound        = "ERR_JOB_CODE_NOT_FOUND"
	ErrJobIDNotFound          = "ERR_JOB_ID_NOT_FOUND"
)

func EmailAlreadyExistsError(email string) *apperror.DomainError {
	return apperror.Conflict(ErrUserEmailAlreadyExists, "email", fmt.Sprintf("user with email %s already exists", email))
}

func JobCodeAlreadyExistsError(code string) *apperror.DomainError {
	return apperror.Conflict(ErrJobCodeAlreadyExists, "code", fmt.Sprintf("job with code %s already exists", code))
}

func JobCodeNotFoundError(code string) *apperror.DomainError {
	return apperror.New(ErrJobCodeNotFound, "code", fmt.Sprintf("job with code %s not found", code), 404)
}

func JobIDNotFoundError(id uint64) *apperror.DomainError {
	return apperror.New(ErrJobIDNotFound, "id", fmt.Sprintf("job with id %d not found", id), 404)
}
