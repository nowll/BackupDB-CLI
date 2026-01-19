package storage

import (
	"fmt"

	"github.com/nowll/db-backup-cli/internal/domain"
)

func NewStorage(config domain.StorageConfig) (domain.StorageRepository, error) {
	switch config.Type {
	case domain.LocalStorage:
		return NewLocalStorage(config), nil
	case domain.S3Storage:
		return NewS3Storage(config)
	case domain.GCSStorage:
		return NewGCSStorage(config)
	case domain.AzureStorage:
		return NewAzureStorage(config)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}
