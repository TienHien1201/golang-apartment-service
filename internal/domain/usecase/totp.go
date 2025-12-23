package usecase

import (
	"context"

	xuser "thomas.vn/apartment_service/internal/domain/model/user"
)

type TotpUsecase interface {
	Generate(ctx context.Context, user *xuser.User) (string, string, error)
	Save(ctx context.Context, user *xuser.User, secret, token string) error
	Verify(ctx context.Context, user *xuser.User, token string) error
	Disable(ctx context.Context, user *xuser.User, token string) error
}
