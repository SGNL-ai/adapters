// Copyright 2025 SGNL.ai, Inc.
package zaplogger

import (
	"github.com/spf13/viper"
)

const (
	LogModeConsole = "console"
	LogModeFile    = "file"
)

type Config struct {
	// Mode sets the logging mode. Valid modes are: "console", "file".
	Mode []string `yaml:"mode" json:"mode" mapstructure:"mode"`
	// Level sets the logging level. Valid levels are: "DEBUG", "INFO", "WARN", "ERROR", "DPANIC", "PANIC", and "FATAL".
	Level string `yaml:"level" json:"level" mapstructure:"level"`

	// The following fields are only used if "file" is included in Mode.
	// FilePath sets the file path for file logging.
	FilePath string `yaml:"file_path" json:"file_path" mapstructure:"file_path"`
	// FileMaxSize sets the maximum size in megabytes of the log file before it gets rotated.
	FileMaxSize int `yaml:"file_max_size" json:"file_max_size" mapstructure:"file_max_size"`
	// FileMaxBackups sets the maximum number of old log files to retain.
	FileMaxBackups int `yaml:"file_max_backups" json:"file_max_backups" mapstructure:"file_max_backups"`
	// FileMaxDays sets the maximum number of days to retain old log files.
	FileMaxDays int `yaml:"file_max_days" json:"file_max_days" mapstructure:"file_max_days"`

	// ServiceName is an optional field that, if set, adds the service name to each log entry.
	ServiceName string `yaml:"service_name" json:"service_name" mapstructure:"service_name"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix("SGNL_LOG")
	v.AutomaticEnv()

	v.SetDefault("level", "INFO")
	v.SetDefault("mode", "console")
	v.SetDefault("file_path", "/var/log/sgnl/adapter-sgnl.log")
	v.SetDefault("file_max_size", 100)
	v.SetDefault("file_max_days", 7)
	v.SetDefault("file_max_backups", 10)
	v.SetDefault("service_name", "")

	var cfg Config

	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
