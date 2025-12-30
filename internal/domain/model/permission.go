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

type UpdatePermissionRequest struct {
	PermissionIDRequest
	Name     string `json:"name" validate:"required"`
	Endpoint string `json:"endpoint" validate:"required"`
	Method   string `json:"method" validate:"required,oneof=GET POST PUT DELETE"`
	Module   string `json:"module" validate:"required"`
}

type PermissionIDRequest struct {
	ID uint `json:"id" param:"id" swaggerignore:"true" validate:"required,gt=0"`
}
type CheckPermissionRequest struct {
	RoleID   int    `json:"role_id"`
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
}

func (Permission) TableName() string {
	return "permissions"
}
