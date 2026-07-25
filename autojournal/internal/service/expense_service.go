package service

import (
	"autojournal/internal/domain"
	"context"
	"fmt"
	"time"
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
		return fmt.Errorf("%w: expense is required", domain.ErrValidation)
	}

	if expense.VehicleID <= 0 {
		return fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}
	if expense.CategoryID <= 0 {
		return fmt.Errorf("%w: category id must be positive", domain.ErrValidation)
	}
	if expense.Amount <= 0 {
		return fmt.Errorf("%w: amount must be positive", domain.ErrValidation)
	}
	if expense.OdometerAt < 0 {
		return fmt.Errorf("%w: odometer cannot be negative", domain.ErrValidation)
	}

	if expense.Date.IsZero() {
		return fmt.Errorf("%w: expense date is required", domain.ErrValidation)
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
		return fmt.Errorf("%w: expense id must be positive", domain.ErrValidation)
	}

	return s.expenseRepo.Delete(ctx, id)
}

func (s *ExpenseService) GetExpense(ctx context.Context, id int64) (*domain.Expense, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: expense id must be positive", domain.ErrValidation)
	}

	return s.expenseRepo.GetByID(ctx, id)
}

func (s *ExpenseService) ListVehicleExpenses(ctx context.Context, vehicleID int64) ([]domain.Expense, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}

	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	return s.expenseRepo.ListByVehicle(ctx, vehicleID)
}

func (s *ExpenseService) ListVehicleExpensesForPeriod(
	ctx context.Context,
	vehicleID int64,
	from, to time.Time,
) ([]domain.Expense, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}
	if from.IsZero() || to.IsZero() {
		return nil, fmt.Errorf("%w: period bounds are required", domain.ErrValidation)
	}
	if from.After(to) {
		return nil, fmt.Errorf("%w: period start must not be after period end", domain.ErrValidation)
	}

	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	return s.expenseRepo.ListByVehicleAndPeriod(ctx, vehicleID, from, to)
}

func (s *ExpenseService) GetVehicleExpenseTotal(ctx context.Context, vehicleID int64) (int64, error) {
	if vehicleID <= 0 {
		return 0, fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}

	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		return 0, fmt.Errorf("get vehicle: %w", err)
	}

	a, err := s.expenseRepo.SumByVehicle(ctx, vehicleID)
	if err != nil {
		return 0, fmt.Errorf("sum vehicle expenses: %w", err)
	}

	return a, nil
}

func (s *ExpenseService) GetVehicleExpenseTotalsByCategory(ctx context.Context, vehicleID int64) ([]domain.ExpenseCategoryTotal, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}

	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	totals, err := s.expenseRepo.TotalsByVehicleCategory(ctx, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("get expense totals by category: %w", err)
	}

	return totals, nil
}

func NewExpenseCategory(repo domain.ExpenseCategoryRepository) *ExpenseCategorySevice {
	return &ExpenseCategorySevice{repo: repo}
}
