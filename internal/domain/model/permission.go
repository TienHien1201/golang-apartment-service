package model

import "time"

type Permission struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Endpoint  string     `json:"endpoint"`
	Method    string     `json:"method"`
	Module    string     `json:"module"`
	DeletedBy int        `json:"deleted_by"`
	CreatedBy int        `json:"created_by"`
	IsDeleted int        `json:"is_deleted"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreatePermissionRequest struct {
	Name     string `json:"name" validate:"required"`
	Endpoint string `json:"endpoint" validate:"required"`
	Method   string `json:"method" validate:"required,oneof=GET POST PUT DELETE"`
	Module   string `json:"module" validate:"required"`
}

type CheckPermissionRequest struct {
	RoleID   int    `json:"role_id"`
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
}

func (Permission) TableName() string {
	return "permissions"
}
