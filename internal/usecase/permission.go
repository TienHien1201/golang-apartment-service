package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/consts"
	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type permissionUsecase struct {
	logger *xlogger.Logger
	repo   repository.PermissionRepository
}

func NewPermissionUsecase(logger *xlogger.Logger, repo repository.PermissionRepository) usecase.PermissionUsecase {
	return &permissionUsecase{logger: logger, repo: repo}
}

func (u *permissionUsecase) CheckPermission(ctx context.Context, request model.CheckPermissionRequest) (bool, error) {
	if request.RoleID == consts.UserAdmin {
		return true, nil
	}
	return u.repo.HasPermission(ctx, request)
}
func (u *permissionUsecase) CreatePermission(ctx context.Context, req *model.CreatePermissionRequest, userID int) (*model.Permission, error) {
	permission := &model.Permission{
		Name:      req.Name,
		Endpoint:  req.Endpoint,
		Method:    req.Method,
		Module:    req.Module,
		CreatedBy: userID,
	}

	return u.repo.CreatePermission(ctx, permission)
}
