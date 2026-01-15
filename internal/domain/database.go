package domain

import (
	"context"
)

type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgresql"
	MongoDB    DatabaseType = "mongodb"
	SQLite     DatabaseType = "sqlite"
)

type DatabaseConfig struct {
	Type     DatabaseType
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
	Options  map[string]string
}

type DatabaseRepository interface {
	TestConnection(ctx context.Context) error
	Backup(ctx context.Context, config BackupConfig) ([]byte, error)
	Restore(ctx context.Context, data []byte, config RestoreConfig) error
	GetDatabases(ctx context.Context) ([]string, error)
	GetTables(ctx context.Context, database string) ([]string, error)
}
