package service

import (
	"autojournal/internal/domain"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BackupService struct {
	repo       domain.BackupRepository
	dbBackuper domain.DatabaseBackuper
	backupDir  string
}

func NewBackupService(
	repo domain.BackupRepository,
	dbBackuper domain.DatabaseBackuper,
	backupDir string,
) *BackupService {
	return &BackupService{
		repo:       repo,
		dbBackuper: dbBackuper,
		backupDir:  backupDir,
	}
}

func (s *BackupService) CreateBackup(ctx context.Context, note string) (*domain.Backup, error) {
	if strings.TrimSpace(s.backupDir) == "" {
		return nil, fmt.Errorf("backup directory is required")
	}

	if err := os.MkdirAll(s.backupDir, 0o755); err != nil {
		return nil, fmt.Errorf("create backup directory: %w", err)
	}

	fileName := fmt.Sprintf(
		"autojournal-%s.db",
		time.Now().UTC().Format("20060102-150405.000000000"),
	)
	filePath := filepath.Join(s.backupDir, fileName)

	if err := s.dbBackuper.CreateSnapshot(ctx, filePath); err != nil {
		return nil, fmt.Errorf("create backup file: %w", err)
	}

	backup := &domain.Backup{
		FilePath: filePath,
		Note:     strings.TrimSpace(note),
	}
	if err := s.repo.Create(ctx, backup); err != nil {
		return nil, fmt.Errorf("save backup metadata: %w", err)
	}

	return backup, nil
}

func (s *BackupService) ListBackups(ctx context.Context) ([]domain.Backup, error) {
	backups, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}

	return backups, nil
}
