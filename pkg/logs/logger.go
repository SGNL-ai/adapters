// Copyright 2025 SGNL.ai, Inc.
package logs

import (
	"context"
	"log"
	"os"
	"slices"

	framework_logs "github.com/sgnl-ai/adapter-framework/pkg/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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

	// Add nanosecond precision to the timestamp.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	zapCores := make([]zapcore.Core, 0, len(cfg.Mode))

	if slices.Contains(cfg.Mode, LogModeFile) {
		zapCores = append(zapCores, zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    cfg.FileMaxSize, // megabytes
				MaxBackups: cfg.FileMaxBackups,
				MaxAge:     cfg.FileMaxDays, // days
				Compress:   true,
			}),
			logLevel,
		))
	}

	if slices.Contains(cfg.Mode, LogModeConsole) {
		zapCores = append(zapCores, zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(os.Stdout),
			logLevel,
		))
	}

	core := zapcore.NewTee(zapCores...)

	logger := zap.New(core, zapOpts...)

	// Replace the global logger zap.L() with the newly created one.
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
//
// The logger from the framework context already has useful request fields attached. See the framework.
func FromContext(ctx context.Context) *zap.Logger {
	if logger := framework_logs.LoggerFromContext(ctx); logger != nil {
		return logger
	}

	// Return the global logger as a fallback.
	return zap.L()
}
