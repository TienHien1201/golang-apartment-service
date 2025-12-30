package repository

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
)

type PermissionRepository interface {
	HasPermission(ctx context.Context, request model.CheckPermissionRequest) (bool, error)
	CreatePermission(ctx context.Context, permission *model.Permission) (*model.Permission, error)
}
