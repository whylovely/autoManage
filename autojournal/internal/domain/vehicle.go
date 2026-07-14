package domain

import "time"

type Vehicle struct {
	ID           int64  `db:"id" json:"id"`
	VIN          string `db:"vin" json:"vin" validate:"required, len=17"`
	Make         string `db:"make" json:"make" validate:"required"`
	Model        string
	Year         int
	Color        string `db:"color" json:"color, omitempty"`
	EngineVolume int64
	FuelType     int
	Odometer     int64
	Notes        string
	CreatedAt    time.Time
	UpdateAt     time.Time
}
