package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	xuser "thomas.vn/apartment_service/internal/domain/model/user"

	"thomas.vn/apartment_service/internal/domain/repository"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xutils "thomas.vn/apartment_service/pkg/utils"
)

type userRepository struct {
	logger    *xlogger.Logger
	userTable *gorm.DB
}

func NewUserRepository(logger *xlogger.Logger, db *gorm.DB) repository.UserRepository {
	return &userRepository{
		logger:    logger,
		userTable: db.Table("users"),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *xuser.User) (*xuser.User, error) {
	user.CreatedAt = xutils.GetTimeNow()
	user.UpdatedAt = xutils.GetTimeNow()

	result := r.userTable.WithContext(ctx).Create(user)
	if result.Error != nil {
		r.logger.Error("Create user failed", xlogger.Error(result.Error))
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("create user failed")
	}

	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*xuser.User, error) {
	var user xuser.User
	result := r.userTable.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("Get user by id failed", xlogger.Error(result.Error))
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*xuser.User, error) {
	var user xuser.User
	result := r.userTable.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("Get user by email failed", xlogger.Error(result.Error))
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *xuser.User) (*xuser.User, error) {
	user.UpdatedAt = xutils.GetTimeNow()

	result := r.userTable.WithContext(ctx).Save(user)
	if result.Error != nil {
		r.logger.Error("Update user failed", xlogger.Error(result.Error))
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("update user failed")
	}

	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	result := r.userTable.WithContext(ctx).Where("id = ?", id).Delete(&xuser.User{})
	if result.Error != nil {
		r.logger.Error("Delete user failed", xlogger.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("delete user failed")
	}

	return nil
}

func (r *userRepository) ListUsers(ctx context.Context, req *xuser.ListUserRequest) ([]*xuser.User, int64, error) {
	var users []*xuser.User
	var total int64

	query := r.userTable.WithContext(ctx)

	// Apply filters
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.FromDate != "" {
		query = query.Where(req.RangeBy+" >= ?", req.FromDate+" 00:00:00")
	}
	if req.ToDate != "" {
		query = query.Where(req.RangeBy+" <= ?", req.ToDate+" 23:59:59")
	}

	// Get total count if not exclude
	if !req.ExcludeTotal {
		if err := query.Count(&total).Error; err != nil {
			r.logger.Error("Count users failed", xlogger.Error(err))
			return nil, 0, err
		}
	}

	// Apply pagination
	if req.Page > 0 && req.Limit > 0 {
		query = query.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}

	// Apply sorting
	if req.SortBy != "" && req.OrderBy != "" {
		query = query.Order(req.SortBy + " " + req.OrderBy)
	}

	// Execute query
	if err := query.Find(&users).Error; err != nil {
		r.logger.Error("List users failed", xlogger.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) UpdateTotpSecret(
	ctx context.Context,
	userID int64,
	secret *string,
) error {
	result := r.userTable.
		WithContext(ctx).
		Model(&xuser.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"totp_secret": secret,
			"updated_at":  xutils.GetTimeNow(),
		})

	if result.Error != nil {
		r.logger.Error(
			"Update totp secret failed",
			xlogger.Int64("user_id", userID),
			xlogger.Error(result.Error),
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("update totp secret failed: user not found")
	}

	return nil
}
