// Copyright 2025 SGNL.ai, Inc.
package crowdstrike_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/crowdstrike"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructRESTEndpoint(t *testing.T) {
	tests := map[string]struct {
		request   *crowdstrike.Request
		path      string
		wantURL   *string
		wantError *framework.Error
	}{
		"nil_request": {
			request: nil,
			path:    "devices/queries/devices-scroll/v1",
			wantError: &framework.Error{
				Message: "Request is nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"empty_path": {
			request: &crowdstrike.Request{
				BaseURL:  "https://api.crowdstrike.com",
				PageSize: 100,
				RESTCursor: &pagination.CompositeCursor[string]{
					Cursor: nil,
				},
			},
			path: "",
			wantError: &framework.Error{
				Message: "The path to fetch the entity from is nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"valid_request_with_cursor_string_offset": {
			request: &crowdstrike.Request{
				BaseURL:          "https://api.crowdstrike.com",
				PageSize:         100,
				EntityExternalID: crowdstrike.Device, // has UseIntCursor = false
				RESTCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("some-string-cursor"),
				},
			},
			path: "devices/queries/devices-scroll/v1",
			wantURL: testutil.GenPtr(
				"https://api.crowdstrike.com/devices/queries/devices-scroll/v1?limit=100&offset=some-string-cursor",
			),
		},
		"valid_request_with_cursor_string_offset_and_filter": {
			request: &crowdstrike.Request{
				BaseURL:          "https://api.crowdstrike.com",
				PageSize:         100,
				EntityExternalID: crowdstrike.Device,
				RESTCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("some-string-cursor"),
				},
				Filter: testutil.GenPtr(`platform_name:'Windows'`),
			},
			path: "devices/queries/devices-scroll/v1",
			wantURL: testutil.GenPtr(
				"https://api.crowdstrike.com/devices/queries/devices-scroll/v1?limit=100" +
					"&offset=some-string-cursor&filter=platform_name%3A%27Windows%27",
			),
		},
		"valid_request_without_cursor_string_offset": {
			request: &crowdstrike.Request{
				BaseURL:          "https://api.crowdstrike.com",
				PageSize:         100,
				EntityExternalID: crowdstrike.Device,
				RESTCursor: &pagination.CompositeCursor[string]{
					Cursor: nil,
				},
			},
			path:    "devices/queries/devices-scroll/v1",
			wantURL: testutil.GenPtr("https://api.crowdstrike.com/devices/queries/devices-scroll/v1?limit=100"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotURL, gotError := crowdstrike.ConstructRESTEndpoint(tt.request, tt.path)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if diff := cmp.Diff(gotURL, tt.wantURL); diff != "" {
				t.Errorf("gotURL: %v, wantURL: %v", *gotURL, *tt.wantURL)
			}
		})
	}
}
