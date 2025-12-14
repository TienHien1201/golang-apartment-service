package model

import (
	"time"

	xhttp "thomas.vn/hr_recruitment/pkg/http"
)

type User struct {
	ID         int        `gorm:"primaryKey;column:id" json:"id" example:"1"`
	Email      string     `gorm:"column:email;uniqueIndex:uniq_email;size:255" json:"email" validate:"required,email" example:"abc@host.com"`
	FullName   string     `gorm:"column:full_name;size:255" json:"full_name,omitempty" example:"John Doe"`
	Avatar     *string    `gorm:"column:avatar;size:255" json:"avatar,omitempty" example:"https://avatar.com/abc.jpg"`
	Password   string     `gorm:"column:password;size:255" json:"-"` // không trả ra JSON
	FacebookID string     `gorm:"column:facebook_id;uniqueIndex:uniq_facebook_id;size:255" json:"facebook_id,omitempty"`
	GoogleID   string     `gorm:"column:google_id;uniqueIndex:uniq_google_id;size:255" json:"google_id,omitempty"`
	TotpSecret *string    `gorm:"column:totp_secret;size:255" json:"totp_secret,omitempty"`
	RoleID     int        `gorm:"column:role_id;index:idx_role_id;default:2" json:"role_id" example:"2"`
	DeletedBy  int        `gorm:"column:deleted_by;default:0" json:"deleted_by"`
	IsDeleted  bool       `gorm:"column:is_deleted;default:false" json:"is_deleted"`
	IsActive   int        `gorm:"column:is_active;default:1" json:"is_active"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	CreatedAt  time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at" example:"2025-01-01T10:00:00Z"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at" example:"2025-01-01T10:00:00Z"`

	// Quan hệ (nếu cần dùng trong Go)
	// Role             Role               `gorm:"foreignKey:RoleID;references:id;constraint:OnUpdate:RESTRICT"`
	// ChatGroupMembers []ChatGroupMember  `gorm:"foreignKey:UserID"`
	// ChatGroups       []ChatGroup        `gorm:"foreignKey:OwnerID"`
	// ChatMessages     []ChatMessage      `gorm:"foreignKey:UserID"`
	// Chats            []Chat             `gorm:"foreignKey:UserID"`
}

type UserIDRequest struct {
	ID int `json:"id" param:"id" swaggerignore:"true" validate:"required,gt=0"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email" example:"abc@host.com"`
	Password string `json:"password" validate:"required,min=8" example:"password"`
	FullName string `json:"full_name" validate:"required" example:"John Doe"`
	RoleID   int    `json:"role_id" validate:"required,oneof=1 2" example:"user"`
}

type UpdateUserRequest struct {
	UserIDRequest
	Password string `json:"password" validate:"omitempty,min=8" example:"password"`
	FullName string `json:"full_name" validate:"omitempty" example:"John Doe"`
	RoleID   int    `json:"role_id" validate:"omitempty,oneof=1 2" example:"user"`
	IsActive int    `json:"is_active" validate:"omitempty,oneof=0 1" example:"1"`
}

type ListUserRequest struct {
	xhttp.PaginationOptions
	xhttp.DateRangeOptions
	xhttp.SortOptions
	Status int `query:"status" validate:"omitempty,oneof=1 2"`
}

type ListUserResponse struct {
	Rows  []*User `json:"rows"`
	Total int64   `json:"total,omitempty"`
}

func (User) TableName() string {
	return "users"
}
