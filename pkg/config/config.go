package config

import (
	"fmt"

	"github.com/nowll/BackupDB-CLI/internal/domain"
	"github.com/spf13/viper"
)

type Config struct {
	Database domain.DatabaseConfig
	Storage  domain.StorageConfig
	Logging  LoggingConfig
	Slack    SlackConfig
}

type LoggingConfig struct {
	Level    string
	FilePath string
}

type SlackConfig struct {
	Token   string
	Channel string
	Enabled bool
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &config, nil
}
