package storage

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestSQLiteBackuper_CreateSnapshot(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "source.db")
	destinationPath := filepath.Join(dir, "backup.db")

	db, err := sqlx.Open("sqlite3", sourcePath)
	if err != nil {
		t.Fatalf("open source database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE records (id INTEGER PRIMARY KEY, name TEXT NOT NULL)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO records (name) VALUES ('first')`); err != nil {
		t.Fatalf("insert record: %v", err)
	}

	backuper := NewSQLiteBackuper(db)
	if err := backuper.CreateSnapshot(context.Background(), destinationPath); err != nil {
		t.Fatalf("create snapshot: %v", err)
	}

	backupDB, err := sqlx.Open("sqlite3", destinationPath)
	if err != nil {
		t.Fatalf("open backup database: %v", err)
	}
	defer backupDB.Close()

	var count int
	if err := backupDB.Get(&count, `SELECT COUNT(*) FROM records`); err != nil {
		t.Fatalf("count backup records: %v", err)
	}
	if count != 1 {
		t.Fatalf("backup contains %d records, want 1", count)
	}
}
