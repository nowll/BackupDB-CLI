package domain

import "context"

type StorageType string

const (
	LocalStorage StorageType = "local"
	S3Storage    StorageType = "s3"
	GCSStorage   StorageType = "gcs"
	AzureStorage StorageType = "azure"
)

type StorageConfig struct {
	Type      StorageType
	Path      string
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
	Options   map[string]string
}

type StorageRepository interface {
	Upload(ctx context.Context, key string, data []byte) error
	Download(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]string, error)
	Exists(ctx context.Context, key string) (bool, error)
}
