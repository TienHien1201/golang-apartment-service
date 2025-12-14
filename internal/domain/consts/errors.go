package consts

import (
	"fmt"
	"net/http"

	xhttp "thomas.vn/apartment_service/pkg/http"
)

const (
	ErrUserEmailAlreadyExists = "ERR_USER_EMAIL_ALREADY_EXISTS"
	ErrJobCodeAlreadyExists   = "ERR_JOB_CODE_ALREADY_EXISTS"
	ErrJobCodeNotFound        = "ERR_JOB_CODE_NOT_FOUND"
	ErrJobIDNotFound          = "ERR_JOB_ID_NOT_FOUND"
)

func EmailAlreadyExistsError(email string) *xhttp.AppError {
	return xhttp.NewAppError(ErrUserEmailAlreadyExists, "email", fmt.Sprintf("user with email %s already exists", email), http.StatusConflict)
}

func JobCodeAlreadyExistsError(code string) *xhttp.AppError {
	return xhttp.NewAppError(ErrJobCodeAlreadyExists, "code", fmt.Sprintf("job with code %s already exists", code), http.StatusConflict)
}

func JobCodeNotFoundError(code string) *xhttp.AppError {
	return xhttp.NewAppError(ErrJobCodeNotFound, "code", fmt.Sprintf("job with code %s not found", code), http.StatusNotFound)
}

func JobIDNotFoundError(id uint64) *xhttp.AppError {
	return xhttp.NewAppError(ErrJobIDNotFound, "id", fmt.Sprintf("job with id %d not found", id), http.StatusNotFound)
}
