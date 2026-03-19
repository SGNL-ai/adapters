// Copyright 2026 SGNL.ai, Inc.

// Contract tests for the DB2 adapter using recorded fixtures.
// These tests verify that the adapter produces consistent output for known inputs.

package db2_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/sgnl-ai/adapters/pkg/db2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFixtureLoader verifies that fixtures can be loaded correctly.
func TestFixtureLoader(t *testing.T) {
	t.Run("load items_with_filter fixture", func(t *testing.T) {
		fixture, err := LoadFixture("items_with_filter")
		require.NoError(t, err)
		assert.Equal(t, "items_with_filter", fixture.Name)
		assert.Equal(t, "ITEMS", fixture.Request.Entity)
		assert.Len(t, fixture.Response.Objects, 10)
		assert.Empty(t, fixture.Response.NextCursor) // Full page, no more data
	})

	t.Run("load items_small_page fixture with pagination", func(t *testing.T) {
		fixture, err := LoadFixture("items_small_page")
		require.NoError(t, err)
		assert.Equal(t, "items_small_page", fixture.Name)
		assert.Len(t, fixture.Response.Objects, 3)
		assert.NotEmpty(t, fixture.Response.NextCursor) // Has more pages
	})

	t.Run("load items_no_filter fixture", func(t *testing.T) {
		fixture, err := LoadFixture("items_no_filter")
		require.NoError(t, err)
		assert.Equal(t, "items_no_filter", fixture.Name)
		assert.Len(t, fixture.Response.Objects, 5)
		assert.NotEmpty(t, fixture.Response.NextCursor) // Has more pages
	})

	t.Run("load all fixtures", func(t *testing.T) {
		fixtures, err := LoadAllFixtures()
		require.NoError(t, err)
		assert.Len(t, fixtures, 4)
	})

	t.Run("error on missing fixture", func(t *testing.T) {
		_, err := LoadFixture("nonexistent")
		assert.Error(t, err)
	})
}

// TestFixtureMockRows verifies the MockRows implementation works correctly.
func TestFixtureMockRows(t *testing.T) {
	objects := []map[string]interface{}{
		{"id": "T1|D1001|L01", "TENANT_ID": "T1", "DOC_NUM": "D1001", "AMOUNT": 150.5},
		{"id": "T1|D1001|L02", "TENANT_ID": "T1", "DOC_NUM": "D1001", "AMOUNT": 275.0},
	}
	columns := []string{"id", "TENANT_ID", "DOC_NUM", "AMOUNT"}

	rows := NewFixtureMockRows(objects, columns)

	// First row
	assert.True(t, rows.Next())
	var id, tenantID, docNum string
	var amount float64
	err := rows.Scan(&id, &tenantID, &docNum, &amount)
	require.NoError(t, err)
	assert.Equal(t, "T1|D1001|L01", id)
	assert.Equal(t, "T1", tenantID)
	assert.Equal(t, "D1001", docNum)
	assert.Equal(t, 150.5, amount)

	// Second row
	assert.True(t, rows.Next())
	err = rows.Scan(&id, &tenantID, &docNum, &amount)
	require.NoError(t, err)
	assert.Equal(t, "T1|D1001|L02", id)
	assert.Equal(t, 275.0, amount)

	// No more rows
	assert.False(t, rows.Next())
	assert.NoError(t, rows.Err())
	assert.NoError(t, rows.Close())
}

// TestContractItemsWithFilter verifies adapter behavior with filter.
func TestContractItemsWithFilter(t *testing.T) {
	fixture, err := LoadFixture("items_with_filter")
	require.NoError(t, err)

	// Create a mock client that simulates both unique key query and data query
	queryCount := 0
	mockClient := createContractMockClient(fixture, &queryCount)

	datasource := db2.NewClient(mockClient)
	adapter := db2.NewAdapter(datasource)

	// Build request matching the fixture
	request := &framework.Request[db2.Config]{
		Auth: &framework.DatasourceAuthCredentials{
			Basic: &framework.BasicAuthCredentials{
				Username: "testuser",
				Password: "testpass",
			},
		},
		Address:  "localhost:50000",
		PageSize: fixture.Request.PageSize,
		Entity: framework.EntityConfig{
			ExternalId: fixture.Request.Entity,
			Attributes: buildAttributeConfigs(fixture.Request.Attributes),
		},
		Config: &db2.Config{
			Database: fixture.Request.Database,
			Schema:   fixture.Request.Schema,
			Filters:  buildFiltersFromFixture(fixture),
		},
	}

	response := adapter.GetPage(context.Background(), request)

	// Verify response structure
	assert.Nil(t, response.Error, "Expected no error")
	require.NotNil(t, response.Success, "Expected success response")

	// Verify object count matches fixture
	assert.Len(t, response.Success.Objects, len(fixture.Response.Objects))

	// Verify pagination - if fixture has next cursor, response should too
	if fixture.Response.NextCursor != "" {
		assert.NotEmpty(t, response.Success.NextCursor, "Expected next cursor")
	} else {
		assert.Empty(t, response.Success.NextCursor, "Expected no next cursor")
	}
}

// TestContractSmallPage verifies adapter behavior with pagination.
func TestContractSmallPage(t *testing.T) {
	fixture, err := LoadFixture("items_small_page")
	require.NoError(t, err)

	queryCount := 0
	mockClient := createContractMockClient(fixture, &queryCount)

	datasource := db2.NewClient(mockClient)
	adapter := db2.NewAdapter(datasource)

	request := &framework.Request[db2.Config]{
		Auth: &framework.DatasourceAuthCredentials{
			Basic: &framework.BasicAuthCredentials{
				Username: "testuser",
				Password: "testpass",
			},
		},
		Address:  "localhost:50000",
		PageSize: fixture.Request.PageSize,
		Entity: framework.EntityConfig{
			ExternalId: fixture.Request.Entity,
			Attributes: buildAttributeConfigs(fixture.Request.Attributes),
		},
		Config: &db2.Config{
			Database: fixture.Request.Database,
			Schema:   fixture.Request.Schema,
			Filters:  buildFiltersFromFixture(fixture),
		},
	}

	response := adapter.GetPage(context.Background(), request)

	assert.Nil(t, response.Error)
	require.NotNil(t, response.Success)
	assert.Len(t, response.Success.Objects, len(fixture.Response.Objects))
	assert.NotEmpty(t, response.Success.NextCursor, "Should have next cursor for pagination")
}

// TestContractPaginationContinuation verifies that using a cursor returns the next page.
func TestContractPaginationContinuation(t *testing.T) {
	// Load page 1 fixture
	page1, err := LoadFixture("items_small_page")
	require.NoError(t, err)
	require.NotEmpty(t, page1.Response.NextCursor, "Page 1 must have a cursor for this test")

	// Load page 2 fixture (uses cursor from page 1)
	page2, err := LoadFixture("items_small_page_2")
	if err != nil {
		t.Skip("Page 2 fixture not recorded yet - run db2_record_fixtures.go to record it")
	}

	// Verify page 2 uses the cursor from page 1
	assert.Equal(t, page1.Response.NextCursor, page2.Request.Cursor,
		"Page 2 should use cursor from page 1")

	// Verify page 2 returns different objects than page 1
	page1IDs := make(map[string]bool)
	for _, obj := range page1.Response.Objects {
		if id, ok := obj["id"].(string); ok {
			page1IDs[id] = true
		}
	}

	for _, obj := range page2.Response.Objects {
		if id, ok := obj["id"].(string); ok {
			assert.False(t, page1IDs[id], "Page 2 should not contain objects from page 1: %s", id)
		}
	}

	// Test that adapter correctly handles the cursor
	queryCount := 0
	mockClient := createContractMockClient(page2, &queryCount)
	datasource := db2.NewClient(mockClient)
	adapter := db2.NewAdapter(datasource)

	request := &framework.Request[db2.Config]{
		Auth: &framework.DatasourceAuthCredentials{
			Basic: &framework.BasicAuthCredentials{
				Username: "testuser",
				Password: "testpass",
			},
		},
		Address:  "localhost:50000",
		PageSize: page2.Request.PageSize,
		Cursor:   page2.Request.Cursor, // Use cursor from page 1
		Entity: framework.EntityConfig{
			ExternalId: page2.Request.Entity,
			Attributes: buildAttributeConfigs(page2.Request.Attributes),
		},
		Config: &db2.Config{
			Database: page2.Request.Database,
			Schema:   page2.Request.Schema,
			Filters:  buildFiltersFromFixture(page2),
		},
	}

	response := adapter.GetPage(context.Background(), request)

	assert.Nil(t, response.Error)
	require.NotNil(t, response.Success)
	assert.Len(t, response.Success.Objects, len(page2.Response.Objects))
}

// TestContractNoFilter verifies adapter behavior without filter.
func TestContractNoFilter(t *testing.T) {
	fixture, err := LoadFixture("items_no_filter")
	require.NoError(t, err)

	queryCount := 0
	mockClient := createContractMockClient(fixture, &queryCount)

	datasource := db2.NewClient(mockClient)
	adapter := db2.NewAdapter(datasource)

	request := &framework.Request[db2.Config]{
		Auth: &framework.DatasourceAuthCredentials{
			Basic: &framework.BasicAuthCredentials{
				Username: "testuser",
				Password: "testpass",
			},
		},
		Address:  "localhost:50000",
		PageSize: fixture.Request.PageSize,
		Entity: framework.EntityConfig{
			ExternalId: fixture.Request.Entity,
			Attributes: buildAttributeConfigs(fixture.Request.Attributes),
		},
		Config: &db2.Config{
			Database: fixture.Request.Database,
			Schema:   fixture.Request.Schema,
		},
	}

	response := adapter.GetPage(context.Background(), request)

	assert.Nil(t, response.Error)
	require.NotNil(t, response.Success)
	assert.Len(t, response.Success.Objects, len(fixture.Response.Objects))
}

// Helper functions

func buildAttributeConfigs(attrs []string) []*framework.AttributeConfig {
	configs := make([]*framework.AttributeConfig, len(attrs))
	for i, attr := range attrs {
		attrType := framework.AttributeTypeString
		if attr == "AMOUNT" {
			attrType = framework.AttributeTypeDouble
		}
		configs[i] = &framework.AttributeConfig{
			ExternalId: attr,
			Type:       attrType,
			UniqueId:   attr == "id",
		}
	}

	return configs
}

func buildFiltersFromFixture(fixture *Fixture) map[string]condexpr.Condition {
	if fixture.Request.Filter == nil {
		return nil
	}

	return map[string]condexpr.Condition{
		fixture.Request.Entity: *fixture.Request.Filter,
	}
}

// createContractMockClient creates a MockSQLClient that handles both
// unique key queries and data queries using fixture data.
func createContractMockClient(fixture *Fixture, queryCount *int) db2.SQLClient {
	return &db2.MockSQLClient{
		ConnectFunc: func(_ string) (*sql.DB, error) { return nil, nil },
		QueryFunc: func(_ context.Context, query string, _ ...interface{}) (db2.Rows, error) {
			*queryCount++

			// First query is for unique key columns (SYSCAT.TABCONST)
			if strings.Contains(query, "SYSCAT.TABCONST") || strings.Contains(query, "SYSCAT.KEYCOLUSE") {
				return createUniqueKeyMockRows(fixture.Request.Entity), nil
			}

			// Subsequent queries are for actual data
			return createDataMockRows(fixture), nil
		},
	}
}

// createUniqueKeyMockRows creates MockRows that return the table's composite key columns.
func createUniqueKeyMockRows(tableName string) db2.Rows {
	// Return composite key for ITEMS table
	if tableName == "ITEMS" {
		return &db2.MockRows{
			Data: []map[string]interface{}{
				{"CONSTNAME": "PK_ITEMS", "TABNAME": "ITEMS", "COLNAME": "TENANT_ID", "COLSEQ": 1},
				{"CONSTNAME": "PK_ITEMS", "TABNAME": "ITEMS", "COLNAME": "DOC_NUM", "COLSEQ": 2},
				{"CONSTNAME": "PK_ITEMS", "TABNAME": "ITEMS", "COLNAME": "LINE_NUM", "COLSEQ": 3},
			},
		}
	}

	return &db2.MockRows{Data: []map[string]interface{}{}}
}

// createDataMockRows creates MockRows from fixture response data.
func createDataMockRows(fixture *Fixture) db2.Rows {
	// Build column list: requested attributes + total_remaining_rows
	columns := append([]string{}, fixture.Request.Attributes...)
	columns = append(columns, "total_remaining_rows")

	// If fixture has next cursor, add an extra row to trigger pagination
	objects := fixture.Response.Objects
	hasMore := fixture.Response.NextCursor != ""

	// Add total_remaining_rows to each object
	objectsWithCount := make([]map[string]interface{}, 0, len(objects)+1)
	for _, obj := range objects {
		objCopy := make(map[string]interface{})
		for k, v := range obj {
			objCopy[k] = v
		}

		if hasMore {
			objCopy["total_remaining_rows"] = int64(1)
		} else {
			objCopy["total_remaining_rows"] = int64(0)
		}

		objectsWithCount = append(objectsWithCount, objCopy)
	}

	// If there's a next cursor, add one extra object to trigger pagination
	if hasMore && len(objectsWithCount) > 0 {
		extra := make(map[string]interface{})
		for k, v := range objectsWithCount[len(objectsWithCount)-1] {
			extra[k] = v
		}

		extra["total_remaining_rows"] = int64(0)
		objectsWithCount = append(objectsWithCount, extra)
	}

	return NewFixtureMockRows(objectsWithCount, columns)
}
