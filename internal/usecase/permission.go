package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
)

type permissionUsecase struct {
	repo repository.PermissionRepository
}

func NewPermissionUsecase(
	repo repository.PermissionRepository,
) usecase.PermissionUsecase {
	return &permissionUsecase{repo: repo}
}

func (u *permissionUsecase) CheckPermission(
	ctx context.Context,
	roleID int,
	method string,
	endpoint string,
) (bool, error) {
	if roleID == 1 {
		return true, nil
	}

	return u.repo.HasPermission(ctx, roleID, method, endpoint)
}
