package storage

import (
	"autojournal/internal/domain"
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// SQLiteBackuper creates consistent SQLite snapshots without manually copying
// WAL files.
type SQLiteBackuper struct {
	db *sqlx.DB
}

func NewSQLiteBackuper(db *sqlx.DB) *SQLiteBackuper {
	return &SQLiteBackuper{db: db}
}

func (b *SQLiteBackuper) CreateSnapshot(ctx context.Context, destination string) error {
	if strings.TrimSpace(destination) == "" {
		return fmt.Errorf("%w: backup destination is required", domain.ErrValidation)
	}

	if _, err := b.db.ExecContext(ctx, `VACUUM INTO ?`, destination); err != nil {
		return fmt.Errorf("%w: create sqlite snapshot: %w", domain.ErrInfrastructure, err)
	}

	return nil
}
