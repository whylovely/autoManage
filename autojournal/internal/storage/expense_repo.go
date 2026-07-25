package storage

import (
	"autojournal/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ExpenseRepo struct {
	db *sqlx.DB
}

func NewExpenseRepo(db *sqlx.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Create(ctx context.Context, expense *domain.Expense) error {
	const query = `
			INSERT INTO expenses (
					vehicle_id, category_id, amount,
					odometer_at, date, description)
			VALUES (?, ?, ?, ?, ?, ?)
			RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		expense.VehicleID,
		expense.CategoryID,
		expense.Amount,
		expense.OdometerAt,
		expense.Date,
		expense.Description,
	).Scan(&expense.ID)
	if err != nil {
		return fmt.Errorf("%w: create expense: %w", domain.ErrInfrastructure, err)
	}

	return nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM expenses WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("%w: delete expense: %w", domain.ErrInfrastructure, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: get deleted rows: %w", domain.ErrInfrastructure, err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: expense %d", domain.ErrNotFound, id)
	}

	return nil
}

func (r *ExpenseRepo) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	const query = `
		SELECT
		id, vehicle_id, category_id, amount, odometer_at,
			date, description, created_at
		FROM expenses
		WHERE id = ?`

	var expense domain.Expense
	if err := r.db.GetContext(ctx, &expense, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: expense %d", domain.ErrNotFound, id)
		}
		return nil, fmt.Errorf("%w: get expense: %w", domain.ErrInfrastructure, err)
	}

	return &expense, nil
}

func (r *ExpenseRepo) ListByVehicle(ctx context.Context, vehicleID int64) ([]domain.Expense, error) {
	const query = `
			SELECT
				id, vehicle_id, category_id, amount, odometer_at,
				date, description, created_at
			FROM expenses
			WHERE vehicle_id = ?
			ORDER BY date DESC`

	var expenses []domain.Expense
	if err := r.db.SelectContext(ctx, &expenses, query, vehicleID); err != nil {
		return nil, fmt.Errorf("%w: list expenses: %w", domain.ErrInfrastructure, err)
	}

	return expenses, nil
}

func (r *ExpenseRepo) ListByVehicleAndPeriod(
	ctx context.Context,
	vehicleID int64,
	from, to time.Time,
) ([]domain.Expense, error) {
	const query = `
		SELECT
			id, vehicle_id, category_id, amount, odometer_at,
			date, description, created_at
		FROM expenses
		WHERE vehicle_id = ?
			AND date >= ?
			AND date <= ?
		ORDER BY date DESC`

	var expenses []domain.Expense
	if err := r.db.SelectContext(ctx, &expenses, query, vehicleID, from, to); err != nil {
		return nil, fmt.Errorf("%w: list expenses by period: %w", domain.ErrInfrastructure, err)
	}

	return expenses, nil
}

func (r *ExpenseRepo) SumByVehicle(ctx context.Context, vehicleID int64) (int64, error) {
	const query = `
			SELECT COALESCE(SUM(amount), 0)
			FROM expenses
			WHERE vehicle_id = ?`

	var total int64
	if err := r.db.GetContext(ctx, &total, query, vehicleID); err != nil {
		return 0, fmt.Errorf("%w: sum expenses: %w", domain.ErrInfrastructure, err)
	}

	return total, nil
}

func (r *ExpenseRepo) TotalsByVehicleCategory(ctx context.Context, vehicleID int64) ([]domain.ExpenseCategoryTotal, error) {
	const query = `
			SELECT
				c.id AS category_id,
				c.name AS category_name,
				SUM(e.amount) AS total_amount
			FROM expenses e
			JOIN expense_categories c ON c.id = e.category_id
			WHERE e.vehicle_id = ?
			GROUP BY c.id, c.name
			ORDER BY total_amount DESC`

	var totals []domain.ExpenseCategoryTotal
	if err := r.db.SelectContext(ctx, &totals, query, vehicleID); err != nil {
		return nil, fmt.Errorf("%w: list expense totals by category: %w", domain.ErrInfrastructure, err)
	}

	return totals, nil
}
