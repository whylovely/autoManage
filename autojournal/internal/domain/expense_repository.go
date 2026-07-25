package domain

import (
	"context"
	"time"
)

type ExpenseRepository interface {
	Create(ctx context.Context, expense *Expense) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Expense, error)
	ListByVehicle(ctx context.Context, vehicleID int64) ([]Expense, error)
	ListByVehicleAndPeriod(ctx context.Context, vehicleID int64, from, to time.Time) ([]Expense, error)
	SumByVehicle(ctx context.Context, vehicleID int64) (int64, error)
	TotalsByVehicleCategory(ctx context.Context, vehicleID int64) ([]ExpenseCategoryTotal, error)
}
