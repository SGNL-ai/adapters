// Copyright 2025 SGNL.ai, Inc.
package rootly_test

import (
	"strings"
	"testing"

	rootly_adapter "github.com/sgnl-ai/adapters/pkg/rootly"
)

func TestConstructEndpoint(t *testing.T) {
	tests := map[string]struct {
		request     *rootly_adapter.Request
		expectedURL string
	}{
		"basic_incidents_request": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           10,
			},
			expectedURL: "https://api.rootly.com/v1/incidents?page%5Bnumber%5D=1&page%5Bsize%5D=10",
		},
		"basic_users_request": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "users",
				PageSize:           25,
			},
			expectedURL: "https://api.rootly.com/v1/users?page%5Bnumber%5D=1&page%5Bsize%5D=25",
		},
		"basic_teams_request": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "teams",
				PageSize:           50,
			},
			expectedURL: "https://api.rootly.com/v1/teams?page%5Bnumber%5D=1&page%5Bsize%5D=50",
		},
		"request_with_cursor": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           20,
				Cursor:             stringPtr("3"),
			},
			expectedURL: "https://api.rootly.com/v1/incidents?page%5Bnumber%5D=3&page%5Bsize%5D=20",
		},
		"request_with_empty_cursor": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "users",
				PageSize:           15,
				Cursor:             stringPtr(""),
			},
			expectedURL: "https://api.rootly.com/v1/users?page%5Bnumber%5D=1&page%5Bsize%5D=15",
		},
		"request_with_simple_filter": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           30,
				Filter:             "status=open",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bstatus%5D=open&page%5Bnumber%5D=1&page%5Bsize%5D=30",
		},
		"request_with_multiple_filters": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           40,
				Filter:             "status=started&severity=high",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bseverity%5D=high&filter%5Bstatus%5D=started&page%5Bnumber%5D=1&page%5Bsize%5D=40",
		},
		"request_with_comma_separated_filter_values": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           100,
				Filter:             "status=started,mitigated,resolved",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bstatus%5D=started%2Cmitigated%2Cresolved&page%5Bnumber%5D=1&page%5Bsize%5D=100",
		},
		"request_with_complex_filters": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           50,
				Filter:             "severity=major,minor&status=started,mitigated&kind=normal",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bkind%5D=normal&filter%5Bseverity%5D=major%2Cminor&filter%5Bstatus%5D=started%2Cmitigated&page%5Bnumber%5D=1&page%5Bsize%5D=50",
		},
		"request_with_cursor_and_filters": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           25,
				Cursor:             stringPtr("2"),
				Filter:             "status=open&severity=critical",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bseverity%5D=critical&filter%5Bstatus%5D=open&page%5Bnumber%5D=2&page%5Bsize%5D=25",
		},
		"request_with_special_characters_in_filter": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "users",
				PageSize:           20,
				Filter:             "email=test@example.com&name=John Doe",
			},
			expectedURL: "https://api.rootly.com/v1/users?filter%5Bemail%5D=test%40example.com&filter%5Bname%5D=John+Doe&page%5Bnumber%5D=1&page%5Bsize%5D=20",
		},
		"request_with_zero_page_size": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "teams",
				PageSize:           0,
			},
			expectedURL: "https://api.rootly.com/v1/teams?page%5Bnumber%5D=1",
		},
		"request_with_large_page_size": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           1000,
			},
			expectedURL: "https://api.rootly.com/v1/incidents?page%5Bnumber%5D=1&page%5Bsize%5D=1000",
		},
		"request_with_empty_filter": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "users",
				PageSize:           15,
				Filter:             "",
			},
			expectedURL: "https://api.rootly.com/v1/users?page%5Bnumber%5D=1&page%5Bsize%5D=15",
		},
		"request_with_malformed_filter": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           30,
				Filter:             "invalid-filter-format",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Binvalid-filter-format%5D=&page%5Bnumber%5D=1&page%5Bsize%5D=30",
		},
		"request_with_different_entity_types": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "services",
				PageSize:           35,
				Filter:             "active=true",
			},
			expectedURL: "https://api.rootly.com/v1/services?filter%5Bactive%5D=true&page%5Bnumber%5D=1&page%5Bsize%5D=35",
		},
		"request_with_nested_path": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incident_types",
				PageSize:           10,
			},
			expectedURL: "https://api.rootly.com/v1/incident_types?page%5Bnumber%5D=1&page%5Bsize%5D=10",
		},
		"request_with_high_page_number": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           10,
				Cursor:             stringPtr("99"),
			},
			expectedURL: "https://api.rootly.com/v1/incidents?page%5Bnumber%5D=99&page%5Bsize%5D=10",
		},
		"request_with_boolean_filter": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           20,
				Filter:             "private=false&resolved=true",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bprivate%5D=false&filter%5Bresolved%5D=true&page%5Bnumber%5D=1&page%5Bsize%5D=20",
		},
		"request_with_url_encoded_characters": {
			request: &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           15,
				Filter:             "title=Database%20Error&status=open",
			},
			expectedURL: "https://api.rootly.com/v1/incidents?filter%5Bstatus%5D=open&filter%5Btitle%5D=Database+Error&page%5Bnumber%5D=1&page%5Bsize%5D=15",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actualURL := rootly_adapter.ConstructEndpoint(tt.request)
			if actualURL != tt.expectedURL {
				t.Errorf("ConstructEndpoint() = %v, want %v", actualURL, tt.expectedURL)
			}
		})
	}
}

func TestConstructEndpointFilterTransformation(t *testing.T) {
	tests := map[string]struct {
		name           string
		inputFilter    string
		expectedParams map[string]string
	}{
		"single_filter": {
			name:        "single status filter",
			inputFilter: "status=open",
			expectedParams: map[string]string{
				"filter%5Bstatus%5D": "open",
			},
		},
		"multiple_filters": {
			name:        "multiple separate filters",
			inputFilter: "status=started&severity=high&kind=normal",
			expectedParams: map[string]string{
				"filter%5Bstatus%5D":   "started",
				"filter%5Bseverity%5D": "high",
				"filter%5Bkind%5D":     "normal",
			},
		},
		"comma_separated_values": {
			name:        "comma separated filter values",
			inputFilter: "status=started,mitigated,resolved",
			expectedParams: map[string]string{
				"filter%5Bstatus%5D": "started%2Cmitigated%2Cresolved",
			},
		},
		"mixed_format": {
			name:        "mixed comma separated and separate filters",
			inputFilter: "status=started,mitigated&severity=high,medium&private=false",
			expectedParams: map[string]string{
				"filter%5Bstatus%5D":   "started%2Cmitigated",
				"filter%5Bseverity%5D": "high%2Cmedium",
				"filter%5Bprivate%5D":  "false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &rootly_adapter.Request{
				BaseURL:            "https://api.rootly.com/v1",
				EntityExternalName: "incidents",
				PageSize:           10,
				Filter:             tt.inputFilter,
			}

			actualURL := rootly_adapter.ConstructEndpoint(request)

			// Parse the URL to check parameters
			for expectedParam := range tt.expectedParams {
				// Check if the URL contains the expected parameter
				if !contains(actualURL, expectedParam) {
					t.Errorf("ConstructEndpoint() URL %v does not contain expected parameter %v", actualURL, expectedParam)
				}
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	// Convert both to lowercase for case-insensitive comparison
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)

	return strings.Contains(sLower, substrLower)
}
