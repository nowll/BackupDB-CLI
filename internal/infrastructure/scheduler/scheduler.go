package scheduler

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/yourusername/db-backup-cli/internal/domain"
)

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) AddBackupJob(
	schedule string,
	usecase domain.BackupUseCase,
	dbConfig domain.DatabaseConfig,
	backupConfig domain.BackupConfig,
) error {
	_, err := s.cron.AddFunc(schedule, func() {
		ctx := context.Background()
		_, err := usecase.CreateBackup(ctx, dbConfig, backupConfig)
		if err != nil {
			fmt.Printf("Scheduled backup failed: %v\n", err)
		}
	})

	return err
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
