// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package identitynow_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/identitynow"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request      *identitynow.Request
		wantEndpoint string
		wantError    *framework.Error
	}{
		"accounts": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "v3",
				EntityExternalID: "accounts",
				PageSize:         100,
				Token:            "Bearer token",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](0),
				},
			},
			wantEndpoint: "https://sgnl-dev.api.identitynow-demo.com/v3/accounts?limit=100&offset=0",
		},
		"accounts_with_filter_and_sorters": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "v3",
				EntityExternalID: "accounts",
				PageSize:         100,
				Token:            "Bearer token",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](0),
				},
				Filter:  testutil.GenPtr[string](`name eq "John Doe"`),
				Sorters: testutil.GenPtr[string]("id,-modified"),
			},
			wantEndpoint: "https://sgnl-dev.api.identitynow-demo.com/v3/accounts?limit=100&offset=0&sorters=id,-modified&filters=name%20eq%20%22John%20Doe%22",
		},
		"accounts_with_filter_with_+_symbol": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "v3",
				EntityExternalID: "accounts",
				PageSize:         100,
				Token:            "Bearer token",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](0),
				},
				Filter: testutil.GenPtr[string](`name eq "John+ Doe"`),
			},
			wantEndpoint: "https://sgnl-dev.api.identitynow-demo.com/v3/accounts?limit=100&offset=0&filters=name%20eq%20%22John%2B%20Doe%22",
		},
		"accounts_nil_composite_cursor": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "v3",
				EntityExternalID: "accounts",
				PageSize:         100,
				Token:            "Bearer token",
				Cursor:           nil,
			},
			wantEndpoint: "",
			wantError: &framework.Error{
				Message: "Request cursor is unexpectedly nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"accounts_nil_cursor": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "v3",
				EntityExternalID: "accounts",
				PageSize:         100,
				Token:            "Bearer token",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: nil,
				},
			},
			wantEndpoint: "",
			wantError: &framework.Error{
				Message: "Request cursor is unexpectedly nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"entitlements": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "beta",
				EntityExternalID: "entitlements",
				PageSize:         100,
				Token:            "Bearer token",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](0),
				},
			},
			wantEndpoint: "https://sgnl-dev.api.identitynow-demo.com/beta/entitlements?limit=100&offset=0",
		},
		"account_entitlements": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "beta",
				EntityExternalID: identitynow.AccountEntitlements,
				PageSize:         100,
				Token:            "Bearer token",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:       testutil.GenPtr[int64](0),
					CollectionID: testutil.GenPtr[string]("13fa129a0bfe4f8db95fcc2a619bb826"),
				},
			},
			wantEndpoint: "https://sgnl-dev.api.identitynow-demo.com/beta/accounts/13fa129a0bfe4f8db95fcc2a619bb826/entitlements?limit=100&offset=0",
		},
		"account_entitlements_not_first_page": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "beta",
				EntityExternalID: identitynow.AccountEntitlements,
				PageSize:         100,
				Token:            "Bearer token",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:       testutil.GenPtr[int64](5),
					CollectionID: testutil.GenPtr[string]("13fa129a0bfe4f8db95fcc2a619bb826"),
				},
			},
			wantEndpoint: "https://sgnl-dev.api.identitynow-demo.com/beta/accounts/13fa129a0bfe4f8db95fcc2a619bb826/entitlements?limit=100&offset=5",
		},
		"account_entitlements_no_collection_id": {
			request: &identitynow.Request{
				BaseURL:          "https://sgnl-dev.api.identitynow-demo.com",
				APIVersion:       "beta",
				EntityExternalID: identitynow.AccountEntitlements,
				PageSize:         100,
				Token:            "Bearer token",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](0),
				},
			},
			wantEndpoint: "",
			wantError: &framework.Error{
				Message: "CollectionId field is unexpectedly nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"nil_request": {
			request:      nil,
			wantEndpoint: "",
			wantError:    nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotEndpoint, gotError := identitynow.ConstructEndpoint(tt.request)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if !reflect.DeepEqual(gotEndpoint, tt.wantEndpoint) {
				t.Errorf("gotEndpoint: %v, wantEndpoint: %v", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}
