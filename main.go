import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/nowll/db-backup-cli/internal/domain"
	"github.com/nowll/db-backup-cli/internal/infrastructure/logger"
	"github.com/nowll/db-backup-cli/internal/infrastructure/notifier"
	"github.com/nowll/db-backup-cli/internal/infrastructure/scheduler"
	"github.com/nowll/db-backup-cli/internal/repository"
	"github.com/nowll/db-backup-cli/internal/storage"
	"github.com/nowll/db-backup-cli/internal/usecase"
	"github.com/nowll/db-backup-cli/pkg/config"
)

var (
	configFile  string
	dbType      string
	dbHost      string
	dbPort      int
	dbUser      string
	dbPass      string
	dbName      string
	backupType  string
	compress    bool
	storageType string
	storagePath string
	scheduleStr string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "db-backup",
		Short: "Database backup CLI utility",
		Long:  `A comprehensive database backup tool supporting multiple database types and storage options`,
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file path")

	// Backup command
	backupCmd := &cobra.Command{
		Use:   "backup",
		Short: "Create a database backup",
		RunE:  runBackup,
	}

	backupCmd.Flags().StringVar(&dbType, "db-type", "mysql", "Database type (mysql, postgresql, mongodb, sqlite)")
	backupCmd.Flags().StringVar(&dbHost, "host", "localhost", "Database host")
	backupCmd.Flags().IntVar(&dbPort, "port", 3306, "Database port")
	backupCmd.Flags().StringVar(&dbUser, "user", "", "Database username")
	backupCmd.Flags().StringVar(&dbPass, "password", "", "Database password")
	backupCmd.Flags().StringVar(&dbName, "database", "", "Database name")
	backupCmd.Flags().StringVar(&backupType, "type", "full", "Backup type (full, incremental, differential)")
	backupCmd.Flags().BoolVar(&compress, "compress", true, "Compress backup")
	backupCmd.Flags().StringVar(&storageType, "storage", "local", "Storage type (local, s3, gcs, azure)")
	backupCmd.Flags().StringVar(&storagePath, "path", "./backups", "Storage path")

	// Restore command
	restoreCmd := &cobra.Command{
		Use:   "restore [backup-id]",
		Short: "Restore from a backup",
		Args:  cobra.ExactArgs(1),
		RunE:  runRestore,
	}

	restoreCmd.Flags().StringVar(&dbType, "db-type", "mysql", "Database type")
	restoreCmd.Flags().StringVar(&dbHost, "host", "localhost", "Database host")
	restoreCmd.Flags().IntVar(&dbPort, "port", 3306, "Database port")
	restoreCmd.Flags().StringVar(&dbUser, "user", "", "Database username")
	restoreCmd.Flags().StringVar(&dbPass, "password", "", "Database password")
	restoreCmd.Flags().StringVar(&dbName, "database", "", "Database name")

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all backups",
		RunE:  runList,
	}

	listCmd.Flags().StringVar(&storagePath, "path", "./backups", "Storage path")

	// Schedule command
	scheduleCmd := &cobra.Command{
		Use:   "schedule",
		Short: "Schedule automatic backups",
		RunE:  runSchedule,
	}

	scheduleCmd.Flags().StringVar(&scheduleStr, "cron", "", "Cron schedule (e.g., '0 2 * * *' for daily at 2am)")
	scheduleCmd.Flags().StringVar(&dbType, "db-type", "mysql", "Database type")
	scheduleCmd.Flags().StringVar(&dbHost, "host", "localhost", "Database host")
	scheduleCmd.Flags().IntVar(&dbPort, "port", 3306, "Database port")
	scheduleCmd.Flags().StringVar(&dbUser, "user", "", "Database username")
	scheduleCmd.Flags().StringVar(&dbPass, "password", "", "Database password")
	scheduleCmd.Flags().StringVar(&dbName, "database", "", "Database name")
	scheduleCmd.Flags().BoolVar(&compress, "compress", true, "Compress backup")
	scheduleCmd.Flags().StringVar(&storageType, "storage", "local", "Storage type")
	scheduleCmd.Flags().StringVar(&storagePath, "path", "./backups", "Storage path")

	// Test connection command
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test database connection",
		RunE:  runTest,
	}

	testCmd.Flags().StringVar(&dbType, "db-type", "mysql", "Database type")
	testCmd.Flags().StringVar(&dbHost, "host", "localhost", "Database host")
	testCmd.Flags().IntVar(&dbPort, "port", 3306, "Database port")
	testCmd.Flags().StringVar(&dbUser, "user", "", "Database username")
	testCmd.Flags().StringVar(&dbPass, "password", "", "Database password")
	testCmd.Flags().StringVar(&dbName, "database", "", "Database name")

	rootCmd.AddCommand(backupCmd, restoreCmd, listCmd, scheduleCmd, testCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runBackup(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize logger
	log, err := logger.NewLogger("logs/backup.log")
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}
	defer log.Sync()

	// Create database config
	dbConfig := domain.DatabaseConfig{
		Type:     domain.DatabaseType(dbType),
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUser,
		Password: dbPass,
		Database: dbName,
	}

	// Create storage config
	storageConfig := domain.StorageConfig{
		Type: domain.StorageType(storageType),
		Path: storagePath,
	}

	// Initialize storage
	store, err := storage.NewStorage(storageConfig)
	if err != nil {
		return fmt.Errorf("initialize storage: %w", err)
	}

	// Initialize repositories
	dbFactory := repository.NewRepositoryFactory()
	metadataRepo := repository.NewMetadataRepository(storagePath)

	// Initialize notifier (optional)
	var notif usecase.Notifier
	if slackToken := os.Getenv("SLACK_TOKEN"); slackToken != "" {
		notif = notifier.NewSlackNotifier(slackToken, os.Getenv("SLACK_CHANNEL"))
	}

	// Create backup use case
	backupUC := usecase.NewBackupUseCase(dbFactory, store, log, notif, metadataRepo)

	// Create backup config
	backupConfig := domain.BackupConfig{
		Type:     domain.BackupType(backupType),
		Database: dbName,
		Compress: compress,
	}

	// Perform backup
	metadata, err := backupUC.CreateBackup(ctx, dbConfig, backupConfig)
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	fmt.Printf("Backup completed successfully!\n")
	fmt.Printf("Backup ID: %s\n", metadata.ID)
	fmt.Printf("Size: %d bytes (compressed: %d bytes)\n", metadata.Size, metadata.CompressedSize)
	fmt.Printf("Duration: %s\n", metadata.EndTime.Sub(metadata.StartTime))

	return nil
}

func runRestore(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	backupID := args[0]

	// Initialize logger
	log, err := logger.NewLogger("logs/restore.log")
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}
	defer log.Sync()

	// Create database config
	dbConfig := domain.DatabaseConfig{
		Type:     domain.DatabaseType(dbType),
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUser,
		Password: dbPass,
		Database: dbName,
	}

	// Create storage config
	storageConfig := domain.StorageConfig{
		Type: domain.StorageType(storageType),
		Path: storagePath,
	}

	// Initialize storage
	store, err := storage.NewStorage(storageConfig)
	if err != nil {
		return fmt.Errorf("initialize storage: %w", err)
	}

	// Initialize repositories
	dbFactory := repository.NewRepositoryFactory()
	metadataRepo := repository.NewMetadataRepository(storagePath)

	// Create restore use case
	restoreUC := usecase.NewRestoreUseCase(dbFactory, store, log, metadataRepo)

	// Create restore config
	restoreConfig := domain.RestoreConfig{
		Database:     dbName,
		DropExisting: false,
		IgnoreErrors: false,
	}

	// Perform restore
	if err := restoreUC.RestoreBackup(ctx, backupID, dbConfig, restoreConfig); err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	fmt.Printf("Restore completed successfully!\n")

	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	metadataRepo := repository.NewMetadataRepository(storagePath)

	backups, err := metadataRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("list backups: %w", err)
	}

	if len(backups) == 0 {
		fmt.Println("No backups found")
		return nil
	}

	fmt.Printf("Found %d backup(s):\n\n", len(backups))

	for _, backup := range backups {
		fmt.Printf("ID: %s\n", backup.ID)
		fmt.Printf("  Type: %s (%s)\n", backup.DatabaseType, backup.BackupType)
		fmt.Printf("  Database: %s\n", backup.DatabaseName)
		fmt.Printf("  Date: %s\n", backup.StartTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Size: %d bytes (compressed: %d bytes)\n", backup.Size, backup.CompressedSize)
		fmt.Printf("  Status: %s\n", backup.Status)
		fmt.Println()
	}

	return nil
}

func runSchedule(cmd *cobra.Command, args []string) error {
	if scheduleStr == "" {
		return fmt.Errorf("--cron flag is required")
	}

	// Initialize logger
	log, err := logger.NewLogger("logs/scheduled.log")
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}
	defer log.Sync()

	// Create configs
	dbConfig := domain.DatabaseConfig{
		Type:     domain.DatabaseType(dbType),
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUser,
		Password: dbPass,
		Database: dbName,
	}

	storageConfig := domain.StorageConfig{
		Type: domain.StorageType(storageType),
		Path: storagePath,
	}

	// Initialize components
	store, err := storage.NewStorage(storageConfig)
	if err != nil {
		return fmt.Errorf("initialize storage: %w", err)
	}

	dbFactory := repository.NewRepositoryFactory()
	metadataRepo := repository.NewMetadataRepository(storagePath)

	var notif usecase.Notifier
	if slackToken := os.Getenv("SLACK_TOKEN"); slackToken != "" {
		notif = notifier.NewSlackNotifier(slackToken, os.Getenv("SLACK_CHANNEL"))
	}

	backupUC := usecase.NewBackupUseCase(dbFactory, store, log, notif, metadataRepo)

	backupConfig := domain.BackupConfig{
		Type:     domain.BackupType(backupType),
		Database: dbName,
		Compress: compress,
	}

	// Create and start scheduler
	sched := scheduler.NewScheduler()

	if err := sched.AddBackupJob(scheduleStr, backupUC, dbConfig, backupConfig); err != nil {
		return fmt.Errorf("add scheduled job: %w", err)
	}

	sched.Start()

	fmt.Printf("Backup scheduled with cron: %s\n", scheduleStr)
	fmt.Println("Press Ctrl+C to stop...")

	// Keep running
	select {}
}

func runTest(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	dbConfig := domain.DatabaseConfig{
		Type:     domain.DatabaseType(dbType),
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUser,
		Password: dbPass,
		Database: dbName,
	}

	factory := repository.NewRepositoryFactory()
	repo, err := factory.Create(dbConfig)
	if err != nil {
		return fmt.Errorf("create repository: %w", err)
	}

	fmt.Printf("Testing connection to %s database at %s:%d...\n", dbType, dbHost, dbPort)

	if err := repo.TestConnection(ctx); err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		return err
	}

	fmt.Println("✓ Connection successful!")

	return nil
}
