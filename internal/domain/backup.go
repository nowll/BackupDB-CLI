package domain

import (
	"context"
	"time"
)

type BackupType string

const (
	FullBackup         BackupType = "full"
	IncrementalBackup  BackupType = "incremental"
	DifferentialBackup BackupType = "differential"
)

type BackupConfig struct {
	Type          BackupType
	Database      string
	Tables        []string
	Compress      bool
	Encryption    bool
	ExcludeTables []string
}

type RestoreConfig struct {
	Database     string
	Tables       []string
	DropExisting bool
	IgnoreErrors bool
}

type BackupMetadata struct {
	ID              string
	DatabaseType    DatabaseType
	BackupType      BackupType
	DatabaseName    string
	StartTime       time.Time
	EndTime         time.Time
	Size            int64
	CompressedSize  int64
	Status          string
	Error           string
	Checksum        string
	StorageLocation string
}

type BackupUseCase interface {
	CreateBackup(ctx context.Context, dbConfig DatabaseConfig, backupConfig BackupConfig) (*BackupMetadata, error)
	ScheduleBackup(schedule string, dbConfig DatabaseConfig, backupConfig BackupConfig) error
	ListBackups(ctx context.Context) ([]*BackupMetadata, error)
}

type RestoreUseCase interface {
	RestoreBackup(ctx context.Context, backupID string, dbConfig DatabaseConfig, restoreConfig RestoreConfig) error
	ValidateBackup(ctx context.Context, backupID string) error
}
