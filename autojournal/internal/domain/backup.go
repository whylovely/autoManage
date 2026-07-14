package domain

import "time"

// Backup stores metadata about a created database backup file.
// The actual backup creation logic belongs in service/storage layers.
type Backup struct {
	ID        int64
	FilePath  string
	Note      string
	CreatedAt time.Time
}
