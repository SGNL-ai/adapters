// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatasource_NewClient(t *testing.T) {
	mockClient := NewMockSQLClient()
	client := NewClient(mockClient)
	assert.NotNil(t, client)
}
