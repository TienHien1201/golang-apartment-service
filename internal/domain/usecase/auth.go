package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	GoogleLogin(ctx context.Context, gUser *model.GoogleUser) (string, string, error)
	Register(ctx context.Context, req *xuser.CreateUserRequest) (*xuser.User, error)
	RefreshToken(ctx context.Context, accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	Logout(ctx context.Context) error
	GetInfo(ctx context.Context, user *xuser.User) (*xuser.User, error)
}
