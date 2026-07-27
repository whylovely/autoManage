package handler

import (
	"autojournal/internal/domain"
	"autojournal/internal/scheduler"
	"autojournal/internal/service"
	"context"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context

	vehicleService    *service.VehicleService
	expenseService    *service.ExpenseService
	categoryService   *service.ExpenseCategorySevice
	backupService     *service.BackupService
	reminderService   *service.ReminderService
	reminderScheduler *scheduler.ReminderScheduler
}

func NewApp(
	vehicleService *service.VehicleService,
	expenseService *service.ExpenseService,
	categoryService *service.ExpenseCategorySevice,
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

func (a *App) AddExpense(expense domain.Expense) (*domain.Expense, error) {
	if err := a.expenseService.AddExpense(a.ctx, &expense); err != nil {
		return nil, err
	}

	return &expense, nil
}

func (a *App) ListVehicleExpenses(vehicleID int64) ([]domain.Expense, error) {
	return a.expenseService.ListVehicleExpenses(a.ctx, vehicleID)
}

func (a *App) CreateBackup(note string) (*domain.Backup, error) {
	return a.backupService.CreateBackup(a.ctx, note)
}

func (a *App) ListBackups() ([]domain.Backup, error) {
	return a.backupService.ListBackups(a.ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Shutdown(ctx context.Context) {
	if a.reminderScheduler != nil {
		a.reminderScheduler.Stop()
	}
}
