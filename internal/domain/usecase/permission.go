package usecase

import "context"

type PermissionUsecase interface {
	CheckPermission(ctx context.Context, roleID int, method string, endpoint string) (bool, error)
}
