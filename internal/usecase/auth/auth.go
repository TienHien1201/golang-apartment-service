package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"thomas.vn/apartment_service/internal/domain/consts"
	"thomas.vn/apartment_service/internal/domain/model"
	xauth "thomas.vn/apartment_service/internal/domain/model/auth"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xqueue "thomas.vn/apartment_service/pkg/queue"
	pkgtotp "thomas.vn/apartment_service/pkg/totp"
)

type authUsecase struct {
	logger       *xlogger.Logger
	userRepo     repository.UserRepository
	tokenUc      usecase.TokenUsecase
	queueService xqueue.QueueService
}

func NewAuthUsecase(logger *xlogger.Logger, userRepo repository.UserRepository, tokenUc usecase.TokenUsecase, queueService xqueue.QueueService) usecase.AuthUsecase {
	return &authUsecase{
		logger:       logger,
		userRepo:     userRepo,
		tokenUc:      tokenUc,
		queueService: queueService,
	}
}
func (u *authUsecase) Register(ctx context.Context, req *xuser.CreateUserRequest) (*xuser.User, error) {

	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, xhttp.BadRequestErrorf("Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &xuser.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		IsActive:  consts.UserStatusActive,
		RoleID:    consts.DefaultUserRoleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = u.queueService.PublishMessage(
		ctx,
		consts.MailJobType,
		&model.MailPayload{
			Type:     consts.QueueMailRegister,
			Email:    req.Email,
			FullName: req.FullName,
		},
	)

	return u.userRepo.CreateUser(ctx, newUser)
}

func (u *authUsecase) Login(ctx context.Context, email string, password string, totpToken *string) (*xauth.AuthLoginResult, error) {

	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, xhttp.BadRequestErrorf("Email not exists")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, xhttp.BadRequestErrorf("Wrong password")
	}

	accessToken, refreshToken, err := u.tokenUc.CreateTokens(uint(user.ID))
	if err != nil {
		return nil, err
	}

	if user.TotpSecret != nil {
		if totpToken == nil {
			return &xauth.AuthLoginResult{
				IsTotp:       true,
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}, nil
		}

		if !pkgtotp.Verify(*totpToken, *user.TotpSecret) {
			return nil, xhttp.BadRequestErrorf("Invalid totp token")
		}
	}
	_ = u.queueService.PublishMessage(
		ctx,
		consts.MailJobType,
		&model.MailPayload{
			Type:     consts.QueueMailLogin,
			Email:    user.Email,
			FullName: user.FullName,
		},
	)

	return &xauth.AuthLoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) RefreshToken(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) (string, string, error) {

	refreshClaims, err := u.tokenUc.VerifyRefreshToken(refreshToken)
	if err != nil {
		u.logger.Error("Invalid refresh token", xlogger.Error(err))
		return "", "", err
	}

	accessClaims, err := u.tokenUc.VerifyAccessToken(accessToken)
	if err != nil {
		u.logger.Error("Invalid access token", xlogger.Error(err))
		return "", "", err
	}

	if accessClaims.UserID != refreshClaims.UserID {
		err = fmt.Errorf("token user mismatch")
		u.logger.Error("Token invalid", xlogger.Error(err))
		return "", "", err
	}

	user, err := u.userRepo.GetUserByID(ctx, refreshClaims.UserID)
	if err != nil {
		u.logger.Error("User not found", xlogger.Error(err))
		return "", "", err
	}
	if user == nil {
		return "", "", xhttp.BadRequestErrorf("user does not exist")
	}

	newAccessToken, newRefreshToken, err := u.tokenUc.CreateTokens(uint(user.ID))
	if err != nil {
		u.logger.Error("Failed to generate new access token", xlogger.Error(err))
		return "", "", err
	}

	u.logger.Info(
		"Token refreshed successfully",
		xlogger.Uint("user_id", uint(user.ID)),
	)

	return newAccessToken, newRefreshToken, nil
}

func (u *authUsecase) Logout(_ context.Context) error {
	u.logger.Info("User logged out (stateless)")
	return nil
}

func (u *authUsecase) GetInfo(_ context.Context, user *xuser.User) (*xauth.AuthInfoResult, error) {

	return &xauth.AuthInfoResult{
		ID:       int64(user.ID),
		Email:    user.Email,
		IsTotp:   user.TotpSecret != nil,
		Avatar:   user.Avatar,
		FullName: user.FullName,
	}, nil
}
