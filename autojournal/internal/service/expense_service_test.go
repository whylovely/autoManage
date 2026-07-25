package service

import (
	"autojournal/internal/domain"
	"context"
	"errors"
	"testing"
	"time"
)

type expenseRepoStub struct{}

func (expenseRepoStub) Create(context.Context, *domain.Expense) error { return nil }
func (expenseRepoStub) Delete(context.Context, int64) error           { return nil }
func (expenseRepoStub) GetByID(context.Context, int64) (*domain.Expense, error) {
	return nil, nil
}
func (expenseRepoStub) ListByVehicle(context.Context, int64) ([]domain.Expense, error) {
	return nil, nil
}
func (expenseRepoStub) ListByVehicleAndPeriod(context.Context, int64, time.Time, time.Time) ([]domain.Expense, error) {
	return nil, nil
}
func (expenseRepoStub) SumByVehicle(context.Context, int64) (int64, error) { return 0, nil }
func (expenseRepoStub) TotalsByVehicleCategory(context.Context, int64) ([]domain.ExpenseCategoryTotal, error) {
	return nil, nil
}

func TestExpenseService_AddExpense_RejectsNonPositiveAmount(t *testing.T) {
	service := NewExpenseService(expenseRepoStub{}, nil, nil)

	err := service.AddExpense(context.Background(), &domain.Expense{
		VehicleID:  1,
		CategoryID: 1,
		Amount:     0,
		Date:       time.Now(),
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}
