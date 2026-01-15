package repository

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/nowll/BackupDB-CLI/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	config domain.DatabaseConfig
}

func NewMongoDBRepository(config domain.DatabaseConfig) domain.DatabaseRepository {
	return &mongoRepository{config: config}
}

func (r *mongoRepository) TestConnection(ctx context.Context) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	return client.Ping(ctx, nil)
}

func (r *mongoRepository) Backup(ctx context.Context, config domain.BackupConfig) ([]byte, error) {
	args := []string{
		fmt.Sprintf("--host=%s", r.config.Host),
		fmt.Sprintf("--port=%d", r.config.Port),
		fmt.Sprintf("--username=%s", r.config.Username),
		fmt.Sprintf("--password=%s", r.config.Password),
		"--archive",
	}

	if config.Database != "" {
		args = append(args, fmt.Sprintf("--db=%s", config.Database))
	}

	if len(config.Tables) > 0 {
		for _, collection := range config.Tables {
			args = append(args, fmt.Sprintf("--collection=%s", collection))
		}
	}

	cmd := exec.CommandContext(ctx, "mongodump", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("mongodump: %w", err)
	}

	return output, nil
}

func (r *mongoRepository) Restore(ctx context.Context, data []byte, config domain.RestoreConfig) error {
	args := []string{
		fmt.Sprintf("--host=%s", r.config.Host),
		fmt.Sprintf("--port=%d", r.config.Port),
		fmt.Sprintf("--username=%s", r.config.Username),
		fmt.Sprintf("--password=%s", r.config.Password),
		"--archive",
	}

	if config.Database != "" {
		args = append(args, fmt.Sprintf("--db=%s", config.Database))
	}

	if config.DropExisting {
		args = append(args, "--drop")
	}

	cmd := exec.CommandContext(ctx, "mongorestore", args...)
	cmd.Stdin = bytes.NewReader(data)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mongorestore: %w, output: %s", err, output)
	}

	return nil
}

func (r *mongoRepository) GetDatabases(ctx context.Context) ([]string, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	return client.ListDatabaseNames(ctx, map[string]interface{}{})
}

func (r *mongoRepository) GetTables(ctx context.Context, database string) ([]string, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	db := client.Database(database)
	return db.ListCollectionNames(ctx, map[string]interface{}{})
}
