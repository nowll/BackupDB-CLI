package repository

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nowll/BackupDB-CLI/internal/domain"
)

type mysqlRepository struct {
	config domain.DatabaseConfig
}

func NewMySQLRepository(config domain.DatabaseConfig) domain.DatabaseRepository {
	return &mysqlRepository{config: config}
}

func (r *mysqlRepository) TestConnection(ctx context.Context) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
		r.config.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.PingContext(ctx)
}

func (r *mysqlRepository) Backup(ctx context.Context, config domain.BackupConfig) ([]byte, error) {
	args := []string{
		fmt.Sprintf("--host=%s", r.config.Host),
		fmt.Sprintf("--port=%d", r.config.Port),
		fmt.Sprintf("--user=%s", r.config.Username),
		fmt.Sprintf("--password=%s", r.config.Password),
	}

	if config.Database != "" {
		args = append(args, config.Database)
	} else {
		args = append(args, "--all-databases")
	}

	if len(config.Tables) > 0 {
		args = append(args, "--tables")
		args = append(args, config.Tables...)
	}

	cmd := exec.CommandContext(ctx, "mysqldump", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("mysqldump: %w", err)
	}

	return output, nil
}

func (r *mysqlRepository) Restore(ctx context.Context, data []byte, config domain.RestoreConfig) error {
	args := []string{
		fmt.Sprintf("--host=%s", r.config.Host),
		fmt.Sprintf("--port=%d", r.config.Port),
		fmt.Sprintf("--user=%s", r.config.Username),
		fmt.Sprintf("--password=%s", r.config.Password),
	}

	if config.Database != "" {
		args = append(args, config.Database)
	}

	cmd := exec.CommandContext(ctx, "mysql", args...)
	cmd.Stdin = bytesReader(data)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mysql restore: %w, output: %s", err, output)
	}

	return nil
}

func (r *mysqlRepository) GetDatabases(ctx context.Context) ([]string, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, "SHOW DATABASES")
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

func (r *mysqlRepository) GetTables(ctx context.Context, database string) ([]string, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
		database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, "SHOW TABLES")
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

func bytesReader(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}
