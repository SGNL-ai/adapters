package logs

import (
	"github.com/spf13/viper"
)

type Config struct {
	// Mode sets the logging mode. Valid modes are: "console", "file".
	Mode []string `yaml:"mode" json:"mode" mapstructure:"mode"`
	// Level sets the logging level. Valid levels are: "DEBUG", "INFO", "WARN", "ERROR", "DPANIC", "PANIC", and "FATAL".
	Level string `json:"level"`

	// The following fields are only used if "file" is included in Mode.
	// FilePath sets the file path for file logging.
	FilePath string `yaml:"file_path" json:"file_path" mapstructure:"file_path"`
	// FileMaxSize sets the maximum size in megabytes of the log file before it gets rotated.
	FileMaxSize int `yaml:"file_max_size" json:"file_max_size" mapstructure:"file_max_size"`
	// FileMaxBackups sets the maximum number of old log files to retain.
	FileMaxBackups int `yaml:"file_max_backups" json:"file_max_backups" mapstructure:"file_max_backups"`
	// FileMaxDays sets the maximum number of days to retain old log files.
	FileMaxDays int `yaml:"file_max_days" json:"file_max_days" mapstructure:"file_max_days"`
}

func LoadConfig() Config {
	v := viper.New()
	v.SetEnvPrefix("SGNL_LOG")
	v.AutomaticEnv()

	v.SetDefault("level", "INFO")
	v.SetDefault("mode", "console")
	v.SetDefault("file_path", "/var/log/sgnl/unconfigured.log")
	v.SetDefault("file_max_size", 100)
	v.SetDefault("file_max_days", 7)
	v.SetDefault("file_max_backups", 10)

	var cfg Config

	v.UnmarshalExact(&cfg)

	return cfg
}
