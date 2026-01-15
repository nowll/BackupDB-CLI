package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nowll/BackupDB-CLI/internal/domain"
)

type sqliteRepository struct {
	config domain.DatabaseConfig
}

func NewSQLiteRepository(config domain.DatabaseConfig) domain.DatabaseRepository {
	return &sqliteRepository{config: config}
}

func (r *sqliteRepository) TestConnection(ctx context.Context) error {
	db, err := sql.Open("sqlite3", r.config.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.PingContext(ctx)
}

func (r *sqliteRepository) Backup(ctx context.Context, config domain.BackupConfig) ([]byte, error) {
	data, err := os.ReadFile(r.config.Database)
	if err != nil {
		return nil, fmt.Errorf("read sqlite file: %w", err)
	}

	return data, nil
}

func (r *sqliteRepository) Restore(ctx context.Context, data []byte, config domain.RestoreConfig) error {
	if config.DropExisting {
		if err := os.Remove(r.config.Database); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove existing db: %w", err)
		}
	}

	if err := os.WriteFile(r.config.Database, data, 0644); err != nil {
		return fmt.Errorf("write sqlite file: %w", err)
	}

	return nil
}

func (r *sqliteRepository) GetDatabases(ctx context.Context) ([]string, error) {
	return []string{r.config.Database}, nil
}

func (r *sqliteRepository) GetTables(ctx context.Context, database string) ([]string, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT name FROM sqlite_master WHERE type='table'`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}
