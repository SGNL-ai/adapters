// Copyright 2026 SGNL.ai, Inc.
package testutil

func GenPtr[T any](v T) *T {
	return &v
}
