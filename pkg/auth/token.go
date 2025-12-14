package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	accessSecret  string
	accessExpire  time.Duration
	refreshSecret string
	refreshExpire time.Duration
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewToken(accessSecret string, accessExpire time.Duration, refreshSecret string, refreshExpire time.Duration) *Token {
	return &Token{
		accessSecret:  accessSecret,
		accessExpire:  accessExpire,
		refreshSecret: refreshSecret,
		refreshExpire: refreshExpire,
	}
}

func (s *Token) CreateTokens(userID uint) (accessToken, refreshToken string, err error) {
	accessToken, err = s.generateToken(userID, s.refreshSecret, s.refreshExpire)
	if err != nil {
		return "", "", err
	}
	return accessToken, accessToken, nil
}

func (s *Token) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return s.verifyToken(tokenStr, s.refreshSecret)
}

func (s *Token) VerifyAccessToken(tokenStr string) (*Claims, error) {
	return s.verifyToken(tokenStr, s.accessSecret)
}

func (s *Token) generateToken(userID uint, secret string, expire time.Duration) (string, error) {
	claims := &Claims{
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

func (s *Token) verifyToken(tokenStr, secret string) (*Claims, error) {
	claims := &Claims{}
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
