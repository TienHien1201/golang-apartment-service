package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
	xauth "thomas.vn/apartment_service/internal/domain/model/auth"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
)

type AuthUsecase interface {
	Login(ctx context.Context, email string, password string, totpToken *string) (*xauth.AuthLoginResult, error)
	GoogleLogin(ctx context.Context, gUser *model.GoogleUser) (string, string, error)
	Register(ctx context.Context, req *xuser.CreateUserRequest) (*xuser.User, error)
	RefreshToken(ctx context.Context, accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	Logout(ctx context.Context) error
	GetInfo(_ context.Context, user *xuser.User) (*xauth.AuthInfoResult, error)
}
