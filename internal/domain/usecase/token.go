package usecase

import (
	"time"

	"thomas.vn/apartment_service/internal/domain/model/token"
)

type TokenUsecase interface {
	CreateTokens(userID uint) (accessToken, refreshToken string, err error)
	GenerateToken(userID uint, secret string, expire time.Duration) (string, error)
	VerifyToken(tokenStr, secret string) (*token.Claims, error)
	VerifyRefreshToken(tokenStr string) (*token.Claims, error)
	VerifyAccessToken(tokenStr string) (*token.Claims, error)
}
