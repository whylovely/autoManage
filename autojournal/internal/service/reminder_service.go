package service

import (
	"autojournal/internal/domain"
	"context"
	"fmt"
	"strings"
	"time"
)

type ReminderService struct {
	repo        domain.ReminderRepository
	vehicleRepo domain.VehicleRepository
}

func NewReminderService(repo domain.ReminderRepository, vehicleRepo domain.VehicleRepository) *ReminderService {
	return &ReminderService{repo: repo, vehicleRepo: vehicleRepo}
}

func isValidReminderType(rT domain.ReminderType) bool {
	switch rT {
	case domain.ReminderTypeCustom,
		domain.ReminderTypeInsurance,
		domain.ReminderTypeOilChange,
		domain.ReminderTypeTireRotation:
		return true
	default:
		return false
	}
}

func validateReminder(reminder *domain.Reminder) error {
	if reminder == nil {
		return fmt.Errorf("%w: reminder is empty", domain.ErrValidation)
	}
	if reminder.VehicleID <= 0 {
		return fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}
	reminder.Title = strings.TrimSpace(reminder.Title)
	if reminder.Title == "" {
		return fmt.Errorf("%w: title is empty", domain.ErrValidation)
	}
	if !isValidReminderType(reminder.ReminderType) {
		return fmt.Errorf("%w: invalid reminder type", domain.ErrValidation)
	}
	if reminder.IntervalDays == nil && reminder.IntervalKM == nil {
		return fmt.Errorf("%w: at least one interval is required", domain.ErrValidation)
	}
	if reminder.IntervalDays != nil && *reminder.IntervalDays <= 0 {
		return fmt.Errorf("%w: interval days must be positive", domain.ErrValidation)
	}
	if reminder.IntervalKM != nil && *reminder.IntervalKM <= 0 {
		return fmt.Errorf("%w: interval km must be positive", domain.ErrValidation)
	}
	if reminder.IntervalDays != nil && reminder.LastDoneDate == nil {
		return fmt.Errorf("%w: interval days and last done date are null", domain.ErrValidation)
	}
	if reminder.LastDoneDate != nil && reminder.LastDoneDate.IsZero() {
		return fmt.Errorf("%w: last done date is required", domain.ErrValidation)
	}
	if reminder.IntervalKM != nil && reminder.LastDoneOdometer == nil {
		return fmt.Errorf("%w: intervals km and last done odometer are nulls", domain.ErrValidation)
	}
	if reminder.LastDoneOdometer != nil && *reminder.LastDoneOdometer < 0 {
		return fmt.Errorf("%w: last done odometer cannot be negative", domain.ErrValidation)
	}
	return nil
}

func (s *ReminderService) CreateReminder(ctx context.Context, reminder *domain.Reminder) error {
	if err := validateReminder(reminder); err != nil {
		return err
	}
	if _, err := s.vehicleRepo.GetByID(ctx, reminder.VehicleID); err != nil {
		return fmt.Errorf("get vehicle: %w", err)
	}
	calculateNextDue(reminder)

	reminder.IsActive = true
	return s.repo.Create(ctx, reminder)
}

func calculateNextDue(reminder *domain.Reminder) {
	reminder.NextDueDate = nil
	if reminder.IntervalDays != nil {
		nextDate := reminder.LastDoneDate.AddDate(0, 0, int(*reminder.IntervalDays))
		reminder.NextDueDate = &nextDate
	}

	reminder.NextDueOdometer = nil
	if reminder.IntervalKM != nil {
		nextOdometer := *reminder.LastDoneOdometer + *reminder.IntervalKM
		reminder.NextDueOdometer = &nextOdometer
	}
}

func (s *ReminderService) UpdateReminder(ctx context.Context, reminder *domain.Reminder) error {
	if err := validateReminder(reminder); err != nil {
		return err
	}
	if reminder.ID <= 0 {
		return fmt.Errorf("%w: reminder id must be positive", domain.ErrValidation)
	}
	if _, err := s.vehicleRepo.GetByID(ctx, reminder.VehicleID); err != nil {
		return fmt.Errorf("get vehicle: %w", err)
	}
	calculateNextDue(reminder)

	return s.repo.Update(ctx, reminder)
}

func (s *ReminderService) DeleteReminder(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: reminder id must be positive", domain.ErrValidation)
	}

	return s.repo.Delete(ctx, id)
}

func (s *ReminderService) ListVehicleReminders(ctx context.Context, vehicleID int64) ([]domain.Reminder, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("%w: vehicle id must be positive", domain.ErrValidation)
	}
	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		return nil, fmt.Errorf("get vehicle: %w", err)
	}

	reminders, err := s.repo.ListByVehicle(ctx, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("list vehicle reminders: %w", err)
	}

	return reminders, nil
}

func (s *ReminderService) GetDueReminders(ctx context.Context, now time.Time) ([]domain.DueReminder, error) {
	reminders, err := s.repo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active reminders: %w", err)
	}

	dateThreshold := now.AddDate(0, 0, 7)
	dueReminders := make([]domain.DueReminder, 0)

	for _, reminder := range reminders {
		dueByDate := reminder.NextDueDate != nil && !reminder.NextDueDate.After(dateThreshold)
		dueByOdometer := false

		if reminder.NextDueOdometer != nil {
			vehicle, err := s.vehicleRepo.GetByID(ctx, reminder.VehicleID)
			if err != nil {
				return nil, fmt.Errorf("get vehicle for reminder %d: %w", reminder.ID, err)
			}

			dueByOdometer = *reminder.NextDueOdometer <= vehicle.Odometer+500
		}

		if dueByDate || dueByOdometer {
			dueReminders = append(dueReminders, domain.DueReminder{
				Reminder:      reminder,
				DueByDate:     dueByDate,
				DueByOdometer: dueByOdometer,
			})
		}
	}

	return dueReminders, nil
}
