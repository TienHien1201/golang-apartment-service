package repository

import "context"

type PermissionRepository interface {
	HasPermission(ctx context.Context, roleID int, method string,
		endpoint string) (bool, error)
}
