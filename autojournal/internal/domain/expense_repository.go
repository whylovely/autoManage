package domain

import (
	"context"
)

type ExpenseRepository interface {
	Create(ctx context.Context, expense *Expense) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Expense, error)
	ListByVehicle(ctx context.Context, vehicleID int64) ([]Expense, error)
}
