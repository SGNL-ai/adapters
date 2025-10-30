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

		if cursorMap := parseCursorFromLog(gotLog, "responseNextCursor"); cursorMap != nil {
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

// parseCursorFromLog takes a log map and dereferences any CompositeCursor pointer
// in the cursorField (typically "responseNextCursor"), converting it to a map.
func parseCursorFromLog(log map[string]any, cursorField string) map[string]any {
	cursorPtr, exists := log[cursorField]
	if !exists {
		return nil
	}

	// Try parsing as either int64 or string cursor.
	if cursor, ok := cursorPtr.(*pagination.CompositeCursor[int64]); ok && cursor != nil {
		return compositeCursorToMap(cursor)
	}

	if cursor, ok := cursorPtr.(*pagination.CompositeCursor[string]); ok && cursor != nil {
		return compositeCursorToMap(cursor)
	}

	return nil
}

// compositeCursorToMap converts any CompositeCursor to a map.
func compositeCursorToMap[T int64 | string](cursor *pagination.CompositeCursor[T]) map[string]any {
	cursorMap := make(map[string]any)

	if cursor.Cursor != nil {
		cursorMap["cursor"] = *cursor.Cursor
	}

	if cursor.CollectionID != nil {
		cursorMap["collectionId"] = *cursor.CollectionID
	}

	if cursor.CollectionCursor != nil {
		cursorMap["collectionCursor"] = *cursor.CollectionCursor
	}

	return cursorMap
}
