package service

import (
	"autojournal/internal/domain"
	"context"
	"errors"
	"testing"
	"time"
)

func TestCalculateNextDue(t *testing.T) {
	lastDate := time.Date(2026, time.January, 10, 0, 0, 0, 0, time.UTC)
	nextDate := lastDate.AddDate(0, 0, 30)
	lastOdometer := int64(10_000)
	nextOdometer := int64(15_000)

	tests := []struct {
		name     string
		reminder domain.Reminder
		wantDate *time.Time
		wantKM   *int64
	}{
		{
			name: "date only",
			reminder: domain.Reminder{
				IntervalDays:    int64Pointer(30),
				LastDoneDate:    &lastDate,
				NextDueOdometer: int64Pointer(99_999),
			},
			wantDate: &nextDate,
			wantKM:   nil,
		},
		{
			name: "odometer only",
			reminder: domain.Reminder{
				IntervalKM:       int64Pointer(5_000),
				LastDoneOdometer: &lastOdometer,
				NextDueDate:      &lastDate,
			},
			wantDate: nil,
			wantKM:   &nextOdometer,
		},
		{
			name: "date and odometer",
			reminder: domain.Reminder{
				IntervalDays:     int64Pointer(30),
				LastDoneDate:     &lastDate,
				IntervalKM:       int64Pointer(5_000),
				LastDoneOdometer: &lastOdometer,
			},
			wantDate: &nextDate,
			wantKM:   &nextOdometer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculateNextDue(&tt.reminder)

			if !sameOptionalTime(tt.reminder.NextDueDate, tt.wantDate) {
				t.Fatalf("next due date = %v, want %v", tt.reminder.NextDueDate, tt.wantDate)
			}
			if !sameOptionalInt64(tt.reminder.NextDueOdometer, tt.wantKM) {
				t.Fatalf("next due odometer = %v, want %v", tt.reminder.NextDueOdometer, tt.wantKM)
			}
		})
	}
}

func TestValidateReminder_RejectsInvalidIntervals(t *testing.T) {
	lastDate := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	lastOdometer := int64(1_000)

	tests := []struct {
		name     string
		reminder domain.Reminder
	}{
		{
			name: "no intervals",
			reminder: domain.Reminder{
				VehicleID:    1,
				Title:        "Oil change",
				ReminderType: domain.ReminderTypeOilChange,
			},
		},
		{
			name: "zero day interval",
			reminder: domain.Reminder{
				VehicleID:    1,
				Title:        "Insurance",
				ReminderType: domain.ReminderTypeInsurance,
				IntervalDays: int64Pointer(0),
				LastDoneDate: &lastDate,
			},
		},
		{
			name: "zero kilometer interval",
			reminder: domain.Reminder{
				VehicleID:        1,
				Title:            "Tires",
				ReminderType:     domain.ReminderTypeTireRotation,
				IntervalKM:       int64Pointer(0),
				LastDoneOdometer: &lastOdometer,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReminder(&tt.reminder)
			if !errors.Is(err, domain.ErrValidation) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
}

func TestReminderService_GetDueReminders(t *testing.T) {
	now := time.Date(2026, time.January, 10, 12, 0, 0, 0, time.UTC)
	dateAtThreshold := now.AddDate(0, 0, 7)
	futureDate := now.AddDate(0, 0, 8)
	nextOdometer := int64(10_500)

	reminderRepo := &dueReminderRepoStub{reminders: []domain.Reminder{
		{
			ID:          1,
			VehicleID:   1,
			Title:       "Insurance",
			NextDueDate: &dateAtThreshold,
			IsActive:    true,
		},
		{
			ID:              2,
			VehicleID:       2,
			Title:           "Oil change",
			NextDueOdometer: &nextOdometer,
			IsActive:        true,
		},
		{
			ID:          3,
			VehicleID:   3,
			Title:       "Future reminder",
			NextDueDate: &futureDate,
			IsActive:    true,
		},
	}}
	vehicleRepo := &dueVehicleRepoStub{vehicles: map[int64]*domain.Vehicle{
		2: {ID: 2, Odometer: 10_000},
	}}
	service := NewReminderService(reminderRepo, vehicleRepo)

	dueReminders, err := service.GetDueReminders(context.Background(), now)
	if err != nil {
		t.Fatalf("get due reminders: %v", err)
	}
	if len(dueReminders) != 2 {
		t.Fatalf("got %d due reminders, want 2", len(dueReminders))
	}
	if dueReminders[0].Reminder.ID != 1 || !dueReminders[0].DueByDate || dueReminders[0].DueByOdometer {
		t.Fatalf("unexpected date due reminder: %#v", dueReminders[0])
	}
	if dueReminders[1].Reminder.ID != 2 || dueReminders[1].DueByDate || !dueReminders[1].DueByOdometer {
		t.Fatalf("unexpected odometer due reminder: %#v", dueReminders[1])
	}
}

type dueReminderRepoStub struct {
	reminders []domain.Reminder
}

func (*dueReminderRepoStub) Create(context.Context, *domain.Reminder) error { return nil }
func (*dueReminderRepoStub) Update(context.Context, *domain.Reminder) error { return nil }
func (*dueReminderRepoStub) Delete(context.Context, int64) error            { return nil }
func (*dueReminderRepoStub) GetByID(context.Context, int64) (*domain.Reminder, error) {
	return nil, domain.ErrNotFound
}
func (r *dueReminderRepoStub) ListByVehicle(context.Context, int64) ([]domain.Reminder, error) {
	return r.reminders, nil
}
func (r *dueReminderRepoStub) ListActive(context.Context) ([]domain.Reminder, error) {
	return r.reminders, nil
}

type dueVehicleRepoStub struct {
	vehicles map[int64]*domain.Vehicle
}

func (*dueVehicleRepoStub) Create(context.Context, *domain.Vehicle) error { return nil }
func (*dueVehicleRepoStub) Update(context.Context, *domain.Vehicle) error { return nil }
func (*dueVehicleRepoStub) Delete(context.Context, int64) error           { return nil }
func (r *dueVehicleRepoStub) GetByID(_ context.Context, id int64) (*domain.Vehicle, error) {
	vehicle, ok := r.vehicles[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return vehicle, nil
}
func (*dueVehicleRepoStub) List(context.Context) ([]domain.Vehicle, error) { return nil, nil }

func int64Pointer(value int64) *int64 {
	return &value
}

func sameOptionalInt64(actual, expected *int64) bool {
	if actual == nil || expected == nil {
		return actual == expected
	}
	return *actual == *expected
}

func sameOptionalTime(actual, expected *time.Time) bool {
	if actual == nil || expected == nil {
		return actual == expected
	}
	return actual.Equal(*expected)
}
