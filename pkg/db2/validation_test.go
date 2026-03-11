// Copyright 2026 SGNL.ai, Inc.
package db2_test

import (
	"strings"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/db2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRequestFromConfig tests the nil checks in NewRequestFromConfig.
func TestNewRequestFromConfig(t *testing.T) {
	tests := []struct {
		name           string
		request        *framework.Request[db2.Config]
		wantErr        bool
		wantErrContain string
	}{
		{
			name: "valid_config",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "localhost",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "test_table",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id", UniqueId: true},
					},
				},
				Config: &db2.Config{
					Database: "TESTDB",
				},
			},
			wantErr: false,
		},
		{
			name:           "nil_request",
			request:        nil,
			wantErr:        true,
			wantErrContain: "request is nil",
		},
		{
			name: "nil_auth",
			request: &framework.Request[db2.Config]{
				Auth:     nil,
				Address:  "localhost",
				PageSize: 100,
				Config:   &db2.Config{Database: "TESTDB"},
			},
			wantErr:        true,
			wantErrContain: "auth is nil",
		},
		{
			name: "nil_basic_auth",
			request: &framework.Request[db2.Config]{
				Auth:     &framework.DatasourceAuthCredentials{},
				Address:  "localhost",
				PageSize: 100,
				Config:   &db2.Config{Database: "TESTDB"},
			},
			wantErr:        true,
			wantErrContain: "Basic authentication is required",
		},
		{
			name: "nil_config",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "localhost",
				PageSize: 100,
				Config:   nil,
			},
			wantErr:        true,
			wantErrContain: "request contains no config",
		},
		{
			name: "address_with_https_scheme_rejected",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "https://mydb2host:50000",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "test_table",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id", UniqueId: true},
					},
				},
				Config: &db2.Config{
					Database: "TESTDB",
				},
			},
			wantErr:        true,
			wantErrContain: "is not supported",
		},
		{
			name: "address_with_http_scheme_rejected",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "http://mydb2host:50000",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "test_table",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id", UniqueId: true},
					},
				},
				Config: &db2.Config{
					Database: "TESTDB",
				},
			},
			wantErr:        true,
			wantErrContain: "is not supported",
		},
		{
			name: "address_with_whitespace_trimmed",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "  mydb2host:50000  ",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "test_table",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id", UniqueId: true},
					},
				},
				Config: &db2.Config{
					Database: "TESTDB",
				},
			},
			wantErr: false,
		},
		{
			name: "address_without_scheme_accepted",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "mydb2host:50000",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "test_table",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id", UniqueId: true},
					},
				},
				Config: &db2.Config{
					Database: "TESTDB",
				},
			},
			wantErr: false,
		},
		{
			name: "nil_attributes_handled_gracefully",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "localhost",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "test_table",
					Attributes: nil,
				},
				Config: &db2.Config{
					Database: "TESTDB",
				},
			},
			wantErr: false,
		},
		{
			name: "empty_attributes_handled_gracefully",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "localhost",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "test_table",
					Attributes: []*framework.AttributeConfig{},
				},
				Config: &db2.Config{
					Database: "TESTDB",
				},
			},
			wantErr: false,
		},
		{
			name: "nil_filters_handled_gracefully",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "localhost",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "test_table",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id", UniqueId: true},
					},
				},
				Config: &db2.Config{
					Database: "TESTDB",
					Filters:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "empty_entity_external_id_handled_gracefully",
			request: &framework.Request[db2.Config]{
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "testuser",
						Password: "testpass",
					},
				},
				Address:  "localhost",
				PageSize: 100,
				Entity: framework.EntityConfig{
					ExternalId: "",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id", UniqueId: true},
					},
				},
				Config: &db2.Config{
					Database: "TESTDB",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := db2.NewRequestFromConfig(tt.request)

			if tt.wantErr {
				require.NotNil(t, err, "expected error but got nil")
				if tt.wantErrContain != "" {
					assert.Contains(t, err.Message, tt.wantErrContain)
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// TestRequestValidate tests deep validation on the Request struct.
func TestRequestValidate(t *testing.T) {
	// Helper to create a valid base request
	validRequest := func() *db2.Request {
		return &db2.Request{
			Username: "testuser",
			Password: "testpass",
			BaseURL:  "localhost",
			Database: "TESTDB",
			PageSize: 100,
			EntityConfig: framework.EntityConfig{
				ExternalId: "test_table",
				Attributes: []*framework.AttributeConfig{
					{ExternalId: "id", UniqueId: true},
				},
			},
			UniqueAttributeExternalID: "id",
		}
	}

	tests := []struct {
		name           string
		modifyRequest  func(*db2.Request)
		wantErr        bool
		wantErrContain string
	}{
		{
			name:          "valid_request",
			modifyRequest: func(_ *db2.Request) {},
			wantErr:       false,
		},
		{
			name:           "empty_username",
			modifyRequest:  func(r *db2.Request) { r.Username = "" },
			wantErr:        true,
			wantErrContain: "username and password are required",
		},
		{
			name:           "empty_password",
			modifyRequest:  func(r *db2.Request) { r.Password = "" },
			wantErr:        true,
			wantErrContain: "username and password are required",
		},
		{
			name:           "empty_address",
			modifyRequest:  func(r *db2.Request) { r.BaseURL = "" },
			wantErr:        true,
			wantErrContain: "address (hostname) is required",
		},
		{
			name:           "missing_database",
			modifyRequest:  func(r *db2.Request) { r.Database = "" },
			wantErr:        true,
			wantErrContain: "database is not set",
		},
		{
			name: "no_unique_attribute",
			modifyRequest: func(r *db2.Request) {
				r.EntityConfig.Attributes = []*framework.AttributeConfig{
					{ExternalId: "name", UniqueId: false},
				}
			},
			wantErr:        true,
			wantErrContain: "no unique attribute defined",
		},
		{
			name: "empty_attributes",
			modifyRequest: func(r *db2.Request) {
				r.EntityConfig.Attributes = []*framework.AttributeConfig{}
			},
			wantErr:        true,
			wantErrContain: "no unique attribute defined",
		},
		{
			name:           "zero_page_size",
			modifyRequest:  func(r *db2.Request) { r.PageSize = 0 },
			wantErr:        true,
			wantErrContain: "page size must be between 1 and",
		},
		{
			name:           "negative_page_size",
			modifyRequest:  func(r *db2.Request) { r.PageSize = -1 },
			wantErr:        true,
			wantErrContain: "page size must be between 1 and",
		},
		{
			name:           "page_size_exceeds_max",
			modifyRequest:  func(r *db2.Request) { r.PageSize = config.MaxPageSize + 1 },
			wantErr:        true,
			wantErrContain: "page size must be between 1 and",
		},
		{
			name:          "max_valid_page_size",
			modifyRequest: func(r *db2.Request) { r.PageSize = config.MaxPageSize },
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := validRequest()
			tt.modifyRequest(request)

			err := request.Validate()

			if tt.wantErr {
				require.NotNil(t, err, "expected error but got nil")
				if tt.wantErrContain != "" {
					assert.Contains(t, err.Message, tt.wantErrContain)
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSimpleSQLValidationShouldValidateSchemaNames(t *testing.T) {
	tests := []struct {
		name           string
		schema         string
		wantErr        bool
		wantErrContain string
	}{
		{
			name:    "valid_schema_alphanumeric",
			schema:  "MYSCHEMA",
			wantErr: false,
		},
		{
			name:    "valid_schema_with_underscore",
			schema:  "my_schema",
			wantErr: false,
		},
		{
			name:    "valid_schema_with_dollar",
			schema:  "schema$1",
			wantErr: false,
		},
		{
			name:    "valid_empty_schema",
			schema:  "",
			wantErr: false,
		},
		{
			name:           "invalid_schema_with_hyphen",
			schema:         "my-schema",
			wantErr:        true,
			wantErrContain: "SQL schema name validation failed",
		},
		{
			name:           "invalid_schema_with_space",
			schema:         "my schema",
			wantErr:        true,
			wantErrContain: "SQL schema name validation failed",
		},
		{
			name:           "invalid_schema_with_slash",
			schema:         "my/schema",
			wantErr:        true,
			wantErrContain: "SQL schema name validation failed",
		},
		{
			name:           "invalid_schema_with_semicolon_injection",
			schema:         "schema; DROP TABLE users;--",
			wantErr:        true,
			wantErrContain: "SQL schema name validation failed",
		},
		{
			name:           "invalid_schema_too_long",
			schema:         strings.Repeat("a", 129),
			wantErr:        true,
			wantErrContain: "SQL schema name validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - create a fully valid request, only varying schema
			request := &db2.Request{
				Username: "testuser",
				Password: "testpass",
				BaseURL:  "localhost",
				Database: "TESTDB",
				PageSize: 100,
				EntityConfig: framework.EntityConfig{
					ExternalId: "valid_table",
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "id", UniqueId: true},
					},
				},
				UniqueAttributeExternalID: "id",
				Schema:                    tt.schema,
			}

			// Act
			err := request.Validate()

			// Assert
			if tt.wantErr {
				require.NotNil(t, err, "expected error but got nil")
				assert.Contains(t, err.Message, tt.wantErrContain)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSimpleSQLValidationShouldValidateColumnNames(t *testing.T) {
	tests := []struct {
		name           string
		attributes     []*framework.AttributeConfig
		wantErr        bool
		wantErrContain string
	}{
		{
			name: "valid_column_alphanumeric",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "username"},
				{ExternalId: "email"},
			},
			wantErr: false,
		},
		{
			name: "valid_column_with_underscore",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "user_name"},
				{ExternalId: "created_at"},
			},
			wantErr: false,
		},
		{
			name: "valid_column_with_hyphen",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "user-name"},
			},
			wantErr: false,
		},
		{
			name: "valid_column_with_slash",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "user/name"},
			},
			wantErr: false,
		},
		{
			name: "valid_column_with_space",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "user name"},
			},
			wantErr: false,
		},
		{
			name: "valid_synthetic_id_skipped",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "id"},
			},
			wantErr: false,
		},
		{
			name: "invalid_column_with_semicolon_injection",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "col; DROP TABLE users;--"},
			},
			wantErr:        true,
			wantErrContain: "SQL column name validation failed",
		},
		{
			name: "invalid_column_with_quote_injection",
			attributes: []*framework.AttributeConfig{
				{ExternalId: "col\""},
			},
			wantErr:        true,
			wantErrContain: "SQL column name validation failed",
		},
		{
			name: "invalid_column_too_long",
			attributes: []*framework.AttributeConfig{
				{ExternalId: strings.Repeat("a", 129)},
			},
			wantErr:        true,
			wantErrContain: "SQL column name validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - ensure at least one attribute has UniqueId for validation
			attrs := tt.attributes
			hasUnique := false
			for _, attr := range attrs {
				if attr.UniqueId {
					hasUnique = true

					break
				}
			}
			if !hasUnique && len(attrs) > 0 {
				attrs[0].UniqueId = true
			}

			request := &db2.Request{
				Username: "testuser",
				Password: "testpass",
				BaseURL:  "localhost",
				Database: "TESTDB",
				PageSize: 100,
				EntityConfig: framework.EntityConfig{
					ExternalId: "valid_table",
					Attributes: attrs,
				},
				UniqueAttributeExternalID: "id",
			}

			// Act
			err := request.Validate()

			// Assert
			if tt.wantErr {
				require.NotNil(t, err, "expected error but got nil")
				assert.Contains(t, err.Message, tt.wantErrContain)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
