package usecase

import (
	"context"

	"thomas.vn/hr_recruitment/internal/domain/model"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	Register(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken string, err error)
	Logout(ctx context.Context) error
}
