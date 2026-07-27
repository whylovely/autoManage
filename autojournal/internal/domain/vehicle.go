package domain

import "time"

type Vehicle struct {
	ID           int64     `db:"id" json:"id"`
	VIN          string    `db:"vin" json:"vin" validate:"required,len=17"`
	Make         string    `db:"make" json:"make" validate:"required"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	Color        string    `db:"color" json:"color,omitempty"`
	EngineVolume int64     `json:"engineVolume"`
	FuelType     int       `json:"fuelType"`
	Odometer     int64     `json:"odometer"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdateAt     time.Time `json:"updatedAt"`
}
