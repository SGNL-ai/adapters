// Copyright 2026 SGNL.ai, Inc.

package zaplogger

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

// Logger wraps a *zap.Logger to implement the framework_logs.Logger interface.
type Logger struct {
	logger *zap.Logger
}

var _ framework_logs.Logger = (*Logger)(nil)

// NewFrameworkLogger creates a new framework_logs.Logger from a *zap.Logger.
func NewFrameworkLogger(logger *zap.Logger) framework_logs.Logger {
	return &Logger{logger: logger}
}

// Info logs an informational message.
func (a *Logger) Info(msg string, fields ...framework_logs.Field) {
	a.logger.Info(msg, toZapFields(fields)...)
}

// Error logs an error message.
func (a *Logger) Error(msg string, fields ...framework_logs.Field) {
	a.logger.Error(msg, toZapFields(fields)...)
}

// Debug logs a debug message.
func (a *Logger) Debug(msg string, fields ...framework_logs.Field) {
	a.logger.Debug(msg, toZapFields(fields)...)
}

// With creates a child logger with pre-attached fields.
func (a *Logger) With(fields ...framework_logs.Field) framework_logs.Logger {
	return &Logger{
		logger: a.logger.With(toZapFields(fields)...),
	}
}

// toZapFields converts framework_logs.Field to zap.Field.
func toZapFields(fields []framework_logs.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}

	return zapFields
}

// Unwrap returns the underlying *zap.Logger.
// This allows consumers to access zap-specific features when needed.
func (a *Logger) Unwrap() *zap.Logger {
	return a.logger
}

// UnwrapLogger attempts to extract a *zap.Logger from a framework_logs.Logger.
// Returns the underlying *zap.Logger and true if the logger is a zaplogger.Logger,
// otherwise returns nil and false.
func UnwrapLogger(logger framework_logs.Logger) (*zap.Logger, bool) {
	if adapter, ok := logger.(*Logger); ok {
		return adapter.Unwrap(), true
	}

	return nil, false
}

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

	if cfg.ServiceName != "" {
		logger = logger.With(zap.String("serviceName", cfg.ServiceName))
	}

	// Replace the global logger zap.L() with the newly created one.
	zap.ReplaceGlobals(logger)

	// Redirect standard library logs to the zap logger for consistency.
	_, err = zap.RedirectStdLogAt(logger, logLevel)
	if err != nil {
		log.Fatalf("Can't redirect std to zap logger: %v", err)
	}

	logger.Info("Zap logger initialized", zap.Any("config", cfg))

	return logger
}

// FromContext returns a logger from the context if available.
// It's a thin wrapper around framework_logs.FromContext.
// If no logger is found in context, returns the global logger as fallback.
//
// The logger from the framework context already has useful request fields attached. See the framework.
func FromContext(ctx context.Context) *zap.Logger {
	if logger := framework_logs.FromContext(ctx); logger != nil {
		zapLogger, ok := UnwrapLogger(logger)
		if ok {
			return zapLogger
		}
	}

	zap.L().Warn("No logger found in context, falling back to global logger")

	// Return the global logger as a fallback.
	return zap.L()
}
