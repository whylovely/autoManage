package service

import (
	"autojournal/internal/domain"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

type reportExpenseRepoStub struct {
	expenses []domain.Expense
	total    int64
	category []domain.ExpenseCategoryTotal
	created  *domain.Expense
	deleted  int64
}

func (r *reportExpenseRepoStub) Create(_ context.Context, expense *domain.Expense) error {
	r.created = expense
	return nil
}
func (r *reportExpenseRepoStub) Delete(_ context.Context, id int64) error {
	r.deleted = id
	return nil
}
func (r *reportExpenseRepoStub) GetByID(context.Context, int64) (*domain.Expense, error) {
	if len(r.expenses) == 0 {
		return nil, domain.ErrNotFound
	}
	return &r.expenses[0], nil
}
func (r *reportExpenseRepoStub) ListByVehicle(context.Context, int64) ([]domain.Expense, error) {
	return r.expenses, nil
}
func (r *reportExpenseRepoStub) ListByVehicleAndPeriod(context.Context, int64, time.Time, time.Time) ([]domain.Expense, error) {
	return r.expenses, nil
}
func (r *reportExpenseRepoStub) SumByVehicle(context.Context, int64) (int64, error) {
	return r.total, nil
}
func (r *reportExpenseRepoStub) TotalsByVehicleCategory(context.Context, int64) ([]domain.ExpenseCategoryTotal, error) {
	return r.category, nil
}

type categoryRepoStub struct {
	category *domain.ExpenseCategory
}

func (r categoryRepoStub) List(context.Context) ([]domain.ExpenseCategory, error) {
	return []domain.ExpenseCategory{*r.category}, nil
}
func (r categoryRepoStub) GetByID(context.Context, int64) (*domain.ExpenseCategory, error) {
	return r.category, nil
}

func TestExpenseService_AddExpense_PersistsValidExpense(t *testing.T) {
	repo := &reportExpenseRepoStub{}
	vehicleRepo := &vehicleRepoStub{vehicle: &domain.Vehicle{ID: 1}}
	categoryRepo := categoryRepoStub{category: &domain.ExpenseCategory{ID: 2, Name: "fuel"}}
	service := NewExpenseService(repo, vehicleRepo, categoryRepo)
	expense := &domain.Expense{
		VehicleID: 1, CategoryID: 2, Amount: 1_500, OdometerAt: 5_000,
		Date: time.Date(2026, time.July, 27, 0, 0, 0, 0, time.UTC),
	}

	require.NoError(t, service.AddExpense(context.Background(), expense))
	assert.Same(t, expense, repo.created)
}

func TestExpenseService_GetExpenseStats_AggregatesByMonth(t *testing.T) {
	repo := &reportExpenseRepoStub{
		expenses: []domain.Expense{
			{Amount: 1_000, Date: time.Date(2026, time.February, 2, 0, 0, 0, 0, time.UTC)},
			{Amount: 500, Date: time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)},
			{Amount: 700, Date: time.Date(2026, time.February, 8, 0, 0, 0, 0, time.UTC)},
		},
		total: 2_200,
	}
	service := NewExpenseService(repo, &vehicleRepoStub{vehicle: &domain.Vehicle{ID: 1}}, nil)

	stats, err := service.GetExpenseStats(context.Background(), 1)
	if err != nil {
		t.Fatalf("get expense stats: %v", err)
	}
	if stats.TotalAmount != 2_200 || len(stats.ByMonth) != 2 {
		t.Fatalf("unexpected stats: %#v", stats)
	}
	if stats.ByMonth[0].Month != "2026-01" || stats.ByMonth[0].TotalAmount != 500 {
		t.Fatalf("unexpected first month: %#v", stats.ByMonth[0])
	}
	if stats.ByMonth[1].Month != "2026-02" || stats.ByMonth[1].TotalAmount != 1_700 {
		t.Fatalf("unexpected second month: %#v", stats.ByMonth[1])
	}
}

func TestExpenseService_ExportVehicleExpenses_CSV(t *testing.T) {
	repo := &reportExpenseRepoStub{expenses: []domain.Expense{{
		ID: 1, VehicleID: 2, CategoryID: 3, Amount: 1_500, OdometerAt: 42_000,
		Date: time.Date(2026, time.July, 27, 0, 0, 0, 0, time.UTC), Description: "Oil",
	}}}
	service := NewExpenseService(repo, &vehicleRepoStub{vehicle: &domain.Vehicle{ID: 2}}, nil)
	destination := filepath.Join(t.TempDir(), "expenses.csv")

	if err := service.ExportVehicleExpenses(context.Background(), 2, "csv", destination); err != nil {
		t.Fatalf("export expenses: %v", err)
	}
	content, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	if !strings.Contains(string(content), "vehicle_id") || !strings.Contains(string(content), "Oil") {
		t.Fatalf("unexpected CSV export: %s", content)
	}
}

func TestExpenseService_ExportVehicleExpenses_JSON(t *testing.T) {
	repo := &reportExpenseRepoStub{expenses: []domain.Expense{{
		ID: 1, VehicleID: 2, CategoryID: 3, Amount: 1_500,
		Date: time.Date(2026, time.July, 27, 0, 0, 0, 0, time.UTC),
	}}}
	service := NewExpenseService(repo, &vehicleRepoStub{vehicle: &domain.Vehicle{ID: 2}}, nil)
	destination := filepath.Join(t.TempDir(), "expenses.json")

	require.NoError(t, service.ExportVehicleExpenses(context.Background(), 2, "json", destination))
	content, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"amount": 1500`)

	err = service.ExportVehicleExpenses(context.Background(), 2, "xml", filepath.Join(t.TempDir(), "bad.xml"))
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestExpenseService_ReadAndDeleteOperations(t *testing.T) {
	date := time.Date(2026, time.July, 27, 0, 0, 0, 0, time.UTC)
	repo := &reportExpenseRepoStub{
		expenses: []domain.Expense{{ID: 8, VehicleID: 1, Amount: 1_500, Date: date}},
		total:    1_500,
		category: []domain.ExpenseCategoryTotal{{CategoryID: 2, CategoryName: "fuel", TotalAmount: 1_500}},
	}
	service := NewExpenseService(repo, &vehicleRepoStub{vehicle: &domain.Vehicle{ID: 1}}, nil)

	got, err := service.GetExpense(context.Background(), 8)
	require.NoError(t, err)
	assert.EqualValues(t, 8, got.ID)
	list, err := service.ListVehicleExpenses(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, repo.expenses, list)
	period, err := service.ListVehicleExpensesForPeriod(context.Background(), 1, date.Add(-time.Hour), date.Add(time.Hour))
	require.NoError(t, err)
	assert.Equal(t, repo.expenses, period)
	total, err := service.GetVehicleExpenseTotal(context.Background(), 1)
	require.NoError(t, err)
	assert.EqualValues(t, 1_500, total)
	totals, err := service.GetVehicleExpenseTotalsByCategory(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, repo.category, totals)
	require.NoError(t, service.DeleteExpense(context.Background(), 8))
	assert.EqualValues(t, 8, repo.deleted)
}

func TestExpenseServices_RejectInvalidInput(t *testing.T) {
	service := NewExpenseService(&reportExpenseRepoStub{}, &vehicleRepoStub{}, nil)
	assert.ErrorIs(t, service.DeleteExpense(context.Background(), 0), domain.ErrValidation)
	_, err := service.GetExpense(context.Background(), -1)
	assert.ErrorIs(t, err, domain.ErrValidation)
	_, err = service.ListVehicleExpenses(context.Background(), 0)
	assert.ErrorIs(t, err, domain.ErrValidation)
	_, err = service.ListVehicleExpensesForPeriod(context.Background(), 1, time.Now(), time.Time{})
	assert.ErrorIs(t, err, domain.ErrValidation)
	_, err = service.GetVehicleExpenseTotal(context.Background(), 0)
	assert.ErrorIs(t, err, domain.ErrValidation)
	_, err = service.GetVehicleExpenseTotalsByCategory(context.Background(), 0)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestExpenseCategoryService_ListCategories(t *testing.T) {
	repo := categoryRepoStub{category: &domain.ExpenseCategory{ID: 1, Name: "fuel"}}
	service := NewExpenseCategory(repo)

	categories, err := service.ListCategories(context.Background())
	require.NoError(t, err)
	require.Len(t, categories, 1)
	assert.Equal(t, "fuel", categories[0].Name)
}
