package domain

import "context"

// DatabaseBackuper creates a consistent snapshot of the application database.
type DatabaseBackuper interface {
	CreateSnapshot(ctx context.Context, destination string) error
}
