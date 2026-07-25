package domain

import "context"

type ReminderRepository interface {
	Create(ctx context.Context, reminder *Reminder) error
	Update(ctx context.Context, reminder *Reminder) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Reminder, error)

	ListByVehicle(ctx context.Context, vehicleID int64) ([]Reminder, error)
	ListActive(ctx context.Context) ([]Reminder, error)
}
