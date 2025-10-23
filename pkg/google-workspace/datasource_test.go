// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst, dupword
package googleworkspace_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	googleworkspace "github.com/sgnl-ai/adapters/pkg/google-workspace"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		entityID     string
		body         []byte
		wantObjects  []map[string]interface{}
		wantNextLink *pagination.CompositeCursor[string]
		wantErr      *framework.Error
	}{
		"single_page": {
			entityID: "User",
			body:     []byte(`{"users": [{"id": "00ub0oNGTSWTBKOLGLNR","status": "ACTIVE"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}], "nextPageToken": "token4"}`),
			wantObjects: []map[string]interface{}{
				{"id": "00ub0oNGTSWTBKOLGLNR", "status": "ACTIVE"},
				{"id": "00ub0oNGTSWTBKOCHDKE", "status": "ACTIVE"},
			},
			wantNextLink: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("token4"),
			},
		},
		"invalid_object_structure": {
			entityID: "User",
			body:     []byte(`[{"id": "00ub0oNGTSWTBKOLGLNR","status": "ACTIVE"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal array into Go value of type googleworkspace.DatasourceResponse.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_objects": {
			entityID: "User",
			body:     []byte(`{"value": [{"00ub0oNGTSWTBKOLGLNR"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]}`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: `Failed to unmarshal the datasource response: invalid character '}' after object key.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextLink, gotErr := googleworkspace.ParseResponse(tt.body, &googleworkspace.Request{EntityExternalID: tt.entityID})

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextLink, tt.wantNextLink) {
				t.Errorf("gotNextLink: %v, wantNextLink: %v", gotNextLink, tt.wantNextLink)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetUsersPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	googleworkspaceClient := googleworkspace.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context      context.Context
		request      *googleworkspace.Request
		wantRes      *googleworkspace.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "User",
				PageSize:              1,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":         "admin#directory#user",
						"id":           "USER987654321",
						"etag":         "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/HFAgBzXOF0pogUtDLHSgZQ9lvls\"",
						"primaryEmail": "user1@sgnldemos.com",
						"name": map[string]any{
							"givenName":  "user1",
							"familyName": "user1",
							"fullName":   "user1 user1",
						},
						"isAdmin":                   true,
						"isDelegatedAdmin":          false,
						"lastLoginTime":             "2024-02-02T23:30:53.000Z",
						"creationTime":              "2024-02-02T23:30:06.000Z",
						"agreedToTerms":             true,
						"suspended":                 false,
						"archived":                  false,
						"changePasswordAtNextLogin": false,
						"ipWhitelisted":             false,
						"emails": []any{
							map[string]any{
								"address": "user1@sgnldemos.com",
								"primary": true,
							},
							map[string]any{
								"address": "user1@sgnldemos.com.test-google-a.com",
							},
						},
						"languages": []any{
							map[string]any{
								"languageCode": "en",
								"preference":   "preferred",
							},
						},
						"nonEditableAliases": []any{
							"user1@sgnldemos.com.test-google-a.com",
						},
						"customerId":                 "CUST123456",
						"orgUnitPath":                "/",
						"isMailboxSetup":             true,
						"isEnrolledIn2Sv":            false,
						"isEnforcedIn2Sv":            false,
						"includeInGlobalAddressList": true,
						"recoveryEmail":              "user1@sgnl.ai",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ=="),
				},
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "User",
					fields.FieldRequestPageSize:         int64(1),
				},
				{
					"level":                             "info",
					"msg":                               "Sending HTTP request to datasource",
					fields.FieldRequestEntityExternalID: "User",
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldURL:                     server.URL + "/admin/directory/v1/users?domain=sgnldemos.com&maxResults=1",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "User",
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(1),
					fields.FieldResponseNextCursor: map[string]any{
						"cursor": "Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ==",
					},
				},
			},
		},
		"second_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "User",
				PageSize:              1,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ=="),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":         "admin#directory#user",
						"id":           "102475661842232156723",
						"etag":         "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/zEsIEnG_BTp0IWmiowHAH5dGdhI\"",
						"primaryEmail": "user2@sgnldemos.com",
						"name": map[string]any{
							"givenName":  "user2",
							"familyName": "user2",
							"fullName":   "user2 user2",
						},
						"isAdmin":                   true,
						"isDelegatedAdmin":          false,
						"lastLoginTime":             "2024-04-18T23:55:53.000Z",
						"creationTime":              "2024-02-02T23:55:44.000Z",
						"agreedToTerms":             true,
						"suspended":                 false,
						"archived":                  false,
						"changePasswordAtNextLogin": false,
						"ipWhitelisted":             false,
						"emails": []any{
							map[string]any{
								"address": "user2@sgnldemos.com",
								"primary": true,
							},
							map[string]any{
								"address": "user2@sgnldemos.com.test-google-a.com",
							},
						},
						"languages": []any{
							map[string]any{
								"languageCode": "en",
								"preference":   "preferred",
							},
						},
						"nonEditableAliases": []any{
							"user2@sgnldemos.com.test-google-a.com",
						},
						"customerId":                 "CUST123456",
						"orgUnitPath":                "/",
						"isMailboxSetup":             true,
						"isEnrolledIn2Sv":            false,
						"isEnforcedIn2Sv":            false,
						"includeInGlobalAddressList": true,
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q0FFUzhRSUJrUHpWQUdiNG04Q1IxMklwSjFJemdWcFhvODg4ZFRjV1RSWTBWK2JoeDNXWUQyQUpLeFkzN1NuQWI4ZDFHUWszMmpESGpxR1I3Um5EcXd4V2REbi9Xc0NJTUYyVXVYR2xZcEgwdUVNRk5ZWCtlVlhzYTRXeXA3MFJ2NUxqT25vM1hCeUMzZ0wvdDRwUXZHa3pnd21QcnVZSm9udFEzMk9zMlcyaEhZOUJ5OGd0UzZmU3BZdHBpeE1uUUtOUWJ6ZlYrTUI0WjVnNFBYVVB4ZjRDZTJVc0pXQ05GSG1FZnYzQkMreU9BRWNYZWRkWCt2U3REMFR0ZjI0SElMY1Z3VHB3SHh3WURzbk84d0N5eTFsNDAwUFNVVGJrNW9BUkFwajJEL3dYcFo4bmRIa3FRdmRGK3Z3b0EwWSt6ZTB6Y3ZkcUlMd3pwVjIzL25GL0tIN2JPcTVqMWFSVVYrRDN0NE4zZzNxaU44clJ5c1dxMWhadkxyT1R2TGJjdWNVdVEwMVgwcHp5cTZlSG5vTUVWeWttMHJUZEFBNGdVOVFtVjUvZStBMWlhVlVrOEQyTVlpMmJJUUtaRnRoUk5jK3lhM1dScyt0ZDFSZWFHd09MVlNnQVNpQWZrQlZYaDJqbFVqblRvL3Y3OUZYWlp1TT0="),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "User",
				PageSize:              1,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q0FFUzhRSUJrUHpWQUdiNG04Q1IxMklwSjFJemdWcFhvODg4ZFRjV1RSWTBWK2JoeDNXWUQyQUpLeFkzN1NuQWI4ZDFHUWszMmpESGpxR1I3Um5EcXd4V2REbi9Xc0NJTUYyVXVYR2xZcEgwdUVNRk5ZWCtlVlhzYTRXeXA3MFJ2NUxqT25vM1hCeUMzZ0wvdDRwUXZHa3pnd21QcnVZSm9udFEzMk9zMlcyaEhZOUJ5OGd0UzZmU3BZdHBpeE1uUUtOUWJ6ZlYrTUI0WjVnNFBYVVB4ZjRDZTJVc0pXQ05GSG1FZnYzQkMreU9BRWNYZWRkWCt2U3REMFR0ZjI0SElMY1Z3VHB3SHh3WURzbk84d0N5eTFsNDAwUFNVVGJrNW9BUkFwajJEL3dYcFo4bmRIa3FRdmRGK3Z3b0EwWSt6ZTB6Y3ZkcUlMd3pwVjIzL25GL0tIN2JPcTVqMWFSVVYrRDN0NE4zZzNxaU44clJ5c1dxMWhadkxyT1R2TGJjdWNVdVEwMVgwcHp5cTZlSG5vTUVWeWttMHJUZEFBNGdVOVFtVjUvZStBMWlhVlVrOEQyTVlpMmJJUUtaRnRoUk5jK3lhM1dScyt0ZDFSZWFHd09MVlNnQVNpQWZrQlZYaDJqbFVqblRvL3Y3OUZYWlp1TT0="),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":         "admin#directory#user",
						"id":           "114211199695002816249",
						"etag":         "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/0VBDYIsSDmoa_zjvxKsXSmzSk2I\"",
						"primaryEmail": "sor-dev@sgnldemos.com",
						"name": map[string]any{
							"givenName":  "SoR",
							"familyName": "Development",
							"fullName":   "SoR Development",
						},
						"isAdmin":                   false,
						"isDelegatedAdmin":          true,
						"lastLoginTime":             "2024-04-20T05:56:11.000Z",
						"creationTime":              "2024-04-18T00:30:35.000Z",
						"agreedToTerms":             true,
						"suspended":                 false,
						"archived":                  false,
						"changePasswordAtNextLogin": false,
						"ipWhitelisted":             false,
						"emails": []any{
							map[string]any{
								"address": "sgnl-demos@sgnl.ai",
								"type":    "work",
							},
							map[string]any{
								"address": "sor-dev@sgnldemos.com",
								"primary": true,
							},
							map[string]any{
								"address": "sor-dev@sgnldemos.com.test-google-a.com",
							},
						},
						"languages": []any{
							map[string]any{
								"languageCode": "en",
								"preference":   "preferred",
							},
						},
						"nonEditableAliases": []any{
							"sor-dev@sgnldemos.com.test-google-a.com",
						},
						"customerId":                 "CUST123456",
						"orgUnitPath":                "/serviceAccounts",
						"isMailboxSetup":             true,
						"isEnrolledIn2Sv":            false,
						"isEnforcedIn2Sv":            false,
						"includeInGlobalAddressList": true,
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.context)

			gotRes, gotErr := googleworkspaceClient.GetPage(ctxWithLogger, tt.request)

			if diff := cmp.Diff(gotRes, tt.wantRes); diff != "" {
				t.Errorf("gotRes: %v, wantRes: %v\n%s", gotRes, tt.wantRes, diff)
			}

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}

func TestGetGroupsPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	googleworkspaceClient := googleworkspace.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *googleworkspace.Request
		wantRes *googleworkspace.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Group",
				PageSize:              1,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":               "admin#directory#group",
						"id":                 "01qoc8b13vgdlqb",
						"etag":               "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/jegUmo9HAPpclSc2UeC_v9oHm2E\"",
						"email":              "emptygroup@sgnldemos.com",
						"name":               "Empty Group",
						"directMembersCount": "0",
						"description":        "",
						"adminCreated":       true,
						"nonEditableAliases": []any{
							"emptygroup@sgnldemos.com.test-google-a.com",
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"),
				},
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Group",
				PageSize:              1,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":               "admin#directory#group",
						"id":                 "048pi1tg0qf1f8g",
						"etag":               "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/n5phxV9NHrbJaK4QhiAp3pDR21k\"",
						"email":              "group2@sgnldemos.com",
						"name":               "Group2",
						"directMembersCount": "3",
						"description":        "",
						"adminCreated":       true,
						"nonEditableAliases": []any{
							"group2@sgnldemos.com.test-google-a.com",
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
				},
			},
			wantErr: nil,
		},
		"third_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Group",
				PageSize:              1,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":               "admin#directory#group",
						"id":                 "030j0zll41obyxg",
						"etag":               "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/knD5dVuTcKRMQdMc1Qr7e37h_n4\"",
						"email":              "hello@sgnldemos.com",
						"name":               "SGNLDemos",
						"directMembersCount": "2",
						"description":        "",
						"adminCreated":       true,
						"nonEditableAliases": []any{
							"hello@sgnldemos.com.test-google-a.com",
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Group",
				PageSize:              1,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects:    nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := googleworkspaceClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(gotRes, tt.wantRes); diff != "" {
				t.Errorf("Mismatch (-want +got):\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetMembersPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	googleworkspaceClient := googleworkspace.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *googleworkspace.Request
		wantRes *googleworkspace.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Member",
				PageSize:              2,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects:    nil,
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor:           nil,
					CollectionCursor: testutil.GenPtr("Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"),
				},
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Member",
				PageSize:              2,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           nil,
					CollectionCursor: testutil.GenPtr("Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":     "admin#directory#member",
						"etag":     "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/1pVaKo8kmZjl-_gkcKlCYvYBPjA\"",
						"id":       "USER987654321",
						"groupId":  "048pi1tg0qf1f8g",
						"uniqueId": "048pi1tg0qf1f8g-USER987654321",
						"email":    "user1@sgnldemos.com",
						"role":     "MEMBER",
						"type":     "USER",
						"status":   "ACTIVE",
					},
					{
						"kind":     "admin#directory#member",
						"etag":     "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/4j9E51Yhqt_d3t65g3OEza89nps\"",
						"id":       "102475661842232156723",
						"groupId":  "048pi1tg0qf1f8g",
						"uniqueId": "048pi1tg0qf1f8g-102475661842232156723",
						"email":    "user2@sgnldemos.com",
						"role":     "OWNER",
						"type":     "USER",
						"status":   "ACTIVE",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr("CjRJaDhLSFFqSXZNaWZ6Z0VTRW0xaGNtTkFjMmR1YkdSbGJXOXpMbU52YlJnQllKeUppTThFIh8KHQjIvMifzgESEm1hcmNAc2dubGRlbW9zLmNvbRgBYJyJiM8E"),
					CollectionID:     testutil.GenPtr("048pi1tg0qf1f8g"),
					CollectionCursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
				},
			},
			wantErr: nil,
		},
		"third_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Member",
				PageSize:              2,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr("CjRJaDhLSFFqSXZNaWZ6Z0VTRW0xaGNtTkFjMmR1YkdSbGJXOXpMbU52YlJnQllKeUppTThFIh8KHQjIvMifzgESEm1hcmNAc2dubGRlbW9zLmNvbRgBYJyJiM8E"),
					CollectionID:     testutil.GenPtr("048pi1tg0qf1f8g"),
					CollectionCursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":     "admin#directory#member",
						"etag":     "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/d566bqA1WOuRhrsagV1MWF5o2f8\"",
						"id":       "114211199695002816249",
						"groupId":  "048pi1tg0qf1f8g",
						"uniqueId": "048pi1tg0qf1f8g-114211199695002816249",
						"email":    "sor-dev@sgnldemos.com",
						"role":     "OWNER",
						"type":     "USER",
						"status":   "ACTIVE",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor:           nil,
					CollectionID:     nil,
					CollectionCursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
				},
			},
			wantErr: nil,
		},
		"fourth_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Member",
				PageSize:              2,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           nil,
					CollectionID:     nil,
					CollectionCursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"kind":     "admin#directory#member",
						"etag":     "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/Z1pIME_BX7RFRQE_z0f1Gfewx70\"",
						"id":       "048pi1tg0qf1f8g",
						"groupId":  "030j0zll41obyxg",
						"uniqueId": "030j0zll41obyxg-048pi1tg0qf1f8g",
						"email":    "group2@sgnldemos.com",
						"role":     "MEMBER",
						"type":     "GROUP",
						"status":   "ACTIVE",
					},
					{
						"kind":     "admin#directory#member",
						"etag":     "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/J-EKKLlSn8MJS_aGpigqCNOJvYs\"",
						"id":       "102475661842232156723",
						"groupId":  "030j0zll41obyxg",
						"uniqueId": "030j0zll41obyxg-102475661842232156723",
						"email":    "user2@sgnldemos.com",
						"role":     "OWNER",
						"type":     "USER",
						"status":   "ACTIVE",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor:           nil,
					CollectionID:     nil,
					CollectionCursor: testutil.GenPtr("Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &googleworkspace.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Member",
				PageSize:              2,
				RequestTimeoutSeconds: 5,
				Domain:                testutil.GenPtr("sgnldemos.com"),
				APIVersion:            "v1",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           nil,
					CollectionID:     nil,
					CollectionCursor: testutil.GenPtr("Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="),
				},
			},
			wantRes: &googleworkspace.Response{
				StatusCode: http.StatusOK,
				Objects:    nil,
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := googleworkspaceClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Mismatch (-want +got):\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
