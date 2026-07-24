package domain

import "time"

type Expense struct {
	ID          int64     `db:"id"`
	VehicleID   int64     `db:"vehicle_id"`
	CategoryID  int64     `db:"category_id"`
	Amount      int64     `db:"amount"`
	OdometerAt  int64     `db:"odometer_at"`
	Date        time.Time `db:"date"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

type ExpenseCategory struct {
	ID   int64
	Name string
	Icon string
}

type ExpenseCategoryTotal struct {
	CategoryID   int64  `db:"category_id"`
	CategoryName string `db:"category_name"`
	TotalAmount  int64  `db:"total_amount"`
}
