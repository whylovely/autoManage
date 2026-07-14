package domain

import (
	"context"
)

type ExpenseCategoryRepository interface {
	List(ctx context.Context) ([]ExpenseCategory, error)
	GetByID(ctx context.Context, id int64) (*ExpenseCategory, error)
}
