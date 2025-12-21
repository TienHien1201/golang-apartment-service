package auth

import (
	"context"
	"errors"
	"time"

	"thomas.vn/apartment_service/internal/domain/model"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	xfile "thomas.vn/apartment_service/pkg/file"
)

func (u *authUsecase) GoogleLogin(ctx context.Context, gUser *model.GoogleUser) (string, string, error) {
	if !gUser.EmailVerified {
		return "", "", errors.New("email not verified by google")
	}

	user, err := u.userRepo.GetUserByEmail(ctx, gUser.Email)
	if err != nil {
		return "", "", err
	}
	avatar := gUser.Avatar
	googleID := gUser.GoogleID

	if user == nil {
		user = &xuser.User{
			RoleID:    xfile.DefaultUserRoleID,
			Email:     gUser.Email,
			FullName:  gUser.FullName,
			Avatar:    avatar,
			GoogleID:  &googleID,
			IsActive:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		user, err = u.userRepo.CreateUser(ctx, user)
		if err != nil {
			return "", "", err
		}
	}

	accessToken, refreshToken, err := u.tokenSvc.CreateTokens(uint(user.ID))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
