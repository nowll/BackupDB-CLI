package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nowll/db-backup-cli/internal/domain"
)

type MetadataRepo struct {
	storagePath string
}

func NewMetadataRepository(storagePath string) *MetadataRepo {
	return &MetadataRepo{
		storagePath: storagePath,
	}
}

func (r *MetadataRepo) Save(ctx context.Context, metadata *domain.BackupMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	metadataDir := filepath.Join(r.storagePath, "metadata")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		return fmt.Errorf("create metadata directory: %w", err)
	}

	filePath := filepath.Join(metadataDir, metadata.ID+".json")
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("write metadata file: %w", err)
	}

	return nil
}

func (r *MetadataRepo) Get(ctx context.Context, id string) (*domain.BackupMetadata, error) {
	filePath := filepath.Join(r.storagePath, "metadata", id+".json")

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read metadata file: %w", err)
	}

	var metadata domain.BackupMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}

	return &metadata, nil
}

func (r *MetadataRepo) List(ctx context.Context) ([]*domain.BackupMetadata, error) {
	metadataDir := filepath.Join(r.storagePath, "metadata")

	files, err := ioutil.ReadDir(metadataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*domain.BackupMetadata{}, nil
		}
		return nil, fmt.Errorf("read metadata directory: %w", err)
	}

	var metadataList []*domain.BackupMetadata

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join(metadataDir, file.Name()))
		if err != nil {
			continue
		}

		var metadata domain.BackupMetadata
		if err := json.Unmarshal(data, &metadata); err != nil {
			continue
		}

		metadataList = append(metadataList, &metadata)
	}

	return metadataList, nil
}
