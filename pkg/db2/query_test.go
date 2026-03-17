// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstructQuery(t *testing.T) {
	tests := []struct {
		name         string
		inputRequest *Request
		wantQuery    string
		wantArgs     []any
		wantErr      bool
	}{
		{
			name: "basic_query_without_cursor",
			inputRequest: &Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				UniqueAttributeExternalID: "id",
				UniqueKeyColumns:          []string{"id"},
				PageSize:                  100,
			},
			wantQuery: `SELECT *, (CAST("id" AS VARCHAR(50))) AS str_id, COUNT(*) OVER() AS total_remaining_rows ` +
				`FROM "users" ORDER BY CAST("id" AS VARCHAR(50)) FETCH FIRST 101 ROWS ONLY`,
			wantArgs: []any{},
			wantErr:  false,
		},
		{
			name: "query_with_cursor",
			inputRequest: &Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				UniqueAttributeExternalID: "id",
				UniqueKeyColumns:          []string{"id"},
				Cursor:                    strPtr("123"),
				PageSize:                  50,
			},
			wantQuery: `SELECT *, (CAST("id" AS VARCHAR(50))) AS str_id, COUNT(*) OVER() AS total_remaining_rows ` +
				`FROM "users" WHERE (CAST("id" AS VARCHAR(50))) > (?) ORDER BY CAST("id" AS VARCHAR(50)) FETCH FIRST 51 ROWS ONLY`,
			wantArgs: []any{"123"},
			wantErr:  false,
		},
		{
			name: "query_with_filter",
			inputRequest: &Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
				},
				UniqueAttributeExternalID: "id",
				UniqueKeyColumns:          []string{"id"},
				PageSize:                  50,
				Filter: &condexpr.Condition{
					Field:    "BUKRS",
					Operator: "=",
					Value:    "US02",
				},
			},
			wantQuery: `SELECT *, (CAST("id" AS VARCHAR(50))) AS str_id, COUNT(*) OVER() AS total_remaining_rows ` +
				`FROM "users" WHERE ("BUKRS" = ?) ORDER BY CAST("id" AS VARCHAR(50)) FETCH FIRST 51 ROWS ONLY`,
			wantArgs: []any{"US02"},
			wantErr:  false,
		},
		{
			name: "query_without_page_size",
			inputRequest: &Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "products",
				},
				UniqueAttributeExternalID: "product_id",
				PageSize:                  0,
			},
			wantQuery: `SELECT *, CAST("product_id" AS VARCHAR(50)) AS str_id, COUNT(*) OVER() AS total_remaining_rows ` +
				`FROM "products" ORDER BY CAST("product_id" AS VARCHAR(50))`,
			wantArgs: []any{},
			wantErr:  false,
		},
		{
			name:         "nil_request",
			inputRequest: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotArgs, gotErr := ConstructQuery(tt.inputRequest)

			if tt.wantErr {
				assert.Error(t, gotErr)

				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.wantQuery, gotQuery)
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}

func TestConstructQueryShouldQuoteSpecialCharacterColumns(t *testing.T) {
	tests := []struct {
		name              string
		attributeExtID    string
		wantQueryContains string
	}{
		{
			name:              "column_with_hyphen_is_quoted",
			attributeExtID:    "user-name",
			wantQueryContains: `"user-name"`,
		},
		{
			name:              "column_with_slash_is_quoted",
			attributeExtID:    "path/resource",
			wantQueryContains: `"path/resource"`,
		},
		{
			name:              "column_with_space_is_quoted",
			attributeExtID:    "First Name",
			wantQueryContains: `"First Name"`,
		},
		{
			name:              "column_with_embedded_quote_is_escaped",
			attributeExtID:    `col"name`,
			wantQueryContains: `"col""name"`,
		},
		{
			name:              "column_alphanumeric_is_quoted",
			attributeExtID:    "username",
			wantQueryContains: `"username"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			request := &Request{
				EntityConfig: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: tt.attributeExtID},
					},
				},
				UniqueAttributeExternalID: "id",
				PageSize:                  100,
			}

			// Act
			gotQuery, _, gotErr := ConstructQuery(request)

			// Assert
			require.NoError(t, gotErr)
			assert.Contains(t, gotQuery, tt.wantQueryContains)
		})
	}
}

func TestBuildSelectColumns(t *testing.T) {
	tests := []struct {
		name                          string
		attributes                    []*framework.AttributeConfig
		uniqueKeyColumns              []string
		expectedColumns               []string
		expectedNeedsCompositeKeyCols bool
	}{
		{
			name:                          "empty_attributes_returns_wildcard",
			attributes:                    []*framework.AttributeConfig{},
			uniqueKeyColumns:              []string{},
			expectedColumns:               []string{"*"},
			expectedNeedsCompositeKeyCols: false,
		},
		{
			name: "single_attribute_returns_column",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "name"},
			},
			uniqueKeyColumns:              []string{},
			expectedColumns:               []string{`"name"`},
			expectedNeedsCompositeKeyCols: false,
		},
		{
			name: "id_attribute_triggers_composite_key_flag",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "id"},
				{ExternalId: "name"},
			},
			uniqueKeyColumns:              []string{"MANDT", "EBELN"},
			expectedColumns:               []string{`"name"`, `"MANDT"`, `"EBELN"`},
			expectedNeedsCompositeKeyCols: true,
		},
		{
			name: "composite_key_columns_deduplicated",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "id"},
				{ExternalId: "MANDT"},
			},
			uniqueKeyColumns:              []string{"MANDT", "EBELN"},
			expectedColumns:               []string{`"MANDT"`, `"EBELN"`},
			expectedNeedsCompositeKeyCols: true,
		},
		{
			name: "special_characters_quoted",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "user-name"},
				{ExternalId: "path/resource"},
			},
			uniqueKeyColumns:              []string{},
			expectedColumns:               []string{`"user-name"`, `"path/resource"`},
			expectedNeedsCompositeKeyCols: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			columns, needsCompositeKeyCols := buildSelectColumns(tt.attributes, tt.uniqueKeyColumns)

			// Assert
			assert.Equal(t, tt.expectedColumns, columns)
			assert.Equal(t, tt.expectedNeedsCompositeKeyCols, needsCompositeKeyCols)
		})
	}
}

func TestBuildPaginationCastExpression(t *testing.T) {
	tests := []struct {
		name             string
		uniqueAttrID     string
		uniqueKeyColumns []string
		expected         string
	}{
		{
			name:             "empty_unique_attr_returns_empty",
			uniqueAttrID:     "",
			uniqueKeyColumns: []string{},
			expected:         "",
		},
		{
			name:             "single_column_key",
			uniqueAttrID:     "user_id",
			uniqueKeyColumns: []string{},
			expected:         `CAST("user_id" AS VARCHAR(50)) AS str_id`,
		},
		{
			name:             "composite_key_concatenates_columns",
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{"MANDT", "EBELN", "EBELP"},
			expected: `(CAST("MANDT" AS VARCHAR(50)) || '|' || ` +
				`CAST("EBELN" AS VARCHAR(50)) || '|' || ` +
				`CAST("EBELP" AS VARCHAR(50))) AS str_id`,
		},
		{
			name:             "id_attr_without_key_columns_returns_empty",
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{},
			expected:         "",
		},
		{
			name:             "special_characters_quoted",
			uniqueAttrID:     "user-id",
			uniqueKeyColumns: []string{},
			expected:         `CAST("user-id" AS VARCHAR(50)) AS str_id`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := buildPaginationCastExpression(tt.uniqueAttrID, tt.uniqueKeyColumns)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildOrderByClause(t *testing.T) {
	tests := []struct {
		name             string
		uniqueAttrID     string
		uniqueKeyColumns []string
		expected         string
	}{
		{
			name:             "empty_unique_attr_returns_empty",
			uniqueAttrID:     "",
			uniqueKeyColumns: []string{},
			expected:         "",
		},
		{
			name:             "single_column_key",
			uniqueAttrID:     "user_id",
			uniqueKeyColumns: []string{},
			expected:         `ORDER BY CAST("user_id" AS VARCHAR(50))`,
		},
		{
			name:             "composite_key_orders_by_all_columns",
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{"MANDT", "EBELN", "EBELP"},
			expected: `ORDER BY CAST("MANDT" AS VARCHAR(50)), ` +
				`CAST("EBELN" AS VARCHAR(50)), CAST("EBELP" AS VARCHAR(50))`,
		},
		{
			name:             "id_without_key_columns_treats_id_as_column",
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{},
			expected:         `ORDER BY CAST("id" AS VARCHAR(50))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := buildOrderByClause(tt.uniqueAttrID, tt.uniqueKeyColumns)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildCursorCondition(t *testing.T) {
	tests := []struct {
		name             string
		cursor           *string
		uniqueAttrID     string
		uniqueKeyColumns []string
		expectedCond     string
		expectedArgs     []any
	}{
		{
			name:             "nil_cursor_returns_empty",
			cursor:           nil,
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{"id"},
			expectedCond:     "",
			expectedArgs:     nil,
		},
		{
			name:             "empty_cursor_returns_empty",
			cursor:           strPtr(""),
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{"id"},
			expectedCond:     "",
			expectedArgs:     nil,
		},
		{
			name:             "single_column_cursor",
			cursor:           strPtr("123"),
			uniqueAttrID:     "user_id",
			uniqueKeyColumns: []string{},
			expectedCond:     `CAST("user_id" AS VARCHAR(50)) > ?`,
			expectedArgs:     []any{"123"},
		},
		{
			name:             "composite_key_cursor",
			cursor:           strPtr("100|4500001234|10"),
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{"MANDT", "EBELN", "EBELP"},
			expectedCond: `(CAST("MANDT" AS VARCHAR(50)), ` +
				`CAST("EBELN" AS VARCHAR(50)), CAST("EBELP" AS VARCHAR(50))) > (?, ?, ?)`,
			expectedArgs: []any{"100", "4500001234", "10"},
		},
		{
			name:             "cursor_parts_mismatch_returns_empty",
			cursor:           strPtr("100|4500001234"),
			uniqueAttrID:     "id",
			uniqueKeyColumns: []string{"MANDT", "EBELN", "EBELP"},
			expectedCond:     "",
			expectedArgs:     nil,
		},
		{
			name:             "empty_unique_attr_returns_empty",
			cursor:           strPtr("123"),
			uniqueAttrID:     "",
			uniqueKeyColumns: []string{},
			expectedCond:     "",
			expectedArgs:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			cond, args := buildCursorCondition(tt.cursor, tt.uniqueAttrID, tt.uniqueKeyColumns)

			// Assert
			assert.Equal(t, tt.expectedCond, cond)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestQuoteIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "alphanumeric_is_quoted",
			input:    "column_name",
			expected: `"column_name"`,
		},
		{
			name:     "hyphen_is_quoted",
			input:    "user-name",
			expected: `"user-name"`,
		},
		{
			name:     "slash_is_quoted",
			input:    "path/resource",
			expected: `"path/resource"`,
		},
		{
			name:     "space_is_quoted",
			input:    "First Name",
			expected: `"First Name"`,
		},
		{
			name:     "embedded_quote_is_escaped",
			input:    `col"name`,
			expected: `"col""name"`,
		},
		{
			name:     "multiple_special_chars",
			input:    `my "col-name`,
			expected: `"my ""col-name"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := quoteIdentifier(tt.input)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

// strPtr is a helper to create string pointers for tests.
func strPtr(s string) *string {
	return &s
}
