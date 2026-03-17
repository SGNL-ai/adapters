//go:build db2
// +build db2

// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	// IBM DB2 Driver - only imported when building with -tags db2
	_ "github.com/ibmdb/go_ibm_db"
)
