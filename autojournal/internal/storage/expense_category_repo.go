package storage

import (
	"autojournal/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ExpenseCategoryRepo struct {
	db *sqlx.DB
}

func NewExpenseCategoryRepo(db *sqlx.DB) *ExpenseCategoryRepo {
	return &ExpenseCategoryRepo{db: db}
}

func (r *ExpenseCategoryRepo) List(ctx context.Context) ([]domain.ExpenseCategory, error) {
	const query = `
			SELECT id, name, icon
			FROM expense_categories
			ORDER BY name ASC`

	var categories []domain.ExpenseCategory
	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, fmt.Errorf("%w: list expense categories: %w", domain.ErrInfrastructure, err)
	}

	return categories, nil
}

func (r *ExpenseCategoryRepo) GetByID(ctx context.Context, id int64) (*domain.ExpenseCategory, error) {
	const query = `
			SELECT id, name, icon
			FROM expense_categories
			WHERE id = ?`

	var category domain.ExpenseCategory
	if err := r.db.GetContext(ctx, &category, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: category %d", domain.ErrNotFound, id)
		}
		return nil, fmt.Errorf("%w: get category: %w", domain.ErrInfrastructure, err)
	}

	return &category, nil
}
