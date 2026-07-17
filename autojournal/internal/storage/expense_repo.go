package storage

import (
	"autojournal/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		return fmt.Errorf("create expense: %w", err)
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
		return fmt.Errorf("delete expense: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("expense %d not found", id)
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
			return nil, fmt.Errorf("expense %d not found", id)
		}
		return nil, fmt.Errorf("get expense: %w", err)
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
		return nil, fmt.Errorf("list expenses: %w", err)
	}

	return expenses, nil
}
