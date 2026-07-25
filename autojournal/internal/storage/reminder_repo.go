package storage

import (
	"autojournal/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ReminderRepo struct {
	db *sqlx.DB
}

func NewReminderRepo(db *sqlx.DB) *ReminderRepo {
	return &ReminderRepo{db: db}
}

func (r *ReminderRepo) Create(ctx context.Context, reminder *domain.Reminder) error {
	const query = `
			INSERT INTO reminders (
				vehicle_id, title, reminder_type,
				interval_km, interval_days,
				last_done_odometer, last_done_date,
				next_due_date, next_due_odometer,
				is_active)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		reminder.VehicleID,
		reminder.Title,
		reminder.ReminderType,
		reminder.IntervalKM,
		reminder.IntervalDays,
		reminder.LastDoneOdometer,
		reminder.LastDoneDate,
		reminder.NextDueDate,
		reminder.NextDueOdometer,
		reminder.IsActive,
	).Scan(&reminder.ID, &reminder.CreatedAt)
	if err != nil {
		return fmt.Errorf("%w: create reminder: %w", domain.ErrInfrastructure, err)
	}

	return nil
}

func (r *ReminderRepo) GetByID(ctx context.Context, id int64) (*domain.Reminder, error) {
	const query = `
			SELECT 
				id, vehicle_id, title, 
				reminder_type, interval_km, interval_days,
				last_done_odometer, last_done_date, next_due_date,
				next_due_odometer, is_active, created_at
			FROM reminders
			WHERE id = ?`

	var reminder domain.Reminder
	if err := r.db.GetContext(ctx, &reminder, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: reminder %d", domain.ErrNotFound, id)
		}
		return nil, fmt.Errorf("%w: get reminder: %w", domain.ErrInfrastructure, err)
	}

	return &reminder, nil
}

func (r *ReminderRepo) Update(ctx context.Context, reminder *domain.Reminder) error {
	const query = `
			UPDATE reminders
			SET vehicle_id = ?, title = ?, reminder_type = ?,
				interval_km = ?, interval_days = ?, 
				last_done_odometer = ?, last_done_date = ?,
				next_due_date = ?, next_due_odometer = ?,
				is_active = ?
			WHERE id = ?`

	result, err := r.db.ExecContext(
		ctx,
		query,
		reminder.VehicleID,
		reminder.Title,
		reminder.ReminderType,
		reminder.IntervalKM,
		reminder.IntervalDays,
		reminder.LastDoneOdometer,
		reminder.LastDoneDate,
		reminder.NextDueDate,
		reminder.NextDueOdometer,
		reminder.IsActive,
		reminder.ID,
	)
	if err != nil {
		return fmt.Errorf("%w: update reminder: %w", domain.ErrInfrastructure, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: get update rows: %w", domain.ErrInfrastructure, err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: reminder %d", domain.ErrNotFound, reminder.ID)
	}

	return nil
}

func (r *ReminderRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM reminders WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("%w: delete reminder: %w", domain.ErrInfrastructure, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: get deleted rows: %w", domain.ErrInfrastructure, err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: reminder %d", domain.ErrNotFound, id)
	}

	return nil
}

func (r *ReminderRepo) ListByVehicle(ctx context.Context, vehicleID int64) ([]domain.Reminder, error) {
	const query = `
			SELECT 
				id, vehicle_id, title, 
				reminder_type, interval_km, interval_days,
				last_done_odometer, last_done_date, next_due_date,
				next_due_odometer, is_active, created_at
			FROM reminders
			WHERE vehicle_id = ?
			ORDER BY is_active DESC, created_at DESC`

	var reminders []domain.Reminder
	if err := r.db.SelectContext(ctx, &reminders, query, vehicleID); err != nil {
		return nil, fmt.Errorf("%w: list reminders: %w", domain.ErrInfrastructure, err)
	}

	return reminders, nil
}

func (r *ReminderRepo) ListActive(ctx context.Context) ([]domain.Reminder, error) {
	const query = `
			SELECT 
				id, vehicle_id, title, 
				reminder_type, interval_km, interval_days,
				last_done_odometer, last_done_date, next_due_date,
				next_due_odometer, is_active, created_at
			FROM reminders
			WHERE is_active = 1
			ORDER BY created_at DESC`

	var reminders []domain.Reminder
	if err := r.db.SelectContext(ctx, &reminders, query); err != nil {
		return nil, fmt.Errorf("%w: list reminders: %w", domain.ErrInfrastructure, err)
	}

	return reminders, nil
}
