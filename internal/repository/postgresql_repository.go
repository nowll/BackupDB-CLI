package repository

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os/exec"

	_ "github.com/lib/pq"
	"github.com/nowll/BackupDB-CLI/internal/domain"
)

type postgresRepository struct {
	config domain.DatabaseConfig
}

func NewPostgreSQLRepository(config domain.DatabaseConfig) domain.DatabaseRepository {
	return &postgresRepository{config: config}
}

func (r *postgresRepository) TestConnection(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		r.config.Host,
		r.config.Port,
		r.config.Username,
		r.config.Password,
		r.config.Database,
		r.config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.PingContext(ctx)
}

func (r *postgresRepository) Backup(ctx context.Context, config domain.BackupConfig) ([]byte, error) {
	args := []string{
		fmt.Sprintf("--host=%s", r.config.Host),
		fmt.Sprintf("--port=%d", r.config.Port),
		fmt.Sprintf("--username=%s", r.config.Username),
		"--format=custom",
	}

	if config.Database != "" {
		args = append(args, fmt.Sprintf("--dbname=%s", config.Database))
	}

	if len(config.Tables) > 0 {
		for _, table := range config.Tables {
			args = append(args, fmt.Sprintf("--table=%s", table))
		}
	}

	cmd := exec.CommandContext(ctx, "pg_dump", args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", r.config.Password))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("pg_dump: %w", err)
	}

	return output, nil
}

func (r *postgresRepository) Restore(ctx context.Context, data []byte, config domain.RestoreConfig) error {
	args := []string{
		fmt.Sprintf("--host=%s", r.config.Host),
		fmt.Sprintf("--port=%d", r.config.Port),
		fmt.Sprintf("--username=%s", r.config.Username),
	}

	if config.Database != "" {
		args = append(args, fmt.Sprintf("--dbname=%s", config.Database))
	}

	if config.IgnoreErrors {
		args = append(args, "--no-owner", "--no-acl")
	}

	cmd := exec.CommandContext(ctx, "pg_restore", args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", r.config.Password))
	cmd.Stdin = bytes.NewReader(data)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_restore: %w, output: %s", err, output)
	}

	return nil
}

func (r *postgresRepository) GetDatabases(ctx context.Context) ([]string, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		r.config.Host,
		r.config.Port,
		r.config.Username,
		r.config.Password,
		r.config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, "SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	return databases, nil
}

func (r *postgresRepository) GetTables(ctx context.Context, database string) ([]string, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		r.config.Host,
		r.config.Port,
		r.config.Username,
		r.config.Password,
		database,
		r.config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT tablename FROM pg_tables WHERE schemaname = 'public'`
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
