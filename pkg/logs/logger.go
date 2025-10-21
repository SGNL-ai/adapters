package logs

import (
	"context"
	"log"

	framework_logs "github.com/sgnl-ai/adapter-framework/pkg/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new zap.Logger based on the provided configuration.
// It uses sensible production defaults with JSON formatting and nanosecond
// precision for timestamps.
// It accepts a user supplied Config and optional zap options.
func New(cfg Config, zapOpts ...zap.Option) *zap.Logger {
	logLevel, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		log.Fatal("Failed to parse log level")
	}

	zapCfg := zap.NewProductionConfig()

	// Disable sampling to ensure all logs are captured.
	zapCfg.Sampling = nil

	// Add nanosecond precision to the timestamp.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	zapCfg.EncoderConfig = encoderCfg

	// Replace the log level.
	zapCfg.Level = zap.NewAtomicLevelAt(logLevel)

	// Build the logger with the above configuration.
	logger, err := zapCfg.Build(zapOpts...)
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	zap.ReplaceGlobals(logger)

	// Redirect standard library logs to the zap logger for consistency.
	_, err = zap.RedirectStdLogAt(logger, logLevel)
	if err != nil {
		log.Fatalf("Can't redirect std to zap logger: %v", err)
	}

	logger.Info("Zap logger initialized")

	return logger
}

// FromContext returns a logger from the context if available.
// It's a thin wrapper around framework_logs.LoggerFromContext.
// If no logger is found in context, returns the global logger as fallback.
func FromContext(ctx context.Context) *zap.Logger {
	if logger := framework_logs.LoggerFromContext(ctx); logger != nil {
		return logger
	}

	// Return the global logger as a fallback.
	return zap.L()
}
