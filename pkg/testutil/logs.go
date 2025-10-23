// Copyright 2025 SGNL.ai, Inc.
package testutil

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	framework_logs "github.com/sgnl-ai/adapter-framework/pkg/logs"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// NewContextWithObservableLogger creates a new context enriched with an observable logger
// which can be used in tests to capture and inspect log output.
func NewContextWithObservableLogger(ctx context.Context) (context.Context, *observer.ObservedLogs) {
	// Create an observable logger to capture log output.
	observedCore, observedLogs := observer.New(zapcore.InfoLevel)
	observableLogger := zap.New(observedCore)

	// Enrich context with logger
	return framework_logs.NewContextWithLogger(ctx, zaplogger.NewFrameworkLogger(observableLogger)), observedLogs
}

// ValidateLogOutput compares the observed logs against the expected logs.
// No comparison happens if expectedLogs is nil or empty.
// The observed logs are extracted from the provided observer.ObservedLogs and
// additional fields like "msg" and "level" are added for comparison.
// If the comparison fails, an error is reported using the testing.T instance.
func ValidateLogOutput(t *testing.T, observedLogs *observer.ObservedLogs, expectedLogs []map[string]any) {
	if len(expectedLogs) == 0 {
		return
	}

	gotLogs := observedLogs.All()

	if len(gotLogs) != len(expectedLogs) {
		t.Errorf("expected %d logs, got %d", len(expectedLogs), len(gotLogs))
	}

	for i, expectedLog := range expectedLogs {
		gotLog := gotLogs[i].ContextMap()
		gotLog["msg"] = gotLogs[i].Message
		gotLog["level"] = gotLogs[i].Level.String()

		if cursorMap := pagination.ParseCursorFromLog(gotLog, "responseNextCursor"); cursorMap != nil {
			gotLog["responseNextCursor"] = cursorMap
		}

		// Parse responseBody if it's a json.RawMessage.
		if responseBody, ok := gotLog[fields.FieldResponseBody]; ok {
			if rawJSON, ok := responseBody.(json.RawMessage); ok {
				var parsed map[string]any

				if err := json.Unmarshal(rawJSON, &parsed); err == nil {
					gotLog[fields.FieldResponseBody] = parsed
				}
			}
		}

		if !reflect.DeepEqual(gotLog, expectedLog) {
			t.Errorf("log %d mismatch: got: %v, want: %v", i, gotLog, expectedLog)
		}
	}
}
