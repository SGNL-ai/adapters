// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst, dupword
package googleworkspace_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	googleworkspace "github.com/sgnl-ai/adapters/pkg/google-workspace"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestAdapterGetPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := googleworkspace.NewAdapter(&googleworkspace.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[googleworkspace.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                        "USER987654321",
							"primaryEmail":              "user1@sgnldemos.com",
							"$.name.fullName":           "user1 user1",
							"isAdmin":                   true,
							"isDelegatedAdmin":          false,
							"creationTime":              time.Date(2024, 2, 2, 23, 30, 6, 0, time.UTC),
							"changePasswordAtNextLogin": false,
							"emails": []framework.Object{
								map[string]any{
									"address": "user1@sgnldemos.com",
									"primary": true,
								},
								map[string]any{
									"address": "user1@sgnldemos.com.test-google-a.com",
								},
							},
							"nonEditableAliases": []string{
								"user1@sgnldemos.com.test-google-a.com",
							},
							"customerId": "CUST123456",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJRMEZGVXpsM1NVSnJVSHBXUVVoVFdXaFBhMlZRYzBkYWJHMXdWbEJ4ZGxkME4yTnRlblpQVERWTU0yOVpOVFpLSzJkNk1FNU5jemRWTlRGcVdYUlVablZXY2xRemQwYzNjVlpXWms4clIwTTBWMUZoVFdweVYyUjNjblp6ZFM5bGQzQlVhbU5rU20xQk1taEVaakZJY0d3NVFVaHlXRzh5YXpWMVdqSk9OemRuTkhoYU1qVTVSMnBTZW5JdmVXcE9lRFpPVDA5M2JrMU5ZMjB6WVc5RFZpdEpNbWMwV1VsRWFqUkRZVzkzUkRrNFQwTTJNak00YUdseGRVRTVVMDVCVDJsbmVtdGxZbnBKU1hvclZHVTVTVUZTWXpSQmQyTTVPVWxKZG5nclEzWXpNa1JMZWtGVFdHcHVTbFIyYUVScVRITXlObnB4TlhaTFpHOXlla3hLU1VaRGRXRlVkamQxT1ZwSlpXTk1XREJYZWpod1ZXZFVNRE5TYlZoVVN6SjFlR294ZVZCYWFuSlhhbGtyVXpaV1dsSjJNblpKZVVVck0zbGtVM2hxVUhGc0wydzFTMEZ3VUc5bFoyMVRiVWgwTUhaaVIzWXZWMDRyYVVReGN6aDNObkZUTm5KQ05GaExTMjB3WkV4VlIyMXpTa05KUjA5VmMybFhRamhvVmpGa1pVa3JabWxETlc4NWFXNWxVV0ZqWWtac1RYVmFjUzkwZWpsV00wSkVibmsxVWs1bVVEazFjbTFFZDJoc1JYUmpNVVpUZUdOb2NVMDJRVXcxT0ZoSFZGVkRUVFZuT1daWllWWkpOMEU0VEhFNVdITklURTlXZDNneVExRlBiV2s1Y1VSTVl5OXBkWHAxUzFOSmVFOHhiVGxGYVhoeVJsRkZaREJSUFE9PSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ=="),
			},
		},
		"valid_request_no_https_prefix": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                        "USER987654321",
							"primaryEmail":              "user1@sgnldemos.com",
							"$.name.fullName":           "user1 user1",
							"isAdmin":                   true,
							"isDelegatedAdmin":          false,
							"creationTime":              time.Date(2024, 2, 2, 23, 30, 6, 0, time.UTC),
							"changePasswordAtNextLogin": false,
							"emails": []framework.Object{
								map[string]any{
									"address": "user1@sgnldemos.com",
									"primary": true,
								},
								map[string]any{
									"address": "user1@sgnldemos.com.test-google-a.com",
								},
							},
							"nonEditableAliases": []string{
								"user1@sgnldemos.com.test-google-a.com",
							},
							"customerId": "CUST123456",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJRMEZGVXpsM1NVSnJVSHBXUVVoVFdXaFBhMlZRYzBkYWJHMXdWbEJ4ZGxkME4yTnRlblpQVERWTU0yOVpOVFpLSzJkNk1FNU5jemRWTlRGcVdYUlVablZXY2xRemQwYzNjVlpXWms4clIwTTBWMUZoVFdweVYyUjNjblp6ZFM5bGQzQlVhbU5rU20xQk1taEVaakZJY0d3NVFVaHlXRzh5YXpWMVdqSk9OemRuTkhoYU1qVTVSMnBTZW5JdmVXcE9lRFpPVDA5M2JrMU5ZMjB6WVc5RFZpdEpNbWMwV1VsRWFqUkRZVzkzUkRrNFQwTTJNak00YUdseGRVRTVVMDVCVDJsbmVtdGxZbnBKU1hvclZHVTVTVUZTWXpSQmQyTTVPVWxKZG5nclEzWXpNa1JMZWtGVFdHcHVTbFIyYUVScVRITXlObnB4TlhaTFpHOXlla3hLU1VaRGRXRlVkamQxT1ZwSlpXTk1XREJYZWpod1ZXZFVNRE5TYlZoVVN6SjFlR294ZVZCYWFuSlhhbGtyVXpaV1dsSjJNblpKZVVVck0zbGtVM2hxVUhGc0wydzFTMEZ3VUc5bFoyMVRiVWgwTUhaaVIzWXZWMDRyYVVReGN6aDNObkZUTm5KQ05GaExTMjB3WkV4VlIyMXpTa05KUjA5VmMybFhRamhvVmpGa1pVa3JabWxETlc4NWFXNWxVV0ZqWWtac1RYVmFjUzkwZWpsV00wSkVibmsxVWs1bVVEazFjbTFFZDJoc1JYUmpNVVpUZUdOb2NVMDJRVXcxT0ZoSFZGVkRUVFZuT1daWllWWkpOMEU0VEhFNVdITklURTlXZDNneVExRlBiV2s1Y1VSTVl5OXBkWHAxUzFOSmVFOHhiVGxGYVhoeVJsRkZaREJSUFE9PSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ=="),
			},
		},
		"invalid_request_invalid_api_version": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1.0",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Google Workspace adapter config is invalid: apiVersion is not supported: v1.0.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"invalid_request_http_prefix": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: "http://" + strings.TrimPrefix(server.URL, "https://"),
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: `Scheme "http" is not supported.`,
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
		"valid_request_with_cursor": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ=="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                        "102475661842232156723",
							"primaryEmail":              "user2@sgnldemos.com",
							"$.name.fullName":           "user2 user2",
							"isAdmin":                   true,
							"isDelegatedAdmin":          false,
							"creationTime":              time.Date(2024, 2, 2, 23, 55, 44, 0, time.UTC),
							"changePasswordAtNextLogin": false,
							"emails": []framework.Object{
								map[string]any{
									"address": "user2@sgnldemos.com",
									"primary": true,
								},
								map[string]any{
									"address": "user2@sgnldemos.com.test-google-a.com",
								},
							},
							"nonEditableAliases": []string{
								"user2@sgnldemos.com.test-google-a.com",
							},
							"customerId": "CUST123456",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJRMEZGVXpoUlNVSnJVSHBXUVVkaU5HMDRRMUl4TWtsd1NqRkplbWRXY0Zodk9EZzRaRlJqVjFSU1dUQldLMkpvZUROWFdVUXlRVXBMZUZrek4xTnVRV0k0WkRGSFVXc3pNbXBFU0dweFIxSTNVbTVFY1hkNFYyUkViaTlYYzBOSlRVWXlWWFZZUjJ4WmNFZ3dkVVZOUms1WldDdGxWbGh6WVRSWGVYQTNNRkoyTlV4cVQyNXZNMWhDZVVNelowd3ZkRFJ3VVhaSGEzcG5kMjFRY25WWlNtOXVkRkV6TWs5ek1sY3lhRWhaT1VKNU9HZDBVelptVTNCWmRIQnBlRTF1VVV0T1VXSjZabFlyVFVJMFdqVm5ORkJZVlZCNFpqUkRaVEpWYzBwWFEwNUdTRzFGWm5ZelFrTXJlVTlCUldOWVpXUmtXQ3QyVTNSRU1GUjBaakkwU0VsTVkxWjNWSEIzU0hoM1dVUnpiazg0ZDBONWVURnNOREF3VUZOVlZHSnJOVzlCVWtGd2FqSkVMM2RZY0ZvNGJtUklhM0ZSZG1SR0szWjNiMEV3V1N0NlpUQjZZM1prY1VsTWQzcHdWakl6TDI1R0wwdElOMkpQY1RWcU1XRlNWVllyUkROME5FNHpaek54YVU0NGNsSjVjMWR4TVdoYWRreHlUMVIyVEdKamRXTlZkVkV3TVZnd2NIcDVjVFpsU0c1dlRVVldlV3R0TUhKVVpFRkJOR2RWT1ZGdFZqVXZaU3RCTVdsaFZsVnJPRVF5VFZscE1tSkpVVXRhUm5Sb1VrNWpLM2xoTTFkU2N5dDBaREZTWldGSGQwOU1WbE5uUVZOcFFXWnJRbFpZYURKcWJGVnFibFJ2TDNZM09VWllXbHAxVFQwPSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q0FFUzhRSUJrUHpWQUdiNG04Q1IxMklwSjFJemdWcFhvODg4ZFRjV1RSWTBWK2JoeDNXWUQyQUpLeFkzN1NuQWI4ZDFHUWszMmpESGpxR1I3Um5EcXd4V2REbi9Xc0NJTUYyVXVYR2xZcEgwdUVNRk5ZWCtlVlhzYTRXeXA3MFJ2NUxqT25vM1hCeUMzZ0wvdDRwUXZHa3pnd21QcnVZSm9udFEzMk9zMlcyaEhZOUJ5OGd0UzZmU3BZdHBpeE1uUUtOUWJ6ZlYrTUI0WjVnNFBYVVB4ZjRDZTJVc0pXQ05GSG1FZnYzQkMreU9BRWNYZWRkWCt2U3REMFR0ZjI0SElMY1Z3VHB3SHh3WURzbk84d0N5eTFsNDAwUFNVVGJrNW9BUkFwajJEL3dYcFo4bmRIa3FRdmRGK3Z3b0EwWSt6ZTB6Y3ZkcUlMd3pwVjIzL25GL0tIN2JPcTVqMWFSVVYrRDN0NE4zZzNxaU44clJ5c1dxMWhadkxyT1R2TGJjdWNVdVEwMVgwcHp5cTZlSG5vTUVWeWttMHJUZEFBNGdVOVFtVjUvZStBMWlhVlVrOEQyTVlpMmJJUUtaRnRoUk5jK3lhM1dScyt0ZDFSZWFHd09MVlNnQVNpQWZrQlZYaDJqbFVqblRvL3Y3OUZYWlp1TT0="),
			},
		},
		"invalid_request_invalid_url": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL + "/invalid",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
				t.Errorf("GetPage() mismatch (-want +got):\n%s", diff)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetUserPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := googleworkspace.NewAdapter(&googleworkspace.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[googleworkspace.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                        "USER987654321",
							"primaryEmail":              "user1@sgnldemos.com",
							"$.name.fullName":           "user1 user1",
							"isAdmin":                   true,
							"isDelegatedAdmin":          false,
							"creationTime":              time.Date(2024, 2, 2, 23, 30, 6, 0, time.UTC),
							"changePasswordAtNextLogin": false,
							"emails": []framework.Object{
								map[string]any{
									"address": "user1@sgnldemos.com",
									"primary": true,
								},
								map[string]any{
									"address": "user1@sgnldemos.com.test-google-a.com",
								},
							},
							"nonEditableAliases": []string{
								"user1@sgnldemos.com.test-google-a.com",
							},
							"customerId": "CUST123456",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJRMEZGVXpsM1NVSnJVSHBXUVVoVFdXaFBhMlZRYzBkYWJHMXdWbEJ4ZGxkME4yTnRlblpQVERWTU0yOVpOVFpLSzJkNk1FNU5jemRWTlRGcVdYUlVablZXY2xRemQwYzNjVlpXWms4clIwTTBWMUZoVFdweVYyUjNjblp6ZFM5bGQzQlVhbU5rU20xQk1taEVaakZJY0d3NVFVaHlXRzh5YXpWMVdqSk9OemRuTkhoYU1qVTVSMnBTZW5JdmVXcE9lRFpPVDA5M2JrMU5ZMjB6WVc5RFZpdEpNbWMwV1VsRWFqUkRZVzkzUkRrNFQwTTJNak00YUdseGRVRTVVMDVCVDJsbmVtdGxZbnBKU1hvclZHVTVTVUZTWXpSQmQyTTVPVWxKZG5nclEzWXpNa1JMZWtGVFdHcHVTbFIyYUVScVRITXlObnB4TlhaTFpHOXlla3hLU1VaRGRXRlVkamQxT1ZwSlpXTk1XREJYZWpod1ZXZFVNRE5TYlZoVVN6SjFlR294ZVZCYWFuSlhhbGtyVXpaV1dsSjJNblpKZVVVck0zbGtVM2hxVUhGc0wydzFTMEZ3VUc5bFoyMVRiVWgwTUhaaVIzWXZWMDRyYVVReGN6aDNObkZUTm5KQ05GaExTMjB3WkV4VlIyMXpTa05KUjA5VmMybFhRamhvVmpGa1pVa3JabWxETlc4NWFXNWxVV0ZqWWtac1RYVmFjUzkwZWpsV00wSkVibmsxVWs1bVVEazFjbTFFZDJoc1JYUmpNVVpUZUdOb2NVMDJRVXcxT0ZoSFZGVkRUVFZuT1daWllWWkpOMEU0VEhFNVdITklURTlXZDNneVExRlBiV2s1Y1VSTVl5OXBkWHAxUzFOSmVFOHhiVGxGYVhoeVJsRkZaREJSUFE9PSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ=="),
			},
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ=="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                        "102475661842232156723",
							"primaryEmail":              "user2@sgnldemos.com",
							"$.name.fullName":           "user2 user2",
							"isAdmin":                   true,
							"isDelegatedAdmin":          false,
							"creationTime":              time.Date(2024, 2, 2, 23, 55, 44, 0, time.UTC),
							"changePasswordAtNextLogin": false,
							"emails": []framework.Object{
								map[string]any{
									"address": "user2@sgnldemos.com",
									"primary": true,
								},
								map[string]any{
									"address": "user2@sgnldemos.com.test-google-a.com",
								},
							},
							"nonEditableAliases": []string{
								"user2@sgnldemos.com.test-google-a.com",
							},
							"customerId": "CUST123456",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJRMEZGVXpoUlNVSnJVSHBXUVVkaU5HMDRRMUl4TWtsd1NqRkplbWRXY0Zodk9EZzRaRlJqVjFSU1dUQldLMkpvZUROWFdVUXlRVXBMZUZrek4xTnVRV0k0WkRGSFVXc3pNbXBFU0dweFIxSTNVbTVFY1hkNFYyUkViaTlYYzBOSlRVWXlWWFZZUjJ4WmNFZ3dkVVZOUms1WldDdGxWbGh6WVRSWGVYQTNNRkoyTlV4cVQyNXZNMWhDZVVNelowd3ZkRFJ3VVhaSGEzcG5kMjFRY25WWlNtOXVkRkV6TWs5ek1sY3lhRWhaT1VKNU9HZDBVelptVTNCWmRIQnBlRTF1VVV0T1VXSjZabFlyVFVJMFdqVm5ORkJZVlZCNFpqUkRaVEpWYzBwWFEwNUdTRzFGWm5ZelFrTXJlVTlCUldOWVpXUmtXQ3QyVTNSRU1GUjBaakkwU0VsTVkxWjNWSEIzU0hoM1dVUnpiazg0ZDBONWVURnNOREF3VUZOVlZHSnJOVzlCVWtGd2FqSkVMM2RZY0ZvNGJtUklhM0ZSZG1SR0szWjNiMEV3V1N0NlpUQjZZM1prY1VsTWQzcHdWakl6TDI1R0wwdElOMkpQY1RWcU1XRlNWVllyUkROME5FNHpaek54YVU0NGNsSjVjMWR4TVdoYWRreHlUMVIyVEdKamRXTlZkVkV3TVZnd2NIcDVjVFpsU0c1dlRVVldlV3R0TUhKVVpFRkJOR2RWT1ZGdFZqVXZaU3RCTVdsaFZsVnJPRVF5VFZscE1tSkpVVXRhUm5Sb1VrNWpLM2xoTTFkU2N5dDBaREZTWldGSGQwOU1WbE5uUVZOcFFXWnJRbFpZYURKcWJGVnFibFJ2TDNZM09VWllXbHAxVFQwPSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q0FFUzhRSUJrUHpWQUdiNG04Q1IxMklwSjFJemdWcFhvODg4ZFRjV1RSWTBWK2JoeDNXWUQyQUpLeFkzN1NuQWI4ZDFHUWszMmpESGpxR1I3Um5EcXd4V2REbi9Xc0NJTUYyVXVYR2xZcEgwdUVNRk5ZWCtlVlhzYTRXeXA3MFJ2NUxqT25vM1hCeUMzZ0wvdDRwUXZHa3pnd21QcnVZSm9udFEzMk9zMlcyaEhZOUJ5OGd0UzZmU3BZdHBpeE1uUUtOUWJ6ZlYrTUI0WjVnNFBYVVB4ZjRDZTJVc0pXQ05GSG1FZnYzQkMreU9BRWNYZWRkWCt2U3REMFR0ZjI0SElMY1Z3VHB3SHh3WURzbk84d0N5eTFsNDAwUFNVVGJrNW9BUkFwajJEL3dYcFo4bmRIa3FRdmRGK3Z3b0EwWSt6ZTB6Y3ZkcUlMd3pwVjIzL25GL0tIN2JPcTVqMWFSVVYrRDN0NE4zZzNxaU44clJ5c1dxMWhadkxyT1R2TGJjdWNVdVEwMVgwcHp5cTZlSG5vTUVWeWttMHJUZEFBNGdVOVFtVjUvZStBMWlhVlVrOEQyTVlpMmJJUUtaRnRoUk5jK3lhM1dScyt0ZDFSZWFHd09MVlNnQVNpQWZrQlZYaDJqbFVqblRvL3Y3OUZYWlp1TT0="),
			},
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultUserEntityConfig(),
				PageSize: 1,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q0FFUzhRSUJrUHpWQUdiNG04Q1IxMklwSjFJemdWcFhvODg4ZFRjV1RSWTBWK2JoeDNXWUQyQUpLeFkzN1NuQWI4ZDFHUWszMmpESGpxR1I3Um5EcXd4V2REbi9Xc0NJTUYyVXVYR2xZcEgwdUVNRk5ZWCtlVlhzYTRXeXA3MFJ2NUxqT25vM1hCeUMzZ0wvdDRwUXZHa3pnd21QcnVZSm9udFEzMk9zMlcyaEhZOUJ5OGd0UzZmU3BZdHBpeE1uUUtOUWJ6ZlYrTUI0WjVnNFBYVVB4ZjRDZTJVc0pXQ05GSG1FZnYzQkMreU9BRWNYZWRkWCt2U3REMFR0ZjI0SElMY1Z3VHB3SHh3WURzbk84d0N5eTFsNDAwUFNVVGJrNW9BUkFwajJEL3dYcFo4bmRIa3FRdmRGK3Z3b0EwWSt6ZTB6Y3ZkcUlMd3pwVjIzL25GL0tIN2JPcTVqMWFSVVYrRDN0NE4zZzNxaU44clJ5c1dxMWhadkxyT1R2TGJjdWNVdVEwMVgwcHp5cTZlSG5vTUVWeWttMHJUZEFBNGdVOVFtVjUvZStBMWlhVlVrOEQyTVlpMmJJUUtaRnRoUk5jK3lhM1dScyt0ZDFSZWFHd09MVlNnQVNpQWZrQlZYaDJqbFVqblRvL3Y3OUZYWlp1TT0="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":                        "114211199695002816249",
							"primaryEmail":              "sor-dev@sgnldemos.com",
							"$.name.fullName":           "SoR Development",
							"isAdmin":                   false,
							"isDelegatedAdmin":          true,
							"creationTime":              time.Date(2024, 4, 18, 0, 30, 35, 0, time.UTC),
							"changePasswordAtNextLogin": false,
							"emails": []framework.Object{
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
							"nonEditableAliases": []string{
								"sor-dev@sgnldemos.com.test-google-a.com",
							},
							"customerId": "CUST123456",
						},
					},
					NextCursor: "",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetGroupPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := googleworkspace.NewAdapter(&googleworkspace.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[googleworkspace.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultGroupEntityConfig(),
				PageSize: 1,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"kind":               "admin#directory#group",
							"id":                 "01qoc8b13vgdlqb",
							"email":              "emptygroup@sgnldemos.com",
							"name":               "Empty Group",
							"directMembersCount": "0",
							"adminCreated":       true,
							"nonEditableAliases": []string{
								"emptygroup@sgnldemos.com.test-google-a.com",
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJRMmx2ZDB4RFNteGlXRUl3WlZka2VXSXpWbmRSU0U1dVltMTRhMXBYTVhaamVUVnFZakl3YVV4RVJYZE9WR015VFZSbmVrMUVXWGhOZW14SlFUSkRPRzExVTJsQ1FUMDkifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"),
			},
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultGroupEntityConfig(),
				PageSize: 1,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"kind":               "admin#directory#group",
							"id":                 "048pi1tg0qf1f8g",
							"email":              "group2@sgnldemos.com",
							"name":               "Group2",
							"directMembersCount": "3",
							"adminCreated":       true,
							"nonEditableAliases": []string{
								"group2@sgnldemos.com.test-google-a.com",
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJRMmxWZDB4RFNtNWpiVGt4WTBSS1FXTXlaSFZpUjFKc1lsYzVla3h0VG5aaVUwbHpUbnBWTTA1RVdURk5WRkV4VFZSbk1GTkJUbWRvTldKbk5tWTNYMTlmWDE5QlVUMDkifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
			},
		},
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultGroupEntityConfig(),
				PageSize: 1,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"kind":               "admin#directory#group",
							"id":                 "030j0zll41obyxg",
							"email":              "hello@sgnldemos.com",
							"name":               "SGNLDemos",
							"directMembersCount": "2",
							"adminCreated":       true,
							"nonEditableAliases": []string{
								"hello@sgnldemos.com.test-google-a.com",
							},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJRMmxOZDB4RFNtOWFWM2h6WWpCQ2Vsb3lOWE5hUjFaMFlqTk5kVmt5T1hSSmFYZDRUVVJWZWs1NlNUVk9WRkY0VFd0blJGbFFZVlo1VEdZNVgxOWZYMTkzUlQwPSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="),
			},
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultGroupEntityConfig(),
				PageSize: 1,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func TestAdapterGetMemberPage(t *testing.T) {
	server := httptest.NewTLSServer(TestServerHandler)
	adapter := googleworkspace.NewAdapter(&googleworkspace.Datasource{
		Client: server.Client(),
	})

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[googleworkspace.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"first_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultMemberEntityConfig(),
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "eyJjb2xsZWN0aW9uQ3Vyc29yIjoiUTJsdmQweERTbXhpV0VJd1pWZGtlV0l6Vm5kUlNFNXVZbTE0YTFwWE1YWmplVFZxWWpJd2FVeEVSWGRPVkdNeVRWUm5lazFFV1hoTmVteEpRVEpET0cxMVUybENRVDA5In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor:           nil,
				CollectionCursor: testutil.GenPtr("Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"),
			},
		},
		"second_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultMemberEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor:           nil,
				CollectionCursor: testutil.GenPtr("Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"kind":     "admin#directory#member",
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
							"id":       "102475661842232156723",
							"groupId":  "048pi1tg0qf1f8g",
							"uniqueId": "048pi1tg0qf1f8g-102475661842232156723",
							"email":    "user2@sgnldemos.com",
							"role":     "OWNER",
							"type":     "USER",
							"status":   "ACTIVE",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJDalJKYURoTFNGRnFTWFpOYVdaNlowVlRSVzB4YUdOdFRrRmpNbVIxWWtkU2JHSlhPWHBNYlU1MllsSm5RbGxLZVVwcFRUaEZJaDhLSFFqSXZNaWZ6Z0VTRW0xaGNtTkFjMmR1YkdSbGJXOXpMbU52YlJnQllKeUppTThFIiwiY29sbGVjdGlvbklkIjoiMDQ4cGkxdGcwcWYxZjhnIiwiY29sbGVjdGlvbkN1cnNvciI6IlEybFZkMHhEU201amJUa3hZMFJLUVdNeVpIVmlSMUpzWWxjNWVreHRUblppVTBselRucFZNMDVFV1RGTlZGRXhUVlJuTUZOQlRtZG9OV0puTm1ZM1gxOWZYMTlCVVQwOSJ9",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("CjRJaDhLSFFqSXZNaWZ6Z0VTRW0xaGNtTkFjMmR1YkdSbGJXOXpMbU52YlJnQllKeUppTThFIh8KHQjIvMifzgESEm1hcmNAc2dubGRlbW9zLmNvbRgBYJyJiM8E"),
				CollectionID:     testutil.GenPtr("048pi1tg0qf1f8g"),
				CollectionCursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
			},
		},
		"third_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultMemberEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor:           testutil.GenPtr("CjRJaDhLSFFqSXZNaWZ6Z0VTRW0xaGNtTkFjMmR1YkdSbGJXOXpMbU52YlJnQllKeUppTThFIh8KHQjIvMifzgESEm1hcmNAc2dubGRlbW9zLmNvbRgBYJyJiM8E"),
				CollectionID:     testutil.GenPtr("048pi1tg0qf1f8g"),
				CollectionCursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"kind":     "admin#directory#member",
							"id":       "114211199695002816249",
							"groupId":  "048pi1tg0qf1f8g",
							"uniqueId": "048pi1tg0qf1f8g-114211199695002816249",
							"email":    "sor-dev@sgnldemos.com",
							"role":     "OWNER",
							"type":     "USER",
							"status":   "ACTIVE",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uQ3Vyc29yIjoiUTJsVmQweERTbTVqYlRreFkwUktRV015WkhWaVIxSnNZbGM1ZWt4dFRuWmlVMGx6VG5wVk0wNUVXVEZOVkZFeFRWUm5NRk5CVG1kb05XSm5ObVkzWDE5ZlgxOUJVVDA5In0=",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor:           nil,
				CollectionID:     nil,
				CollectionCursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
			},
		},
		"fourth_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultMemberEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor:           nil,
				CollectionID:     nil,
				CollectionCursor: testutil.GenPtr("Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"kind":     "admin#directory#member",
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
							"id":       "102475661842232156723",
							"groupId":  "030j0zll41obyxg",
							"uniqueId": "030j0zll41obyxg-102475661842232156723",
							"email":    "user2@sgnldemos.com",
							"role":     "OWNER",
							"type":     "USER",
							"status":   "ACTIVE",
						},
					},
					NextCursor: "eyJjb2xsZWN0aW9uQ3Vyc29yIjoiUTJsTmQweERTbTlhVjNoellqQkNlbG95TlhOYVIxWjBZak5OZFZreU9YUkphWGQ0VFVSVmVrNTZTVFZPVkZGNFRXdG5SRmxRWVZaNVRHWTVYMTlmWDE5M1JUMD0ifQ==",
				},
			},
			wantCursor: &pagination.CompositeCursor[string]{
				Cursor:           nil,
				CollectionID:     nil,
				CollectionCursor: testutil.GenPtr("Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="),
			},
		},
		"last_page": {
			ctx: context.Background(),
			request: &framework.Request[googleworkspace.Config]{
				Address: server.URL,
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer Testtoken",
				},
				Config: &googleworkspace.Config{
					APIVersion: "v1",
					Domain:     testutil.GenPtr("sgnldemos.com"),
				},
				Entity:   *PopulateDefaultMemberEntityConfig(),
				PageSize: 2,
			},
			inputRequestCursor: &pagination.CompositeCursor[string]{
				Cursor:           nil,
				CollectionID:     nil,
				CollectionCursor: testutil.GenPtr("Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="),
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects:    nil,
					NextCursor: "",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			if diff := cmp.Diff(tt.wantResponse, gotResponse); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}
