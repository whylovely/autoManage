package domain

import "time"

type Expense struct {
	ID          int64     `db:"id" json:"id"`
	VehicleID   int64     `db:"vehicle_id" json:"vehicleId"`
	CategoryID  int64     `db:"category_id" json:"categoryId"`
	Amount      int64     `db:"amount" json:"amount"`
	OdometerAt  int64     `db:"odometer_at" json:"odometerAt"`
	Date        time.Time `db:"date" json:"date"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type ExpenseCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type ExpenseCategoryTotal struct {
	CategoryID   int64  `db:"category_id" json:"categoryId"`
	CategoryName string `db:"category_name" json:"categoryName"`
	TotalAmount  int64  `db:"total_amount" json:"totalAmount"`
}

type MonthlyExpenseTotal struct {
	Month       string `json:"month"`
	TotalAmount int64  `json:"totalAmount"`
}

type ExpenseStats struct {
	TotalAmount int64                  `json:"totalAmount"`
	ByCategory  []ExpenseCategoryTotal `json:"byCategory"`
	ByMonth     []MonthlyExpenseTotal  `json:"byMonth"`
}
