package domain

import "context"

// BackupRepository stores metadata for database backup files.
type BackupRepository interface {
	Create(ctx context.Context, backup *Backup) error
	GetByID(ctx context.Context, id int64) (*Backup, error)
	List(ctx context.Context) ([]Backup, error)
}
