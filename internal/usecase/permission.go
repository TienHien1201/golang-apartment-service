package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/consts"
	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
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
func (u *permissionUsecase) GetPermissionByID(ctx context.Context, permissionID uint) (*model.Permission, error) {
	permission, err := u.repo.GetPermissionByID(ctx, permissionID)
	if err != nil {
		u.logger.Error("Failed to get permission", xlogger.Error(err))
		return nil, err
	}
	if permission == nil {
		return nil, xhttp.NotFoundErrorf("User with ID %d not found", permissionID)
	}

	return permission, nil
}

func (u *permissionUsecase) UpdatePermission(ctx context.Context, req *model.UpdatePermissionRequest) (*model.Permission, error) {
	permission, err := u.GetPermissionByID(ctx, req.ID)
	if err != nil {
		u.logger.Error("Failed to get permission", xlogger.Error(err))
		return nil, err
	}
	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Endpoint != "" {
		permission.Endpoint = req.Endpoint
	}
	if req.Method != "" {
		permission.Method = req.Method
	}
	updatedPermission, err := u.repo.UpdatePermission(ctx, permission)
	if err != nil {
		u.logger.Error("Failed to update permission", xlogger.Error(err))
		return nil, err
	}
	return updatedPermission, nil

}
