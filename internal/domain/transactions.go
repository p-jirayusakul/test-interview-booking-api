package domain

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}
