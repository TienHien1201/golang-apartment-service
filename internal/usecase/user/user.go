package user

import (
	"context"
	"path/filepath"
	"time"

	"thomas.vn/apartment_service/internal/domain/consts"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/service"
	user2 "thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xqueue "thomas.vn/apartment_service/pkg/queue"
)

type userUsecase struct {
	logger      *xlogger.Logger
	userRepo    repository.UserRepository
	cacheSvc    service.CacheService
	fileService service.FileService
	queue       xqueue.QueueService
}

func NewUserUsecase(logger *xlogger.Logger, userRepo repository.UserRepository, cacheSvc service.CacheService, fileService service.FileService, queue xqueue.QueueService) user2.UserUsecase {
	return &userUsecase{
		logger:      logger,
		userRepo:    userRepo,
		cacheSvc:    cacheSvc,
		fileService: fileService,
		queue:       queue,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, req *xuser.CreateUserRequest) (*xuser.User, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		u.logger.Error("Failed to check existing user", xlogger.Error(err))
		return nil, err
	}
	if existingUser != nil {
		return nil, consts.EmailAlreadyExistsError(req.Email)
	}

	user := &xuser.User{
		Email:    req.Email,
		Password: req.Password, // Note: Password should be hashed before storing
		FullName: req.FullName,
		RoleID:   consts.DefaultUserRoleID,
		IsActive: consts.UserStatusActive,
	}

	createdUser, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		u.logger.Error("Failed to create user", xlogger.Error(err))
		return nil, err
	}

	return createdUser, nil
}

func (u *userUsecase) GetUser(ctx context.Context, id uint) (*xuser.User, error) {
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

func (u *userUsecase) UpdateUser(ctx context.Context, req *xuser.UpdateUserRequest) (*xuser.User, error) {
	user, err := u.GetUser(ctx, req.ID)
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

func (u *userUsecase) ListUsers(ctx context.Context, req *xuser.ListUserRequest, filters *xuser.UserFilters) ([]*xuser.User, int64, error) {
	return u.userRepo.ListUsers(ctx, req, filters)
}

func (u *userUsecase) DeleteUsersCreatedBefore(_ context.Context, days time.Time) error {
	u.logger.Info("Deleting users created before", xlogger.Object("days", days))

	return nil
}

func (u *userUsecase) UploadLocal(ctx context.Context, req *xuser.UploadAvatarLocalRequest) error {
	if _, err := u.userRepo.GetUserByID(ctx, req.UserID); err != nil {
		return err
	}

	err := u.queue.PublishMessage(
		ctx,
		consts.UploadUserAvatarJobType,
		&xuser.UploadAvatarLocalQueuePayload{
			UserID: req.UserID,
			File:   req.File,
		},
	)
	if err != nil {
		u.logger.Error("Publish upload avatar job failed", xlogger.Error(err))
		return err
	}

	u.logger.Info("Upload avatar job published",
		xlogger.Uint("user_id", req.UserID),
	)

	return nil
}

func (u *userUsecase) ProcessUploadLocal(ctx context.Context, req *xuser.UploadAvatarLocalInput) error {
	userEntity, err := u.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	oldAvatar := userEntity.Avatar
	userEntity.Avatar = req.Filename

	if _, err := u.userRepo.UpdateUser(ctx, userEntity); err != nil {
		return err
	}

	if oldAvatar != "" {
		oldPath := filepath.Join("attachments/images/avatar", oldAvatar)
		_ = u.fileService.Delete(oldPath)
	}

	return nil
}

func (u *userUsecase) UploadCloud(ctx context.Context, req *xuser.UploadAvatarCloudRequest) error {
	if req.File == nil {
		return xhttp.BadRequestErrorf("File is required")
	}

	user, err := u.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if err := u.queue.PublishMessage(
		ctx,
		consts.UploadAvatarCloudJobType,
		&xuser.UploadAvatarCloudQueuePayload{
			UserID:    uint(user.ID),
			File:      req.File,
			OldAvatar: user.Avatar,
		},
	); err != nil {
		u.logger.Error(
			"Publish upload cloud avatar job failed",
			xlogger.Error(err),
			xlogger.Uint("user_id", uint(user.ID)),
		)
		return err
	}

	u.logger.Info(
		"Upload avatar cloud job published",
		xlogger.Uint("user_id", uint(user.ID)),
	)

	return nil
}
func (u *userUsecase) ProcessUploadCloud(ctx context.Context, req *xuser.UploadAvatarCloudInput) error {
	user, err := u.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	oldAvatar := user.Avatar
	user.Avatar = req.SecureURL

	if _, err := u.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}

	if oldAvatar != "" {
		if err := u.queue.PublishMessage(
			ctx,
			consts.DeleteCloudinaryAssetJobType,
			&xuser.DeleteCloudAssetPayload{
				PublicID: oldAvatar,
			},
		); err != nil {
			u.logger.Warn(
				"Publish delete old cloud avatar job failed",
				xlogger.Error(err),
				xlogger.String("public_id", oldAvatar),
			)
		}
	}

	return nil
}
