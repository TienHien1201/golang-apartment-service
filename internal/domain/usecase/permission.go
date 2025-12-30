package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
)

type PermissionUsecase interface {
	CheckPermission(ctx context.Context, request model.CheckPermissionRequest) (bool, error)
	CreatePermission(ctx context.Context, req *model.CreatePermissionRequest, userID int) (*model.Permission, error)
	GetPermissionByID(ctx context.Context, permissionID uint) (*model.Permission, error)
	UpdatePermission(ctx context.Context, req *model.UpdatePermissionRequest) (*model.Permission, error)
}
