package repository

import (
	"fmt"

	"github.com/nowll/db-backup-cli/internal/domain"
)

type RepositoryFactory struct{}

func NewRepositoryFactory() *RepositoryFactory {
	return &RepositoryFactory{}
}

func (f *RepositoryFactory) Create(config domain.DatabaseConfig) (domain.DatabaseRepository, error) {
	switch config.Type {
	case domain.MySQL:
		return NewMySQLRepository(config), nil
	case domain.PostgreSQL:
		return NewPostgreSQLRepository(config), nil
	case domain.MongoDB:
		return NewMongoDBRepository(config), nil
	case domain.SQLite:
		return NewSQLiteRepository(config), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}
