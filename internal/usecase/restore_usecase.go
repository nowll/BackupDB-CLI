package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/nowll/BackupDB-CLI/internal/domain"
	"github.com/nowll/BackupDB-CLI/pkg/compression"
)

type restoreUseCase struct {
    dbRepoFactory DatabaseRepositoryFactory
    storage       domain.StorageRepository
    logger        Logger
    metadataRepo  MetadataRepository
}

func NewRestoreUseCase(
    factory DatabaseRepositoryFactory,
    storage domain.StorageRepository,
    logger Logger,
    metadataRepo MetadataRepository,
) domain.RestoreUseCase {
    return &restoreUseCase{
        dbRepoFactory: factory,
        storage:       storage,
        logger:        logger,
        metadataRepo:  metadataRepo,
    }
}

func (u *restoreUseCase) RestoreBackup(
    ctx context.Context,
    backupID string,
    dbConfig domain.DatabaseConfig,
    restoreConfig domain.RestoreConfig,
) error {
    u.logger.Info("Starting restore", "backup_id", backupID)

    // Get backup metadata
    metadata, err := u.metadataRepo.Get(ctx, backupID)
    if err != nil {
        return fmt.Errorf("get metadata: %w", err)
    }

    // Download backup
    compressedData, err := u.storage.Download(ctx, metadata.StorageLocation)
    if err != nil {
        return fmt.Errorf("download backup: %w", err)
    }

    // Verify checksum
    hash := sha256.Sum256(compressedData)
    checksum := hex.EncodeToString(hash[:])

    if checksum != metadata.Checksum {
        return fmt.Errorf("checksum mismatch: expected %s, got %s", metadata.Checksum, checksum)
    }

    // Decompress
    data, err := compression.Decompress(compressedData)
    if err != nil {
        return fmt.Errorf("decompress: %w", err)
    }

    // Create database repository
    dbRepo, err := u.dbRepoFactory.Create(dbConfig)
    if err != nil {
        return fmt.Errorf("create repository: %w", err)
    }

    // Test connection
    if err := dbRepo.TestConnection(ctx); err != nil {
        return fmt.Errorf("test connection: %w", err)
    }

    // Restore
    if err := dbRepo.Restore(ctx, data, restoreConfig); err != nil {
        return fmt.Errorf("restore: %w", err)
    }

    u.logger.Info("Restore completed successfully", "backup_id", backupID);
