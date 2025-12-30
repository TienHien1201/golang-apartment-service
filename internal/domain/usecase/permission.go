package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
)

type PermissionUsecase interface {
	CheckPermission(ctx context.Context, request model.CheckPermissionRequest) (bool, error)
	CreatePermission(ctx context.Context, req *model.CreatePermissionRequest, userID int) (*model.Permission, error)
}
