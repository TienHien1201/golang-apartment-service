package user

import (
	"context"
	"time"

	"thomas.vn/apartment_service/internal/domain/consts"
	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/service"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type userUsecase struct {
	logger   *xlogger.Logger
	userRepo repository.UserRepository
	cacheSvc service.CacheService
}

func NewUserUsecase(logger *xlogger.Logger, userRepo repository.UserRepository, cacheSvc service.CacheService) usecase.UserUsecase {
	return &userUsecase{
		logger:   logger,
		userRepo: userRepo,
		cacheSvc: cacheSvc,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Failed to check existing user", xlogger.Error(err))
		return nil, err
	}
	if existingUser != nil {
		return nil, consts.EmailAlreadyExistsError(req.Email)
	}

	user := &model.User{
		Email:    req.Email,
		Password: req.Password, // Note: Password should be hashed before storing
		FullName: req.FullName,
		RoleID:   req.RoleID,
		IsActive: consts.UserStatusActive,
	}

	createdUser, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		u.logger.Error("Failed to create user", xlogger.Error(err))
		return nil, err
	}

	return createdUser, nil
}

func (u *userUsecase) GetUser(ctx context.Context, id uint) (*model.User, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		u.logger.Error("Failed to get user", xlogger.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, xhttp.NotFoundErrorf("User with ID %d not found", id)
	}

	return user, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := u.GetUser(ctx, uint(req.ID))
	if err != nil {
		u.logger.Error("Failed to get user", xlogger.Error(err))
		return nil, err
	}

	if req.Password != "" {
		// Note: Password should be hashed before storing
		user.Password = req.Password
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.RoleID != 0 {
		user.RoleID = req.RoleID
	}
	if req.IsActive != 0 {
		user.IsActive = req.IsActive
	}

	updatedUser, err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		u.logger.Error("Failed to update user", xlogger.Error(err))
		return nil, err
	}

	return updatedUser, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	user, err := u.GetUser(ctx, id)
	if err != nil {
		u.logger.Error("Failed to get user", xlogger.Error(err))
		return err
	}

	if err := u.userRepo.DeleteUser(ctx, uint(user.ID)); err != nil {
		u.logger.Error("Failed to delete user", xlogger.Error(err))
		return err
	}

	return nil
}

func (u *userUsecase) ListUsers(ctx context.Context, req *model.ListUserRequest) ([]*model.User, int64, error) {
	users, total, err := u.userRepo.ListUsers(ctx, req)
	if err != nil {
		u.logger.Error("Failed to list user", xlogger.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}

func (u *userUsecase) DeleteUsersCreatedBefore(_ context.Context, days time.Time) error {
	u.logger.Info("Deleting users created before", xlogger.Object("days", days))

	return nil
}
