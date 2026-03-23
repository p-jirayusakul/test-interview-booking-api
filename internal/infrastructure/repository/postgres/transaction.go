package postgres

import (
	"context"

	"github.com/p-jirayusakul/test-interview-booking-api/internal/domain"
	"gorm.io/gorm"
)

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) domain.TransactionManager {
	return &transactionManager{db: db}
}

func (t *transactionManager) WithTransaction(
	ctx context.Context,
	fn func(tx *gorm.DB) error,
) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
