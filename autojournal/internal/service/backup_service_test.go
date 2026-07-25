package service

import (
	"autojournal/internal/domain"
	"context"
	"os"
	"testing"
	"time"
)

type backupRepoStub struct {
	created *domain.Backup
}

func (r *backupRepoStub) Create(_ context.Context, backup *domain.Backup) error {
	backup.ID = 1
	backup.CreatedAt = time.Now()
	r.created = backup
	return nil
}

func (backupRepoStub) GetByID(context.Context, int64) (*domain.Backup, error) { return nil, nil }
func (backupRepoStub) List(context.Context) ([]domain.Backup, error)          { return nil, nil }

type databaseBackuperStub struct {
	destination string
}

func (b *databaseBackuperStub) CreateSnapshot(_ context.Context, destination string) error {
	b.destination = destination
	return os.WriteFile(destination, []byte("backup"), 0o600)
}

func TestBackupService_CreateBackup_CreatesMetadata(t *testing.T) {
	repo := &backupRepoStub{}
	backuper := &databaseBackuperStub{}
	service := NewBackupService(repo, backuper, t.TempDir())

	backup, err := service.CreateBackup(context.Background(), "  before upgrade  ")
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}
	if backup.ID != 1 || backup.Note != "before upgrade" {
		t.Fatalf("unexpected backup metadata: %#v", backup)
	}
	if repo.created != backup || backuper.destination != backup.FilePath {
		t.Fatal("backup file and metadata must refer to the same backup")
	}
	if _, err := os.Stat(backup.FilePath); err != nil {
		t.Fatalf("backup file was not created: %v", err)
	}
}
