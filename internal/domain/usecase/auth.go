package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	Register(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken string, err error)
	Logout(ctx context.Context) error
	GetInfo(ctx context.Context, user *model.User) (*model.User, error)
}
