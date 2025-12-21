package repository

import (
	"context"

	xuser "thomas.vn/apartment_service/internal/domain/model/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *xuser.User) (*xuser.User, error)
	GetUserByID(ctx context.Context, id uint) (*xuser.User, error)
	GetUserByEmail(ctx context.Context, email string) (*xuser.User, error)
	UpdateUser(ctx context.Context, user *xuser.User) (*xuser.User, error)
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, req *xuser.ListUserRequest) ([]*xuser.User, int64, error)
}
