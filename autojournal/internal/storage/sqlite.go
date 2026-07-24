package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

func AppDataDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}

	return filepath.Join(configDir, "autojournal"), nil
}

func OpenSQLite() (*sqlx.DB, error) {
	appDir, err := AppDataDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return nil, fmt.Errorf("create app dir: %w", err)
	}

	dbPath := filepath.Join(appDir, "data.db")

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		return nil, fmt.Errorf("enable wal: %w", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	return db, nil
}
