package totp

import (
	"context"

	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	"thomas.vn/apartment_service/internal/domain/repository"
	utotp "thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	pkgtotp "thomas.vn/apartment_service/pkg/totp"
)

type totpUsecase struct {
	logger   *xlogger.Logger
	userRepo repository.UserRepository
}

func NewTotpUsecase(
	logger *xlogger.Logger,
	userRepo repository.UserRepository,
) utotp.TotpUsecase {
	return &totpUsecase{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (u *totpUsecase) Generate(
	_ context.Context,
	user *xuser.User,
) (string, string, error) {
	if user.TotpSecret != nil {
		return "", "", xhttp.BadRequestErrorf("Totp already enabled")
	}

	result, err := pkgtotp.Generate(user.Email)
	if err != nil {
		return "", "", err
	}

	return result.Secret, result.QRCode, nil

}

func (u *totpUsecase) Save(
	ctx context.Context,
	user *xuser.User,
	secret,
	token string,
) error {
	if user.TotpSecret != nil {
		return xhttp.BadRequestErrorf("Totp already enabled")
	}

	if !pkgtotp.Verify(token, secret) {
		return xhttp.BadRequestErrorf("Invalid token")
	}

	return u.userRepo.UpdateTotpSecret(ctx, int64(user.ID), &secret)
}

func (u *totpUsecase) Verify(
	_ context.Context,
	user *xuser.User,
	token string,
) error {
	if user.TotpSecret == nil {
		return xhttp.BadRequestErrorf("Totp already enabled")
	}

	if !pkgtotp.Verify(token, *user.TotpSecret) {
		return xhttp.BadRequestErrorf("Invalid token")
	}

	return nil
}

func (u *totpUsecase) Disable(
	ctx context.Context,
	user *xuser.User,
	token string,
) error {
	if user.TotpSecret == nil {
		return xhttp.BadRequestErrorf("Totp already enabled")
	}

	if !pkgtotp.Verify(token, *user.TotpSecret) {
		return xhttp.BadRequestErrorf("Invalid token")
	}

	return u.userRepo.UpdateTotpSecret(ctx, int64(user.ID), nil)
}
