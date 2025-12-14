package repository

import (
	"context"

	"gorm.io/gorm"
)

type ITransaction interface {
	Begin(ctx context.Context) (*gorm.DB, error)
	Commit(ctx context.Context, tx *gorm.DB) error
	Rollback(ctx context.Context, tx *gorm.DB) error
}
