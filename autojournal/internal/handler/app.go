package handler

import (
	"autojournal/internal/domain"
	"autojournal/internal/scheduler"
	"autojournal/internal/service"
	"context"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context

	vehicleService    *service.VehicleService
	expenseService    *service.ExpenseService
	categoryService   *service.ExpenseCategoryService
	backupService     *service.BackupService
	reminderService   *service.ReminderService
	reminderScheduler *scheduler.ReminderScheduler
}

func NewApp(
	vehicleService *service.VehicleService,
	expenseService *service.ExpenseService,
	categoryService *service.ExpenseCategoryService,
	backupService *service.BackupService,
	reminderService *service.ReminderService,
) *App {
	return &App{
		vehicleService:  vehicleService,
		expenseService:  expenseService,
		categoryService: categoryService,
		backupService:   backupService,
		reminderService: reminderService,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.reminderScheduler = scheduler.NewReminderScheduler(
		a.reminderService,
		func(ctx context.Context, event string, payload any) {
			runtime.EventsEmit(ctx, event, payload)
		},
	)

	if err := a.reminderScheduler.Start(ctx); err != nil {
		log.Printf("start reminder scheduler: %v", err)
	}
}

func (a *App) CreateVehicle(vehicle domain.Vehicle) (*domain.Vehicle, error) {
	if err := a.vehicleService.CreateVehicle(a.ctx, &vehicle); err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (a *App) ListVehicles() ([]domain.Vehicle, error) {
	return a.vehicleService.ListVehicles(a.ctx)
}

func (a *App) GetVehicles() ([]domain.Vehicle, error) {
	return a.vehicleService.ListVehicles(a.ctx)
}

func (a *App) GetVehicle(id int64) (*domain.Vehicle, error) {
	return a.vehicleService.GetVehicle(a.ctx, id)
}

func (a *App) UpdateVehicle(vehicle domain.Vehicle) (*domain.Vehicle, error) {
	if err := a.vehicleService.UpdateVehicle(a.ctx, &vehicle); err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (a *App) DeleteVehicle(id int64) error {
	return a.vehicleService.DeleteVehicle(a.ctx, id)
}

func (a *App) UpdateOdometer(vehicleID, odometer int64) error {
	return a.vehicleService.UpdateOdometer(a.ctx, vehicleID, odometer)
}

func (a *App) AddExpense(expense domain.Expense) (*domain.Expense, error) {
	if err := a.expenseService.AddExpense(a.ctx, &expense); err != nil {
		return nil, err
	}

	return &expense, nil
}

func (a *App) ListVehicleExpenses(vehicleID int64) ([]domain.Expense, error) {
	return a.expenseService.ListVehicleExpenses(a.ctx, vehicleID)
}

func (a *App) DeleteExpense(id int64) error {
	return a.expenseService.DeleteExpense(a.ctx, id)
}

func (a *App) GetExpenseStats(vehicleID int64) (*domain.ExpenseStats, error) {
	return a.expenseService.GetExpenseStats(a.ctx, vehicleID)
}

func (a *App) ListExpenseCategories() ([]domain.ExpenseCategory, error) {
	return a.categoryService.ListCategories(a.ctx)
}

func (a *App) ExportVehicleExpenses(vehicleID int64, format string) (string, error) {
	extension := format
	if extension != "csv" && extension != "json" {
		extension = "json"
	}

	destination, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Экспорт расходов",
		DefaultFilename: "expenses." + extension,
		Filters: []runtime.FileFilter{
			{DisplayName: "Export (*." + extension + ")", Pattern: "*." + extension},
		},
	})
	if err != nil {
		return "", err
	}
	if destination == "" {
		return "", nil
	}

	if err := a.expenseService.ExportVehicleExpenses(a.ctx, vehicleID, extension, destination); err != nil {
		return "", err
	}

	return destination, nil
}

func (a *App) CreateReminder(reminder domain.Reminder) (*domain.Reminder, error) {
	if err := a.reminderService.CreateReminder(a.ctx, &reminder); err != nil {
		return nil, err
	}

	return &reminder, nil
}

func (a *App) UpdateReminder(reminder domain.Reminder) (*domain.Reminder, error) {
	if err := a.reminderService.UpdateReminder(a.ctx, &reminder); err != nil {
		return nil, err
	}

	return &reminder, nil
}

func (a *App) DeleteReminder(id int64) error {
	return a.reminderService.DeleteReminder(a.ctx, id)
}

func (a *App) ListVehicleReminders(vehicleID int64) ([]domain.Reminder, error) {
	return a.reminderService.ListVehicleReminders(a.ctx, vehicleID)
}

func (a *App) GetDueReminders() ([]domain.DueReminder, error) {
	return a.reminderService.GetDueReminders(a.ctx, time.Now())
}

func (a *App) CreateBackup(note string) (*domain.Backup, error) {
	return a.backupService.CreateBackup(a.ctx, note)
}

func (a *App) ListBackups() ([]domain.Backup, error) {
	return a.backupService.ListBackups(a.ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	if a.reminderScheduler != nil {
		a.reminderScheduler.Stop()
	}
}
