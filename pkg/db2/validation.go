// Copyright 2026 SGNL.ai, Inc.

package db2

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/validation"
)

// NewRequestFromConfig validates the framework request and constructs an internal Request.
// It performs nil checks, constructs the Request, and validates all fields.
// Returns a fully validated Request or an error.
func NewRequestFromConfig(request *framework.Request[Config]) (*Request, *framework.Error) {
	// Nil checks to safely access nested fields
	if request == nil {
		return nil, &framework.Error{
			Message: "DB2 request is invalid: request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if request.Auth == nil {
		return nil, &framework.Error{
			Message: "DB2 auth is invalid: auth is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_AUTH,
		}
	}

	if request.Auth.Basic == nil {
		return nil, &framework.Error{
			Message: "DB2 auth is invalid: Basic authentication is required.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_AUTH,
		}
	}

	if request.Config == nil {
		return nil, &framework.Error{
			Message: "DB2 config is invalid: request contains no config.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Construct the internal Request
	req := &Request{
		Username:     request.Auth.Basic.Username,
		Password:     request.Auth.Basic.Password,
		BaseURL:      request.Address,
		PageSize:     request.PageSize,
		EntityConfig: request.Entity,
		Database:     request.Config.Database,
		Schema:       request.Config.Schema,
		ConfigStruct: request.Config,
	}

	if request.Cursor != "" {
		req.Cursor = &request.Cursor
	}

	for _, attribute := range request.Entity.Attributes {
		if attribute.UniqueId {
			req.UniqueAttributeExternalID = attribute.ExternalId

			break
		}
	}

	if request.Config.Filters != nil {
		if curFilter, ok := request.Config.Filters[request.Entity.ExternalId]; ok {
			req.Filter = &curFilter
		}
	}

	return req, nil
}

// Validate performs deep validation on the Request to ensure all fields are properly populated.
// Called after the Request is constructed.
func (r *Request) Validate() *framework.Error {
	if r == nil {
		return &framework.Error{
			Message: "DB2 request is invalid: request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if r.Username == "" || r.Password == "" {
		return &framework.Error{
			Message: "DB2 auth is invalid: username and password are required.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_AUTH,
		}
	}

	if r.BaseURL == "" {
		return &framework.Error{
			Message: "DB2 config is invalid: address (hostname) is required.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if err := validateAddressPort(r.BaseURL); err != nil {
		return err
	}

	if r.Database == "" {
		return &framework.Error{
			Message: "DB2 config is invalid: database is not set.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if !hasUniqueAttribute(r.EntityConfig.Attributes) {
		return &framework.Error{
			Message: "DB2 entity config is invalid: no unique attribute defined.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if r.PageSize <= 0 || r.PageSize > config.MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf("DB2 request is invalid: page size must be between 1 and %d.", config.MaxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	// Validate SQL identifiers to prevent SQL injection
	if err := r.validateSQLIdentifiers(); err != nil {
		return err
	}

	return nil
}

// validateSQLIdentifiers validates SQL identifiers to prevent SQL injection.
// It validates table names, schema names (strict), and column names (permissive).
func (r *Request) validateSQLIdentifiers() *framework.Error {
	if !validation.IsValidSQLIdentifier(r.EntityConfig.ExternalId) {
		return &framework.Error{
			Message: "SQL table name validation failed: unsupported characters found or length is not in range 1-128.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if !validation.IsValidSQLIdentifier(r.UniqueAttributeExternalID) {
		return &framework.Error{
			Message: "SQL unique attribute validation failed: unsupported characters found or length is not in range 1-128.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate schema name if provided. Schema names use strict validation
	// since they are part of the table qualifier (SCHEMA.TABLE).
	if r.Schema != "" {
		if !validation.IsValidSQLIdentifier(r.Schema) {
			return &framework.Error{
				Message: "SQL schema name validation failed: unsupported characters found or length is not in range 1-128.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}
	}

	// Validate column names. Uses permissive validation that allows /, -, and space
	// since quoteIdentifier() safely handles these by wrapping in double quotes.
	for _, attr := range r.EntityConfig.Attributes {
		// Skip the synthetic "id" attribute used for composite key generation
		if attr.ExternalId == "id" {
			continue
		}

		if !validation.IsValidColumnIdentifier(attr.ExternalId) {
			return &framework.Error{
				Message: fmt.Sprintf(
					"SQL column name validation failed for '%s': "+
						"unsupported characters found or length is not in range 1-128.",
					attr.ExternalId),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}
	}

	return nil
}

// validateAddressPort validates the port format in the address if specified.
func validateAddressPort(address string) *framework.Error {
	if !strings.Contains(address, ":") {
		return nil
	}

	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("DB2 config is invalid: invalid address format: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if host == "" {
		return &framework.Error{
			Message: "DB2 config is invalid: hostname is empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return &framework.Error{
			Message: fmt.Sprintf("DB2 config is invalid: invalid port '%s'. Must be between 1 and 65535.", portStr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	return nil
}

// hasUniqueAttribute checks if the attributes contain a unique ID attribute.
func hasUniqueAttribute(attributes []*framework.AttributeConfig) bool {
	for _, attr := range attributes {
		if attr != nil && attr.UniqueId {
			return true
		}
	}

	return false
}
