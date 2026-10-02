package domain

import "context"

type DatabaseBackuper interface {
	CreateSnapshot(ctx context.Context, destination string) error
}
