package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"thomas.vn/apartment_service/internal/domain/consts"
	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
	"thomas.vn/apartment_service/pkg/auth"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xqueue "thomas.vn/apartment_service/pkg/queue"
)

type authUsecase struct {
	logger       *xlogger.Logger
	userRepo     repository.UserRepository
	tokenSvc     *auth.Token
	queueService xqueue.QueueService
}

func NewAuthUsecase(logger *xlogger.Logger, userRepo repository.UserRepository, token *auth.Token, queueService xqueue.QueueService) usecase.AuthUsecase {
	return &authUsecase{
		logger:       logger,
		userRepo:     userRepo,
		tokenSvc:     token,
		queueService: queueService,
	}
}
func (u *authUsecase) Register(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {

	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		RoleID:    req.RoleID,
		IsActive:  1,
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

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		u.logger.Error("Login Failed - user not found ", xlogger.Error(err))
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		u.logger.Error("Login Failed - wrong password", xlogger.Error(err))
		return "", "", err
	}

	accessToken, refreshToken, err := u.tokenSvc.CreateTokens(uint(user.ID))
	if err != nil {
		u.logger.Error("Login Failed - token creation failed", xlogger.Error(err))
		return "", "", err
	}
	u.logger.Info("User login successfully", xlogger.String("email", email), xlogger.String("access_token", accessToken))
	_ = u.queueService.PublishMessage(
		ctx,
		consts.MailJobType,
		&model.MailPayload{
			Type:     consts.QueueMailLogin,
			Email:    user.Email,
			FullName: user.FullName,
		},
	)

	return accessToken, refreshToken, nil
}
func (u *authUsecase) RefreshToken(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) (string, string, error) {

	refreshClaims, err := u.tokenSvc.VerifyRefreshToken(refreshToken)
	if err != nil {
		u.logger.Error("Invalid refresh token", xlogger.Error(err))
		return "", "", err
	}

	accessClaims, err := u.tokenSvc.VerifyAccessToken(accessToken)
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
		return "", "", fmt.Errorf("user does not exist")
	}

	newAccessToken, newRefreshToken, err := u.tokenSvc.CreateTokens(uint(user.ID))
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

func (u *authUsecase) GetInfo(_ context.Context, user *model.User) (*model.User, error) {
	return user, nil
}
