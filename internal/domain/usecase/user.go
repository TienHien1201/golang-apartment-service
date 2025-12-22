package usecase

import (
	"context"
	"time"

	xuser "thomas.vn/apartment_service/internal/domain/model/user"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *xuser.CreateUserRequest) (*xuser.User, error)
	GetUser(ctx context.Context, id uint) (*xuser.User, error)
	UpdateUser(ctx context.Context, req *xuser.UpdateUserRequest) (*xuser.User, error)
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, req *xuser.ListUserRequest) ([]*xuser.User, int64, error)
	DeleteUsersCreatedBefore(ctx context.Context, days time.Time) error
	UploadLocal(ctx context.Context, req *xuser.UploadAvatarLocalRequest) error
	ProcessUploadLocal(ctx context.Context, req *xuser.UploadAvatarLocalInput) error
	UploadCloud(ctx context.Context, req *xuser.UploadAvatarCloudRequest) error
	ProcessUploadCloud(ctx context.Context, req *xuser.UploadAvatarCloudInput) error
}
