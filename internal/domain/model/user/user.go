package xuser

import (
	"time"

	"thomas.vn/apartment_service/pkg/query"
)

type User struct {
	ID         int        `json:"id" gorm:"primary_key" example:"1"`
	Email      string     `json:"email" gorm:"unique" example:"abc@host.com"`
	FullName   string     `json:"full_name" example:"John Doe"`
	Avatar     string     `json:"avatar,omitempty" example:"https://avatar.com/abc.jpg"`
	Password   string     `json:"password" example:"password"`
	FacebookID *string    `json:"facebook_id" gorm:"unique"`
	GoogleID   *string    `json:"google_id" gorm:"unique"`
	TotpSecret *string    `json:"totp_secret" example:"secret"`
	RoleID     int        `json:"role_id" example:"2"`
	DeletedBy  int        `json:"deleted_by" example:"1"`
	IsDeleted  bool       `json:"is_deleted" example:"0"`
	IsActive   int        `json:"is_active" example:"1"`
	DeletedAt  *time.Time `json:"deleted_at" example:"2020-09-06T10:10:10Z"`
	CreatedAt  time.Time  `json:"created_at" example:"2025-01-01T10:00:00Z"`
	UpdatedAt  time.Time  `json:"updated_at" example:"2025-01-01T10:00:00Z"`
}

type UserIDRequest struct {
	ID uint `json:"id" param:"id" swaggerignore:"true" validate:"required,gt=0"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email" example:"abc@host.com"`
	Password string `json:"password" validate:"required,min=8" example:"password"`
	FullName string `json:"full_name" validate:"required" example:"John Doe"`
}

type UpdateUserRequest struct {
	UserIDRequest
	Password string `json:"password" validate:"omitempty,min=8" example:"password"`
	FullName string `json:"full_name" validate:"omitempty" example:"John Doe"`
	RoleID   int    `json:"role_id" validate:"omitempty,oneof=1 2" example:"user"`
	IsActive int    `json:"is_active" validate:"omitempty,oneof=0 1" example:"1"`
}
type UserFilters struct {
	FullName string `json:"full_name"`
}

type ListUserRequest struct {
	query.PaginationOptions
	query.DateRangeOptions
	query.SortOptions

	IsDeleted int    `query:"is_deleted"`
	Filters   string `query:"filters"`
}

type ListUserResponse struct {
	Rows  []*User `json:"rows"`
	Total int64   `json:"total,omitempty"`
}

func (User) TableName() string {
	return "users"
}
