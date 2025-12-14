package usecase

import (
	"context"
	"time"

	"thomas.vn/apartment_service/internal/domain/model"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetUser(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, req *model.ListUserRequest) ([]*model.User, int64, error)
	DeleteUsersCreatedBefore(ctx context.Context, days time.Time) error
}
