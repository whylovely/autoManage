package storage

import (
	"autojournal/internal/domain"
	"autojournal/migrations"
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestReminderRepo_CRUD(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "autojournal.db")
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
	if err := migrations.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	vehicleID := createReminderTestVehicle(t, db)
	repo := NewReminderRepo(db)

	intervalKM := int64(10_000)
	lastOdometer := int64(50_000)
	lastDate := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	nextDate := lastDate.AddDate(0, 6, 0)
	nextOdometer := int64(60_000)
	reminder := &domain.Reminder{
		VehicleID:        vehicleID,
		Title:            "Insurance renewal",
		ReminderType:     domain.ReminderTypeInsurance,
		IntervalKM:       &intervalKM,
		LastDoneOdometer: &lastOdometer,
		LastDoneDate:     &lastDate,
		NextDueDate:      &nextDate,
		NextDueOdometer:  &nextOdometer,
		IsActive:         true,
	}

	ctx := context.Background()
	if err := repo.Create(ctx, reminder); err != nil {
		t.Fatalf("create reminder: %v", err)
	}
	if reminder.ID == 0 || reminder.CreatedAt.IsZero() {
		t.Fatalf("created reminder is incomplete: %#v", reminder)
	}

	active, err := repo.ListActive(ctx)
	if err != nil {
		t.Fatalf("list active reminders: %v", err)
	}
	if len(active) != 1 || active[0].NextDueOdometer == nil || *active[0].NextDueOdometer != nextOdometer {
		t.Fatalf("unexpected active reminders: %#v", active)
	}

	reminder.Title = "Insurance renewed"
	reminder.IsActive = false
	if err := repo.Update(ctx, reminder); err != nil {
		t.Fatalf("update reminder: %v", err)
	}

	updated, err := repo.GetByID(ctx, reminder.ID)
	if err != nil {
		t.Fatalf("get reminder: %v", err)
	}
	if updated.Title != reminder.Title || updated.IsActive {
		t.Fatalf("reminder update was not persisted: %#v", updated)
	}

	active, err = repo.ListActive(ctx)
	if err != nil {
		t.Fatalf("list active reminders after update: %v", err)
	}
	if len(active) != 0 {
		t.Fatalf("inactive reminder must not be listed as active: %#v", active)
	}

	if err := repo.Delete(ctx, reminder.ID); err != nil {
		t.Fatalf("delete reminder: %v", err)
	}
	if _, err := repo.GetByID(ctx, reminder.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func createReminderTestVehicle(t *testing.T, db *sqlx.DB) int64 {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO vehicles (vin, make, model, year)
		VALUES (?, ?, ?, ?)`,
		"1HGCM82633A004352",
		"Honda",
		"Accord",
		2020,
	)
	if err != nil {
		t.Fatalf("create test vehicle: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get test vehicle id: %v", err)
	}
	return id
}
