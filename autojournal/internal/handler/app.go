package handler

import (
	"autojournal/internal/service"
	"context"
	"fmt"
)

type App struct {
	ctx context.Context

	vehicleService  *service.VehicleService
	expenseService  *service.ExpenseService
	categoryService *service.ExpenseCategorySevice
}

func NewApp(
	vehicleService *service.VehicleService,
	expenseService *service.ExpenseService,
	categoryService *service.ExpenseCategorySevice,
) *App {
	return &App{
		vehicleService:  vehicleService,
		expenseService:  expenseService,
		categoryService: categoryService,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
