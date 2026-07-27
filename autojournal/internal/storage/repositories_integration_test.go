package storage

import (
	"autojournal/internal/domain"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetIntegrationDatabase(t *testing.T) {
	t.Helper()
	_, err := integrationDB.Exec(`
		DELETE FROM reminders;
		DELETE FROM expenses;
		DELETE FROM backups;
		DELETE FROM vehicles;
		DELETE FROM sqlite_sequence WHERE name IN ('reminders', 'expenses', 'backups', 'vehicles');
	`)
	require.NoError(t, err)
}

func integrationVehicle(t *testing.T) *domain.Vehicle {
	t.Helper()
	vehicle := &domain.Vehicle{
		VIN: "JH4KA8260MC000001", Make: "Honda", Model: "Legend", Year: 2021,
		Color: "Black", EngineVolume: 3500, FuelType: 1, Odometer: 40_000,
	}
	require.NoError(t, NewVehicleRepo(integrationDB).Create(context.Background(), vehicle))
	return vehicle
}

func TestVehicleRepo_IntegrationCRUD(t *testing.T) {
	resetIntegrationDatabase(t)
	ctx := context.Background()
	repo := NewVehicleRepo(integrationDB)
	vehicle := integrationVehicle(t)

	assert.Positive(t, vehicle.ID)
	stored, err := repo.GetByID(ctx, vehicle.ID)
	require.NoError(t, err)
	assert.Equal(t, vehicle.VIN, stored.VIN)
	assert.False(t, stored.CreatedAt.IsZero())

	vehicle.Make = "Acura"
	vehicle.Odometer = 41_500
	require.NoError(t, repo.Update(ctx, vehicle))
	stored, err = repo.GetByID(ctx, vehicle.ID)
	require.NoError(t, err)
	assert.Equal(t, "Acura", stored.Make)
	assert.EqualValues(t, 41_500, stored.Odometer)

	vehicles, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, vehicles, 1)

	require.NoError(t, repo.Delete(ctx, vehicle.ID))
	_, err = repo.GetByID(ctx, vehicle.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.ErrorIs(t, repo.Delete(ctx, vehicle.ID), domain.ErrNotFound)
}

func TestExpenseAndCategoryRepos_IntegrationQueries(t *testing.T) {
	resetIntegrationDatabase(t)
	ctx := context.Background()
	vehicle := integrationVehicle(t)
	categoryRepo := NewExpenseCategoryRepo(integrationDB)
	categories, err := categoryRepo.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, categories)
	category, err := categoryRepo.GetByID(ctx, categories[0].ID)
	require.NoError(t, err)
	assert.Equal(t, categories[0].Name, category.Name)
	_, err = categoryRepo.GetByID(ctx, 999_999)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	repo := NewExpenseRepo(integrationDB)
	january := &domain.Expense{
		VehicleID: vehicle.ID, CategoryID: category.ID, Amount: 2_000,
		OdometerAt: 40_100, Date: time.Date(2026, time.January, 10, 12, 0, 0, 0, time.UTC), Description: "January",
	}
	february := &domain.Expense{
		VehicleID: vehicle.ID, CategoryID: category.ID, Amount: 3_000,
		OdometerAt: 41_000, Date: time.Date(2026, time.February, 10, 12, 0, 0, 0, time.UTC), Description: "February",
	}
	require.NoError(t, repo.Create(ctx, january))
	require.NoError(t, repo.Create(ctx, february))
	assert.Positive(t, january.ID)

	stored, err := repo.GetByID(ctx, january.ID)
	require.NoError(t, err)
	assert.Equal(t, january.Description, stored.Description)

	all, err := repo.ListByVehicle(ctx, vehicle.ID)
	require.NoError(t, err)
	require.Len(t, all, 2)
	assert.Equal(t, february.ID, all[0].ID)

	period, err := repo.ListByVehicleAndPeriod(
		ctx,
		vehicle.ID,
		time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, time.January, 31, 23, 59, 59, 0, time.UTC),
	)
	require.NoError(t, err)
	require.Len(t, period, 1)
	assert.Equal(t, january.ID, period[0].ID)

	total, err := repo.SumByVehicle(ctx, vehicle.ID)
	require.NoError(t, err)
	assert.EqualValues(t, 5_000, total)
	totals, err := repo.TotalsByVehicleCategory(ctx, vehicle.ID)
	require.NoError(t, err)
	require.Len(t, totals, 1)
	assert.EqualValues(t, 5_000, totals[0].TotalAmount)

	require.NoError(t, repo.Delete(ctx, january.ID))
	_, err = repo.GetByID(ctx, january.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.ErrorIs(t, repo.Delete(ctx, january.ID), domain.ErrNotFound)

	invalid := &domain.Expense{VehicleID: 999_999, CategoryID: category.ID, Amount: 1, Date: time.Now()}
	assert.True(t, errors.Is(repo.Create(ctx, invalid), domain.ErrInfrastructure))
}

func TestBackupRepo_IntegrationCRUD(t *testing.T) {
	resetIntegrationDatabase(t)
	ctx := context.Background()
	repo := NewBackupRepo(integrationDB)
	backup := &domain.Backup{FilePath: "/tmp/autojournal-backup.db", Note: "before update"}

	require.NoError(t, repo.Create(ctx, backup))
	assert.Positive(t, backup.ID)
	assert.False(t, backup.CreatedAt.IsZero())

	stored, err := repo.GetByID(ctx, backup.ID)
	require.NoError(t, err)
	assert.Equal(t, backup.FilePath, stored.FilePath)
	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 1)
	_, err = repo.GetByID(ctx, 999_999)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestVehicleRepo_DeleteCascadesRelatedData(t *testing.T) {
	resetIntegrationDatabase(t)
	ctx := context.Background()
	vehicle := integrationVehicle(t)
	categories, err := NewExpenseCategoryRepo(integrationDB).List(ctx)
	require.NoError(t, err)

	expense := &domain.Expense{
		VehicleID: vehicle.ID, CategoryID: categories[0].ID, Amount: 100,
		Date: time.Now(),
	}
	require.NoError(t, NewExpenseRepo(integrationDB).Create(ctx, expense))
	interval := int64(5_000)
	last := int64(40_000)
	reminder := &domain.Reminder{
		VehicleID: vehicle.ID, Title: "Oil", ReminderType: domain.ReminderTypeOilChange,
		IntervalKM: &interval, LastDoneOdometer: &last, IsActive: true,
	}
	require.NoError(t, NewReminderRepo(integrationDB).Create(ctx, reminder))

	require.NoError(t, NewVehicleRepo(integrationDB).Delete(ctx, vehicle.ID))
	var expenseCount, reminderCount int
	require.NoError(t, integrationDB.Get(&expenseCount, `SELECT COUNT(*) FROM expenses`))
	require.NoError(t, integrationDB.Get(&reminderCount, `SELECT COUNT(*) FROM reminders`))
	assert.Zero(t, expenseCount)
	assert.Zero(t, reminderCount)
}
