package usecase

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"thomas.vn/apartment_service/internal/config"
	"thomas.vn/apartment_service/internal/domain/model/token"
)

type Token struct {
	accessSecret  string
	accessExpire  time.Duration
	refreshSecret string
	refreshExpire time.Duration
}

func NewToken(cfg config.TokenConfig) *Token {
	return &Token{
		accessSecret:  cfg.AccessSecret,
		accessExpire:  cfg.AccessExpire,
		refreshSecret: cfg.RefreshSecret,
		refreshExpire: cfg.RefreshExpire,
	}
}

func (s *Token) CreateTokens(userID uint) (accessToken, refreshToken string, err error) {
	accessToken, err = s.GenerateToken(userID, s.accessSecret, s.accessExpire)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = s.GenerateToken(userID, s.refreshSecret, s.refreshExpire)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *Token) VerifyRefreshToken(tokenStr string) (*token.Claims, error) {
	return s.VerifyToken(tokenStr, s.refreshSecret)
}

func (s *Token) VerifyAccessToken(tokenStr string) (*token.Claims, error) {
	return s.VerifyToken(tokenStr, s.accessSecret)
}

func (s *Token) GenerateToken(userID uint, secret string, expire time.Duration) (string, error) {
	claims := &token.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *Token) VerifyToken(tokenStr, secret string) (*token.Claims, error) {
	claims := &token.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidId
	}
	return claims, nil
}
