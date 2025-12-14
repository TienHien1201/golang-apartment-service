package repository

import (
	"context"

	"thomas.vn/hr_recruitment/internal/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, req *model.ListUserRequest) ([]*model.User, int64, error)
}
