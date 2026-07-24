package service

import (
	"autojournal/internal/domain"
	"context"
	"fmt"
)

type ExpenseService struct {
	expenseRepo  domain.ExpenseRepository
	vehicleRepo  domain.VehicleRepository
	categoryRepo domain.ExpenseCategoryRepository
}

type ExpenseCategorySevice struct {
	repo domain.ExpenseCategoryRepository
}

func NewExpenseService(
	expenseRepo domain.ExpenseRepository,
	vehicleRepo domain.VehicleRepository,
	categoryRepo domain.ExpenseCategoryRepository,
) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo, vehicleRepo: vehicleRepo, categoryRepo: categoryRepo}
}

func (s *ExpenseService) AddExpense(ctx context.Context, expense *domain.Expense) error {
	if expense == nil {
		return fmt.Errorf("Expense is null")
	}

	if expense.VehicleID <= 0 {
		return fmt.Errorf("vehicle id must be positive")
	}
	if expense.CategoryID <= 0 {
		return fmt.Errorf("category id must be positive")
	}
	if expense.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if expense.OdometerAt < 0 {
		return fmt.Errorf("odometer cannot be negative")
	}

	if expense.Date.IsZero() {
		return fmt.Errorf("expense date is required")
	}

	if _, err := s.vehicleRepo.GetByID(ctx, expense.VehicleID); err != nil {
		return fmt.Errorf("get vehicle: %w", err)
	}
	if _, err := s.categoryRepo.GetByID(ctx, expense.CategoryID); err != nil {
		return fmt.Errorf("get category: %w", err)
	}

	return s.expenseRepo.Create(ctx, expense)
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("expense id must be positive")
	}

	return s.expenseRepo.Delete(ctx, id)
}

func (s *ExpenseService) GetExpense(ctx context.Context, id int64) (*domain.Expense, error) {
	if id <= 0 {
		return nil, fmt.Errorf("expense id must be positive")
	}

	return s.expenseRepo.GetByID(ctx, id)
}

func (s *ExpenseService) ListVehicleExpenses(ctx context.Context, vehicleID int64) ([]domain.Expense, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("vehicle id must be positive")
	}

	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	return s.expenseRepo.ListByVehicle(ctx, vehicleID)
}

func (s *ExpenseService) GetVehicleExpenseTotal(ctx context.Context, vehicleID int64) (int64, error) {
	if vehicleID <= 0 {
		return 0, fmt.Errorf("vehicle id must be positive")
	}

	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		return 0, fmt.Errorf("No vehicle with id: %w", err)
	}

	a, err := s.expenseRepo.SumByVehicle(ctx, vehicleID)
	if err != nil {
		return 0, fmt.Errorf("trouble with repo: %w", err)
	}

	return a, nil
}

func NewExpenseCategory(repo domain.ExpenseCategoryRepository) *ExpenseCategorySevice {
	return &ExpenseCategorySevice{repo: repo}
}
