// Copyright 2025 SGNL.ai, Inc.
package zaplogger_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	framework_logs "github.com/sgnl-ai/adapter-framework/pkg/logs"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var (
	MockClockTime      = time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC)
	MockClockTimestamp = MockClockTime.Format(time.RFC3339)
)

func TestNew(t *testing.T) {
	mockClock := newMockClock()

	tests := map[string]struct {
		config            zaplogger.Config
		writeLogs         func(logger *zap.Logger)
		expectedLogs      []map[string]any
		expectedFileLines []map[string]any
	}{
		"console_mode_only": {
			config: zaplogger.Config{
				Mode:  []string{"console"},
				Level: "INFO",
			},
			writeLogs: func(logger *zap.Logger) {
				logger.Debug("debug message")
				logger.Info("info message")
				logger.Warn("warn message")
				logger.Error("error message")
			},
			expectedLogs: []map[string]any{
				{
					"level": "info",
					"ts":    MockClockTimestamp,
					"msg":   "Zap logger initialized",
				},
				{
					"level": "info",
					"ts":    MockClockTimestamp,
					"msg":   "info message",
				},
				{
					"level": "warn",
					"ts":    MockClockTimestamp,
					"msg":   "warn message",
				},
				{
					"level": "error",
					"ts":    MockClockTimestamp,
					"msg":   "error message",
				},
			},
		},
		"both_console_and_file_mode": {
			config: zaplogger.Config{
				Mode:           []string{"console", "file"},
				Level:          "DEBUG",
				FilePath:       filepath.Join(t.TempDir(), "test.log"),
				FileMaxSize:    100,
				FileMaxBackups: 10,
				FileMaxDays:    7,
			},
			writeLogs: func(logger *zap.Logger) {
				logger.Debug("debug message")
				logger.Info("info message")
				logger.Warn("warn message")
				logger.Error("error message", zap.Int("code", 500))
			},
			expectedLogs: []map[string]any{
				{
					"level": "info",
					"ts":    MockClockTimestamp,
					"msg":   "Zap logger initialized",
				},
				{
					"level": "debug",
					"ts":    MockClockTimestamp,
					"msg":   "debug message",
				},
				{
					"level": "info",
					"ts":    MockClockTimestamp,
					"msg":   "info message",
				},
				{
					"level": "warn",
					"ts":    MockClockTimestamp,
					"msg":   "warn message",
				},
				{
					"level": "error",
					"ts":    MockClockTimestamp,
					"msg":   "error message",
					"code":  int64(500),
				},
			},
			expectedFileLines: []map[string]any{
				{
					"level": "info",
					"ts":    MockClockTimestamp,
					"msg":   "Zap logger initialized",
				},
				{
					"level": "debug",
					"ts":    MockClockTimestamp,
					"msg":   "debug message",
				},
				{
					"level": "info",
					"ts":    MockClockTimestamp,
					"msg":   "info message",
				},
				{
					"level": "warn",
					"ts":    MockClockTimestamp,
					"msg":   "warn message",
				},
				{
					"level": "error",
					"ts":    MockClockTimestamp,
					"msg":   "error message",
					"code":  float64(500), // JSON unmarshals numbers as float64.
				},
			},
		},
		"debug_level": {
			config: zaplogger.Config{
				Mode:  []string{"console"},
				Level: "DEBUG",
			},
			writeLogs: func(logger *zap.Logger) {
				logger.Debug("debug message")
				logger.Info("info message")
			},
			expectedLogs: []map[string]any{
				{
					"level": "info",
					"ts":    MockClockTimestamp,
					"msg":   "Zap logger initialized",
				},
				{
					"level": "debug",
					"ts":    MockClockTimestamp,
					"msg":   "debug message",
				},
				{
					"level": "info",
					"ts":    MockClockTimestamp,
					"msg":   "info message",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			level, err := zapcore.ParseLevel(test.config.Level)
			if err != nil {
				t.Fatalf("failed to parse log level: %v", err)
			}

			// Create an observable core to capture log output.
			observedCore, observedLogs := observer.New(level)

			// Create the logger with the observable core and mock clock.
			logger := zaplogger.New(test.config,
				zap.WrapCore(func(c zapcore.Core) zapcore.Core {
					return zapcore.NewTee(c, observedCore)
				}),
				zap.WithClock(mockClock),
			)
			defer logger.Sync()

			if logger == nil {
				t.Fatal("expected logger to be created")
			}

			// Write logs using the test-specific function.
			if test.writeLogs != nil {
				test.writeLogs(logger)
			}

			// Validate expected logs.
			if len(test.expectedLogs) > 0 {
				gotLogs := observedLogs.All()

				if len(gotLogs) != len(test.expectedLogs) {
					t.Errorf("expected %d logs, got %d", len(test.expectedLogs), len(gotLogs))
				}

				for i, expectedLog := range test.expectedLogs {
					gotLog := gotLogs[i].ContextMap()           // Get all the log fields as a map.
					gotLog["msg"] = gotLogs[i].Message          // Add the "msg" field since that's not included in ContextMap().
					gotLog["level"] = gotLogs[i].Level.String() // Add the "level" field.
					gotLog["ts"] = MockClockTimestamp           // Add the "ts" field to match expected logs.

					if !reflect.DeepEqual(gotLog, expectedLog) {
						t.Errorf("log %d mismatch:\ngot:  %#v\nwant: %#v", i, gotLog, expectedLog)
					}
				}
			}

			// Validate file contents if requested.
			if len(test.expectedFileLines) > 0 {
				content, err := os.ReadFile(test.config.FilePath)
				if err != nil {
					t.Fatalf("failed to read log file: %v", err)
				}

				if len(content) == 0 {
					t.Fatal("log file is empty")
				}

				// Split into lines and parse each as JSON.
				lines := strings.Split(strings.TrimSpace(string(content)), "\n")

				if len(lines) != len(test.expectedFileLines) {
					t.Fatalf("expected %d log lines, got %d", len(test.expectedFileLines), len(lines))
				}

				for i, rawLine := range lines {
					var log map[string]any
					if err := json.Unmarshal([]byte(rawLine), &log); err != nil {
						t.Fatalf("failed to parse JSON on line %d: %v\nLine: %s", i+1, err, rawLine)
					}

					expectedLog := test.expectedFileLines[i]

					for field, expectedValue := range expectedLog {
						gotValue, ok := log[field]
						if !ok {
							t.Errorf("line %d: missing field %q", i+1, field)

							continue
						}

						if gotValue != expectedValue {
							t.Errorf("line %d: field %q = %v, want %v", i+1, field, gotValue, expectedValue)
						}
					}
				}
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	// Create a global logger first so we have something to compare against.
	config := zaplogger.Config{
		Mode:  []string{"console"},
		Level: "INFO",
	}
	globalLogger := zaplogger.New(config)

	// Create a separate logger to store in context.
	contextLogger := zap.NewNop()

	tests := map[string]struct {
		setupCtx   func() context.Context
		wantLogger *zap.Logger
	}{
		"returns_global_logger_when_not_present_in_context": {
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantLogger: globalLogger, // Should return global logger as fallback
		},
		"returns_logger_from_context_when_present": {
			setupCtx: func() context.Context {
				ctx := context.Background()

				return framework_logs.NewContextWithLogger(ctx, zaplogger.NewFrameworkLogger(contextLogger))
			},
			wantLogger: contextLogger,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := tc.setupCtx()
			retrievedLogger := zaplogger.FromContext(ctx)

			if tc.wantLogger != retrievedLogger {
				t.Error("FromContext returned different logger than expected")
			}
		})
	}
}

// mockClock implements zapcore.Clock for testing with a fixed time.
type mockClock struct {
	now time.Time
}

var _ zapcore.Clock = (*mockClock)(nil)

func newMockClock() *mockClock {
	return &mockClock{now: MockClockTime}
}

func (m *mockClock) Now() time.Time {
	return m.now
}

func (m *mockClock) NewTicker(duration time.Duration) *time.Ticker {
	return time.NewTicker(duration)
}
