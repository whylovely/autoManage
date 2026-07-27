package domain

import "time"

// Backup stores metadata about a created database backup file.
// The actual backup creation logic belongs in service/storage layers.
type Backup struct {
	ID        int64     `json:"id"`
	FilePath  string    `json:"filePath"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
}
