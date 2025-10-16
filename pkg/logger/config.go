package logger

import "github.com/spf13/viper"

type Config struct {
	// Level sets the logging level. Valid levels are: "DEBUG", "INFO", "WARN", "ERROR", "DPANIC", "PANIC", and "FATAL".
	Level string `json:"level"`
}

func LoadConfig() Config {
	v := viper.New()
	v.SetEnvPrefix("SGNL_LOG")
	v.AutomaticEnv()

	v.SetDefault("level", "INFO")

	cfg := Config{
		Level: v.GetString("level"),
	}

	return cfg
}
