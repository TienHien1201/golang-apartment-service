package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"thomas.vn/apartment_service/internal/domain/model"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xutils "thomas.vn/apartment_service/pkg/utils"
)

type PermissionRepository struct {
	logger *xlogger.Logger
	db     *gorm.DB
}

func NewPermissionRepository(
	logger *xlogger.Logger,
	db *gorm.DB,
) *PermissionRepository {
	return &PermissionRepository{
		logger: logger,
		db:     db,
	}
}

func (r *PermissionRepository) HasPermission(ctx context.Context, request model.CheckPermissionRequest) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("role_permission rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where(`
			rp.role_id = ?
			AND rp.is_active = 1
			AND p.method = ?
			AND p.endpoint = ?
		`, request.RoleID, request.Method, request.Endpoint).
		Count(&count).Error

	if err != nil {
		r.logger.Error("Check permission failed", xlogger.Error(err))
		return false, err
	}

	return count > 0, nil
}

func (r *PermissionRepository) CreatePermission(ctx context.Context, permission *model.Permission) (*model.Permission, error) {
	permission.CreatedAt = xutils.GetTimeNow()
	permission.UpdatedAt = xutils.GetTimeNow()

	result := r.db.WithContext(ctx).Create(permission)
	if result.Error != nil {
		r.logger.Error("Create permission failed", xlogger.Error(result.Error))
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("Create permission failed, no rows affected")
	}
	return permission, nil
}
