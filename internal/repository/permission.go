package repository

import (
	"context"

	"gorm.io/gorm"
	xlogger "thomas.vn/apartment_service/pkg/logger"
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

func (r *PermissionRepository) HasPermission(ctx context.Context, roleID int, method string, endpoint string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("role_permission rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where(`
			rp.role_id = ?
			AND rp.is_active = 1
			AND p.method = ?
			AND p.endpoint = ?
		`, roleID, method, endpoint).
		Count(&count).Error

	if err != nil {
		r.logger.Error("Check permission failed", xlogger.Error(err))
		return false, err
	}

	return count > 0, nil
}
