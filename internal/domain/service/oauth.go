package service

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
)

type GoogleOAuthService interface {
	GetProfile(ctx context.Context, code string) (*model.GoogleUser, error)
}
