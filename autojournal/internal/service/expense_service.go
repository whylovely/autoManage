package service

import (
	"autojournal/internal/domain"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ExpenseService struct {
	expenseRepo  domain.ExpenseRepository
	vehicleRepo  domain.VehicleRepository
	categoryRepo domain.ExpenseCategoryRepository
}

type ExpenseCategoryService struct {
	repo domain.ExpenseCategoryRepository
}

// ExpenseCategorySevice is kept as an alias for backward compatibility.
type ExpenseCategorySevice = ExpenseCategoryService

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

func (s *ExpenseService) GetExpenseStats(ctx context.Context, vehicleID int64) (*domain.ExpenseStats, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}
	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	total, err := s.expenseRepo.SumByVehicle(ctx, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("sum vehicle expenses: %w", err)
	}
	byCategory, err := s.expenseRepo.TotalsByVehicleCategory(ctx, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("get expense totals by category: %w", err)
	}
	expenses, err := s.expenseRepo.ListByVehicle(ctx, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("list vehicle expenses: %w", err)
	}

	monthly := make(map[string]int64)
	for _, expense := range expenses {
		monthly[expense.Date.Format("2006-01")] += expense.Amount
	}
	months := make([]string, 0, len(monthly))
	for month := range monthly {
		months = append(months, month)
	}
	sort.Strings(months)

	byMonth := make([]domain.MonthlyExpenseTotal, 0, len(months))
	for _, month := range months {
		byMonth = append(byMonth, domain.MonthlyExpenseTotal{
			Month:       month,
			TotalAmount: monthly[month],
		})
	}

	return &domain.ExpenseStats{
		TotalAmount: total,
		ByCategory:  byCategory,
		ByMonth:     byMonth,
	}, nil
}

func (s *ExpenseService) ExportVehicleExpenses(ctx context.Context, vehicleID int64, format, destination string) error {
	if strings.TrimSpace(destination) == "" {
		return fmt.Errorf("%w: export destination is required", domain.ErrValidation)
	}
	normalizedFormat := strings.ToLower(strings.TrimSpace(format))
	if normalizedFormat != "csv" && normalizedFormat != "json" {
		return fmt.Errorf("%w: export format must be csv or json", domain.ErrValidation)
	}

	expenses, err := s.ListVehicleExpenses(ctx, vehicleID)
	if err != nil {
		return err
	}

	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("%w: create export file: %w", domain.ErrInfrastructure, err)
	}
	defer file.Close()

	switch normalizedFormat {
	case "json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(expenses); err != nil {
			return fmt.Errorf("%w: encode JSON export: %w", domain.ErrInfrastructure, err)
		}
	case "csv":
		writer := csv.NewWriter(file)
		if err := writer.Write([]string{"id", "vehicle_id", "category_id", "amount", "odometer", "date", "description"}); err != nil {
			return fmt.Errorf("%w: write CSV header: %w", domain.ErrInfrastructure, err)
		}
		for _, expense := range expenses {
			record := []string{
				strconv.FormatInt(expense.ID, 10),
				strconv.FormatInt(expense.VehicleID, 10),
				strconv.FormatInt(expense.CategoryID, 10),
				strconv.FormatInt(expense.Amount, 10),
				strconv.FormatInt(expense.OdometerAt, 10),
				expense.Date.Format(time.RFC3339),
				expense.Description,
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("%w: write CSV export: %w", domain.ErrInfrastructure, err)
			}
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			return fmt.Errorf("%w: flush CSV export: %w", domain.ErrInfrastructure, err)
		}
	}

	return nil
}

func NewExpenseCategory(repo domain.ExpenseCategoryRepository) *ExpenseCategoryService {
	return &ExpenseCategoryService{repo: repo}
}

func (s *ExpenseCategoryService) ListCategories(ctx context.Context) ([]domain.ExpenseCategory, error) {
	categories, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list expense categories: %w", err)
	}

	return categories, nil
}
