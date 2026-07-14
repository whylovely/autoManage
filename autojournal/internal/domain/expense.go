package domain

import "time"

type Expense struct {
	ID          int64
	VehicleID   int64
	CategoryID  int64
	Amount      int64 `db:"amount" json:"amount" validate:"required,gt=0"`
	OdometerAt  int64
	Date        time.Time
	Description string
	CreatedAt   time.Time
}

type ExpenseCategory struct {
	ID   int64
	Name string
	Icon string
}
