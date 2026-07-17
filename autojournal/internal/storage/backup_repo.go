package storage

import (
	"autojournal/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// BackupRepo persists metadata about backup files, not the files themselves.
type BackupRepo struct {
	db *sqlx.DB
}

func NewBackupRepo(db *sqlx.DB) *BackupRepo {
	return &BackupRepo{db: db}
}

func (r *BackupRepo) Create(ctx context.Context, backup *domain.Backup) error {
	const query = `
		INSERT INTO backups (file_path, note)
		VALUES (?, ?)
		RETURNING id, created_at
	`

	err := r.db.QueryRowxContext(ctx, query, backup.FilePath, backup.Note).Scan(
		&backup.ID,
		&backup.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create backup metadata: %w", err)
	}

	return nil
}

func (r *BackupRepo) GetByID(ctx context.Context, id int64) (*domain.Backup, error) {
	const query = `
		SELECT id, file_path AS filepath, note, created_at AS createdat
		FROM backups
		WHERE id = ?
	`

	var backup domain.Backup
	if err := r.db.GetContext(ctx, &backup, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("backup %d not found", id)
		}
		return nil, fmt.Errorf("get backup metadata: %w", err)
	}

	return &backup, nil
}

func (r *BackupRepo) List(ctx context.Context) ([]domain.Backup, error) {
	const query = `
		SELECT id, file_path AS filepath, note, created_at AS createdat
		FROM backups
		ORDER BY created_at DESC
	`

	var backups []domain.Backup
	if err := r.db.SelectContext(ctx, &backups, query); err != nil {
		return nil, fmt.Errorf("list backup metadata: %w", err)
	}

	return backups, nil
}
