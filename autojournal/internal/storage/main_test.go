package storage

import (
	"autojournal/migrations"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var integrationDB *sqlx.DB

func TestMain(m *testing.M) {
	db, err := sqlx.Open("sqlite3", "file:autojournal-integration?mode=memory&cache=shared&_foreign_keys=on")
	if err != nil {
		fmt.Fprintf(os.Stderr, "open integration database: %v\n", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "ping integration database: %v\n", err)
		os.Exit(1)
	}
	if err := migrations.RunMigrations(db); err != nil {
		fmt.Fprintf(os.Stderr, "run integration migrations: %v\n", err)
		os.Exit(1)
	}

	integrationDB = db
	code := m.Run()
	if err := db.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "close integration database: %v\n", err)
		if code == 0 {
			code = 1
		}
	}
	os.Exit(code)
}
