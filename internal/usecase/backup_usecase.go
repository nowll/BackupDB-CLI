package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/nowll/BackupDB-CLI/internal/domain"
	"github.com/nowll/BackupDB-CLI/pkg/compression"
)

type backupUseCase struct {
	dbRepoFactory DatabaseRepositoryFactory
	storage       domain.StorageRepository
	logger        Logger
	notifier      Notifier
	metadataRepo  MetadataRepository
}

type DatabaseRepositoryFactory interface {
	Create(config domain.DatabaseConfig) (domain.DatabaseRepository, error)
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

type Notifier interface {
	Notify(ctx context.Context, message string) error
}

type MetadataRepository interface {
	Save(ctx context.Context, metadata *domain.BackupMetadata) error
	Get(ctx context.Context, id string) (*domain.BackupMetadata, error)
	List(ctx context.Context) ([]*domain.BackupMetadata, error)
}

func NewBackupUseCase(
	factory DatabaseRepositoryFactory,
	storage domain.StorageRepository,
	logger Logger,
	notifier Notifier,
	metadataRepo MetadataRepository,
) domain.BackupUseCase {
	return &backupUseCase{
		dbRepoFactory: factory,
		storage:       storage,
		logger:        logger,
		notifier:      notifier,
		metadataRepo:  metadataRepo,
	}
}

func (u *backupUseCase) CreateBackup(
	ctx context.Context,
	dbConfig domain.DatabaseConfig,
	backupConfig domain.BackupConfig,
) (*domain.BackupMetadata, error) {
	startTime := time.Now()

	u.logger.Info("Starting backup",
		"database_type", dbConfig.Type,
		"backup_type", backupConfig.Type,
		"database", backupConfig.Database,
	)

	// Create database repository
	dbRepo, err := u.dbRepoFactory.Create(dbConfig)
	if err != nil {
		u.logger.Error("Failed to create database repository", "error", err)
		return nil, fmt.Errorf("create repository: %w", err)
	}

	// Test connection
	if err := dbRepo.TestConnection(ctx); err != nil {
		u.logger.Error("Database connection failed", "error", err)
		return nil, fmt.Errorf("test connection: %w", err)
	}

	data, err := dbRepo.Backup(ctx, backupConfig)
	if err != nil {
		u.logger.Error("Backup failed", "error", err)
		return nil, fmt.Errorf("backup: %w", err)
	}

	originalSize := int64(len(data))

	var compressedData []byte
	var compressedSize int64

	if backupConfig.Compress {
		compressedData, err = compression.Compress(data)
		if err != nil {
			u.logger.Error("Compression failed", "error", err)
			return nil, fmt.Errorf("compress: %w", err)
		}
		compressedSize = int64(len(compressedData))
	} else {
		compressedData = data
		compressedSize = originalSize
	}

	hash := sha256.Sum256(compressedData)
	checksum := hex.EncodeToString(hash[:])

	backupID := fmt.Sprintf("%s_%s_%s",
		dbConfig.Type,
		backupConfig.Database,
		time.Now().Format("20060102_150405"),
	)

	storageKey := fmt.Sprintf("backups/%s/%s.backup", dbConfig.Type, backupID)

	// Upload to storage
	if err := u.storage.Upload(ctx, storageKey, compressedData); err != nil {
		u.logger.Error("Upload failed", "error", err)
		return nil, fmt.Errorf("upload: %w", err)
	}

	endTime := time.Now()

	metadata := &domain.BackupMetadata{
		ID:              backupID,
		DatabaseType:    dbConfig.Type,
		BackupType:      backupConfig.Type,
		DatabaseName:    backupConfig.Database,
		StartTime:       startTime,
		EndTime:         endTime,
		Size:            originalSize,
		CompressedSize:  compressedSize,
		Status:          "completed",
		Checksum:        checksum,
		StorageLocation: storageKey,
	}

	if err := u.metadataRepo.Save(ctx, metadata); err != nil {
		u.logger.Error("Failed to save metadata", "error", err)
		return nil, fmt.Errorf("save metadata: %w", err)
	}

	u.logger.Info("Backup completed successfully",
		"backup_id", backupID,
		"duration", endTime.Sub(startTime),
		"size", originalSize,
		"compressed_size", compressedSize,
	)

	// Send notification
	if u.notifier != nil {
		msg := fmt.Sprintf("Backup completed: %s (Duration: %s, Size: %d bytes)",
			backupID,
			endTime.Sub(startTime),
			compressedSize,
		)
		if err := u.notifier.Notify(ctx, msg); err != nil {
			u.logger.Warn("Failed to send notification", "error", err)
		}
	}

	return metadata, nil
}

func (u *backupUseCase) ScheduleBackup(
	schedule string,
	dbConfig domain.DatabaseConfig,
	backupConfig domain.BackupConfig,
) error {
	return fmt.Errorf("not implemented")
}

func (u *backupUseCase) ListBackups(ctx context.Context) ([]*domain.BackupMetadata, error) {
	return u.metadataRepo.List(ctx)
}
