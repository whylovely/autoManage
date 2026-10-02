package domain

import "time"

type Backup struct {
	ID        int64     `json:"id"`
	FilePath  string    `json:"filePath"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
}
