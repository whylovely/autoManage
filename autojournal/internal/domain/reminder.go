package domain

import "time"

type ReminderType string

const (
	ReminderTypeOilChange    ReminderType = "oil_change"
	ReminderTypeTireRotation ReminderType = "tire_rotation"
	ReminderTypeInsurance    ReminderType = "insurance"
	ReminderTypeCustom       ReminderType = "custom"
)

type Reminder struct {
	ID               int64        `db:"id"`
	VehicleID        int64        `db:"vehicle_id"`
	Title            string       `db:"title"`
	ReminderType     ReminderType `db:"reminder_type"`
	IntervalKM       *int64       `db:"interval_km"`
	IntervalDays     *int64       `db:"interval_days"`
	LastDoneOdometer *int64       `db:"last_done_odometer"`
	LastDoneDate     *time.Time   `db:"last_done_date"`
	NextDueDate      *time.Time   `db:"next_due_date"`
	NextDueOdometer  *int64       `db:"next_due_odometer"`
	IsActive         bool         `db:"is_active"`
	CreatedAt        time.Time    `db:"created_at"`
}
