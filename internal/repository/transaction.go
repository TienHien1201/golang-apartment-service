package repository

import (
	"context"

	"gorm.io/gorm"

	"thomas.vn/apartment_service/internal/domain/repository"
)

type Transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) repository.ITransaction {
	return &Transaction{db: db}
}

func (t *Transaction) Begin(ctx context.Context) (*gorm.DB, error) {
	tx := t.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (t *Transaction) Commit(ctx context.Context, tx *gorm.DB) error {
	return tx.WithContext(ctx).Commit().Error
}

func (t *Transaction) Rollback(ctx context.Context, tx *gorm.DB) error {
	return tx.WithContext(ctx).Rollback().Error
}
