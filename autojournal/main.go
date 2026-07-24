package main

import (
	"autojournal/internal/handler"
	"autojournal/internal/service"
	"autojournal/internal/storage"
	"autojournal/migrations"
	"embed"
	"log"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	db, err := storage.OpenSQLite()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migrations.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	vehicleRepo := storage.NewVehicleRepo(db)
	expenseRepo := storage.NewExpenseRepo(db)
	categoryRepo := storage.NewExpenseCategoryRepo(db)
	backupRepo := storage.NewBackupRepo(db)

	vehicleService := service.NewVehicleService(vehicleRepo)
	expenseService := service.NewExpenseService(expenseRepo, vehicleRepo, categoryRepo)
	categoryService := service.NewExpenseCategory(categoryRepo)

	appDataDir, err := storage.AppDataDir()
	if err != nil {
		log.Fatal(err)
	}
	backupService := service.NewBackupService(
		backupRepo,
		storage.NewSQLiteBackuper(db),
		filepath.Join(appDataDir, "backups"),
	)

	app := handler.NewApp(
		vehicleService,
		expenseService,
		categoryService,
		backupService,
	)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "autojournal",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
