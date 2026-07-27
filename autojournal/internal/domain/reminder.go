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
	ID               int64        `db:"id" json:"id"`
	VehicleID        int64        `db:"vehicle_id" json:"vehicleId"`
	Title            string       `db:"title" json:"title"`
	ReminderType     ReminderType `db:"reminder_type" json:"reminderType"`
	IntervalKM       *int64       `db:"interval_km" json:"intervalKm"`
	IntervalDays     *int64       `db:"interval_days" json:"intervalDays"`
	LastDoneOdometer *int64       `db:"last_done_odometer" json:"lastDoneOdometer"`
	LastDoneDate     *time.Time   `db:"last_done_date" json:"lastDoneDate"`
	NextDueDate      *time.Time   `db:"next_due_date" json:"nextDueDate"`
	NextDueOdometer  *int64       `db:"next_due_odometer" json:"nextDueOdometer"`
	IsActive         bool         `db:"is_active" json:"isActive"`
	CreatedAt        time.Time    `db:"created_at" json:"createdAt"`
}

type DueReminder struct {
	Reminder      Reminder `json:"reminder"`
	DueByDate     bool     `json:"dueByDate"`
	DueByOdometer bool     `json:"dueByOdometer"`
}
