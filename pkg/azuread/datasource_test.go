// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package azuread_test

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
	"github.com/sgnl-ai/adapters/pkg/azuread"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		body         []byte
		wantObjects  []map[string]interface{}
		wantNextLink *string
		wantErr      *framework.Error
	}{
		"single_page": {
			body: []byte(`{"value": [{"id": "00ub0oNGTSWTBKOLGLNR","status": "ACTIVE"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}], "@odata.nextLink": "https://graph.microsoft.com/v1.0/applications?$top=1&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9hY2M4NDhlOS1lOGVjLTRmZWItYTUyMS04ZDU4YjU0ODJlMDkwQXBwbGljYXRpb25fYWNjODQ4ZTktZThlYy00ZmViLWE1MjEtOGQ1OGI1NDgyZTA5AAAAAAAAAAAAAAA"}`),
			wantObjects: []map[string]interface{}{
				{"id": "00ub0oNGTSWTBKOLGLNR", "status": "ACTIVE"},
				{"id": "00ub0oNGTSWTBKOCHDKE", "status": "ACTIVE"},
			},
			wantNextLink: testutil.GenPtr("https://graph.microsoft.com/v1.0/applications?$top=1&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9hY2M4NDhlOS1lOGVjLTRmZWItYTUyMS04ZDU4YjU0ODJlMDkwQXBwbGljYXRpb25fYWNjODQ4ZTktZThlYy00ZmViLWE1MjEtOGQ1OGI1NDgyZTA5AAAAAAAAAAAAAAA"),
		},
		"invalid_object_structure": {
			body: []byte(`[{"id": "00ub0oNGTSWTBKOLGLNR","status": "ACTIVE"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal array into Go value of type azuread.DatasourceResponse.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_objects": {
			body: []byte(`{"value": [{"00ub0oNGTSWTBKOLGLNR"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]}`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: `Failed to unmarshal the datasource response: invalid character '}' after object key.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextLink, gotErr := azuread.ParseResponse(tt.body)

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

	azureadClient := azuread.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *azuread.Request
		wantRes *azuread.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "User",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"businessPhones":    []any{},
						"displayName":       "Conf Room Adams",
						"givenName":         nil,
						"jobTitle":          nil,
						"mail":              "Adams@M365x214355.onmicrosoft.com",
						"mobilePhone":       nil,
						"officeLocation":    nil,
						"preferredLanguage": nil,
						"surname":           nil,
						"userPrincipalName": "Adams@M365x214355.onmicrosoft.com",
						"id":                "6e7b768e-07e2-4810-8459-485f84f8f204",
					},
					{
						"businessPhones": []any{
							"+1 425 555 0109",
						},
						"displayName":       "Adele Vance",
						"givenName":         "Adele",
						"jobTitle":          "Product Marketing Manager",
						"mail":              "AdeleV@M365x214355.onmicrosoft.com",
						"mobilePhone":       nil,
						"officeLocation":    "18/2111",
						"preferredLanguage": "en-US",
						"surname":           "Vance",
						"userPrincipalName": "AdeleV@M365x214355.onmicrosoft.com",
						"id":                "87d349ed-44d7-43e1-9a83-5f2406dee5bd",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "User",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"businessPhones": []any{
							"8006427676",
						},
						"displayName":       "MOD Administrator",
						"givenName":         "MOD",
						"jobTitle":          nil,
						"mail":              "admin@M365x214355.onmicrosoft.com",
						"mobilePhone":       "5555555555",
						"officeLocation":    nil,
						"preferredLanguage": "en-US",
						"surname":           "Administrator",
						"userPrincipalName": "admin@M365x214355.onmicrosoft.com",
						"id":                "5bde3e51-d13b-4db1-9948-fe4b109d11a7",
					},
					{
						"businessPhones": []any{
							"+1 858 555 0110",
						},
						"displayName":       "Alex Wilber",
						"givenName":         "Alex",
						"jobTitle":          "Marketing Assistant",
						"mail":              "AlexW@M365x214355.onmicrosoft.com",
						"mobilePhone":       nil,
						"officeLocation":    "131/1104",
						"preferredLanguage": "en-US",
						"surname":           "Wilber",
						"userPrincipalName": "AlexW@M365x214355.onmicrosoft.com",
						"id":                "4782e723-f4f4-4af3-a76e-25e3bab0d896",
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := azureadClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetRolesPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	azureadClient := azuread.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *azuread.Request
		wantRes *azuread.Response
		wantErr *framework.Error
	}{
		"single_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Role",
				PageSize:              100,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "description",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "roleTemplateId",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":             "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
						"description":    "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
						"displayName":    "Directory Readers",
						"roleTemplateId": "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
					},
					{
						"id":             "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
						"description":    "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
						"displayName":    "Global Administrator",
						"roleTemplateId": "62e90394-3621-4004-a7cb-012177145e10",
					},
					{
						"id":             "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
						"description":    "Can create application registrations independent of the 'Users can register applications' setting.",
						"displayName":    "Application Developer",
						"roleTemplateId": "cf1c38e5-69f5-4237-9190-879624dced7c",
					},
					{
						"id":             "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
						"description":    "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
						"displayName":    "Privileged Role Administrator",
						"roleTemplateId": "e8611ab8-8f55-4a1e-953a-60213ab1f814",
					},
					{
						"id":             "62ceaa28-7382-48f9-j386-f8ed6e9a7c84",
						"description":    "Can read and write basic directory information. For granting access to applications, not intended for users.",
						"displayName":    "Directory Writers",
						"roleTemplateId": "9360feb5-f418-4baa-8175-e2a00bac4301",
					},
				},
			},
			wantErr: nil,
		},

		"first_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Role",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "description",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "roleTemplateId",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":             "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
						"description":    "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
						"displayName":    "Directory Readers",
						"roleTemplateId": "88d8e3e3-c189-46e8-94e1-9b9898b8876b",
					},
					{
						"id":             "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
						"description":    "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
						"displayName":    "Global Administrator",
						"roleTemplateId": "62e90394-3621-4004-a7cb-012177145e10",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("2"),
				},
			},
			wantErr: nil,
		},

		"middle_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:            "Bearer Testtoken",
				BaseURL:          server.URL,
				EntityExternalID: "Role",
				PageSize:         2,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("2"),
				},
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "description",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "roleTemplateId",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":             "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
						"description":    "Can create application registrations independent of the 'Users can register applications' setting.",
						"displayName":    "Application Developer",
						"roleTemplateId": "cf1c38e5-69f5-4237-9190-879624dced7c",
					},
					{
						"id":             "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
						"description":    "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
						"displayName":    "Privileged Role Administrator",
						"roleTemplateId": "e8611ab8-8f55-4a1e-953a-60213ab1f814",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("4"),
				},
			},
			wantErr: nil,
		},

		"last_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:            "Bearer Testtoken",
				BaseURL:          server.URL,
				EntityExternalID: "Role",
				PageSize:         2,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("4"),
				},
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "displayName",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "description",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "roleTemplateId",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":             "62ceaa28-7382-48f9-j386-f8ed6e9a7c84",
						"description":    "Can read and write basic directory information. For granting access to applications, not intended for users.",
						"displayName":    "Directory Writers",
						"roleTemplateId": "9360feb5-f418-4baa-8175-e2a00bac4301",
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := azureadClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetApplicationsPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	azureadClient := azuread.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *azuread.Request
		wantRes *azuread.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Application",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                        "acc848e9-e8ec-4feb-a521-8d58b5482e09",
						"deletedDateTime":           nil,
						"appId":                     "05b10a2d-62db-420c-8626-55f3a5e7865b",
						"applicationTemplateId":     nil,
						"disabledByMicrosoftStatus": nil,
						"createdDateTime":           "2020-02-13T21:09:15Z",
						"displayName":               "apisandboxproxy",
						"description":               nil,
						"groupMembershipClaims":     nil,
						"identifierUris": []any{
							"https://M365x214355.onmicrosoft.com/apisandboxproxy",
						},
						"isDeviceOnlyAuthSupported":         nil,
						"isFallbackPublicClient":            true,
						"notes":                             nil,
						"publisherDomain":                   "M365x214355.onmicrosoft.com",
						"serviceManagementReference":        nil,
						"signInAudience":                    "AzureADMyOrg",
						"tags":                              []any{},
						"tokenEncryptionKeyId":              nil,
						"samlMetadataUrl":                   nil,
						"defaultRedirectUri":                nil,
						"certification":                     nil,
						"optionalClaims":                    nil,
						"servicePrincipalLockConfiguration": nil,
						"requestSignatureVerification":      nil,
						"addIns":                            []any{},
						"api": map[string]any{
							"acceptMappedClaims":          nil,
							"knownClientApplications":     []any{},
							"requestedAccessTokenVersion": nil,
							"oauth2PermissionScopes":      []any{},
							"preAuthorizedApplications":   []any{},
						},
						"appRoles": []any{},
						"info": map[string]any{
							"logoUrl":             nil,
							"marketingUrl":        nil,
							"privacyStatementUrl": nil,
							"supportUrl":          nil,
							"termsOfServiceUrl":   nil,
						},
						"keyCredentials": []any{},
						"parentalControlSettings": map[string]any{
							"countriesBlockedForMinors": []any{},
							"legalAgeGroupRule":         "Allow",
						},
					},
					{
						"id":                        "cfa98ac0-a32c-4b4c-a78b-94c9912ed7b2",
						"deletedDateTime":           nil,
						"appId":                     "c305b21c-fda6-4ecb-aa01-8a8141fdfd51",
						"applicationTemplateId":     nil,
						"disabledByMicrosoftStatus": nil,
						"createdDateTime":           "2018-03-27T02:45:04Z",
						"displayName":               "EduPopulationHelper",
						"description":               nil,
						"groupMembershipClaims":     nil,
						"identifierUris": []any{
							"https://M365x214355.onmicrosoft.com/a60c216f-657f-4925-980a-d8ef69942167",
						},
						"isDeviceOnlyAuthSupported":         nil,
						"isFallbackPublicClient":            false,
						"notes":                             nil,
						"publisherDomain":                   nil,
						"serviceManagementReference":        nil,
						"signInAudience":                    "AzureADMyOrg",
						"tags":                              []any{},
						"tokenEncryptionKeyId":              nil,
						"samlMetadataUrl":                   nil,
						"defaultRedirectUri":                nil,
						"certification":                     nil,
						"optionalClaims":                    nil,
						"servicePrincipalLockConfiguration": nil,
						"requestSignatureVerification":      nil,
						"addIns":                            []any{},
						"appRoles":                          []any{},
						"info": map[string]any{
							"logoUrl":             nil,
							"marketingUrl":        nil,
							"privacyStatementUrl": nil,
							"supportUrl":          nil,
							"termsOfServiceUrl":   nil,
						},
						"keyCredentials": []any{},
						"parentalControlSettings": map[string]any{
							"countriesBlockedForMinors": []any{},
							"legalAgeGroupRule":         "Allow",
						},
						"publicClient": map[string]any{
							"redirectUris": []any{},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/applications?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA"),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Application",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr(server.URL + "/v1.0/applications?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                                "d6dbf9e0-98a4-4eea-b4c1-df8695277868",
						"deletedDateTime":                   nil,
						"appId":                             "377a3df7-6ff0-42b2-ad55-194dcc7aacd9",
						"applicationTemplateId":             nil,
						"disabledByMicrosoftStatus":         nil,
						"createdDateTime":                   "2020-07-24T07:54:39Z",
						"displayName":                       "permissions-scraper-app",
						"description":                       nil,
						"groupMembershipClaims":             nil,
						"identifierUris":                    []any{},
						"isDeviceOnlyAuthSupported":         nil,
						"isFallbackPublicClient":            nil,
						"notes":                             nil,
						"publisherDomain":                   "M365x214355.onmicrosoft.com",
						"serviceManagementReference":        nil,
						"signInAudience":                    "AzureADandPersonalMicrosoftAccount",
						"tags":                              []any{},
						"tokenEncryptionKeyId":              nil,
						"samlMetadataUrl":                   nil,
						"defaultRedirectUri":                nil,
						"certification":                     nil,
						"optionalClaims":                    nil,
						"servicePrincipalLockConfiguration": nil,
						"requestSignatureVerification":      nil,
						"addIns":                            []any{},
						"api": map[string]any{
							"acceptMappedClaims":          nil,
							"knownClientApplications":     []any{},
							"requestedAccessTokenVersion": float64(2),
							"oauth2PermissionScopes":      []any{},
							"preAuthorizedApplications":   []any{},
						},
						"appRoles": []any{},
						"info": map[string]any{
							"logoUrl":             nil,
							"marketingUrl":        nil,
							"privacyStatementUrl": nil,
							"supportUrl":          nil,
							"termsOfServiceUrl":   nil,
						},
						"keyCredentials": []any{},
						"parentalControlSettings": map[string]any{
							"countriesBlockedForMinors": []any{},
							"legalAgeGroupRule":         "Allow",
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := azureadClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetDevicesPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	azureadClient := azuread.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *azuread.Request
		wantRes *azuread.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Device",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                     "0357837b-ca6e-402d-9429-9e54dd51d97a",
						"accountEnabled":         true,
						"deviceId":               "00000000-0000-0000-0000-000000000000",
						"deviceVersion":          float64(1),
						"displayName":            "contoso_pixel",
						"Manufacturer":           "Google",
						"Model":                  "Pixel 3a",
						"operatingSystemVersion": "10.0",
					},
					{
						"id":                     "4d1ed9a4-519e-421b-b9f6-158991feff5b",
						"accountEnabled":         true,
						"deviceId":               "00000000-0000-0000-0000-000000000001",
						"deviceVersion":          float64(1),
						"displayName":            "contoso_galaxy",
						"Manufacturer":           "Samsung",
						"Model":                  "Galaxy Note 7",
						"operatingSystemVersion": "8.2",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/devices?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA"),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Device",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "id",
						Type:       framework.AttributeTypeString,
					},
				},
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr(server.URL + "/v1.0/devices?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                     "6a59ea83-02bd-468f-a40b-f2c3d1821983",
						"accountEnabled":         true,
						"deviceId":               "00000000-0000-0000-0000-000000000002",
						"deviceVersion":          float64(1),
						"displayName":            "contoso_iphone",
						"Manufacturer":           "Apple",
						"Model":                  "iPhone 11 Pro Max",
						"operatingSystemVersion": "11.2",
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := azureadClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetGroupMembersPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	azureadClient := azuread.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *azuread.Request
		wantRes *azuread.Response
		wantErr *framework.Error
	}{
		"group_page_1_of_2_member_page_1_of_2": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "6e7b768e-07e2-4810-8459-485f84f8f204-02bd9fd6-8f93-4758-87c3-1fb73740a315",
						"memberId": "6e7b768e-07e2-4810-8459-485f84f8f204",
						"groupId":  "02bd9fd6-8f93-4758-87c3-1fb73740a315",
					},
					{
						"id":       "87d349ed-44d7-43e1-9a83-5f2406dee5bd-02bd9fd6-8f93-4758-87c3-1fb73740a315",
						"memberId": "87d349ed-44d7-43e1-9a83-5f2406dee5bd",
						"groupId":  "02bd9fd6-8f93-4758-87c3-1fb73740a315",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// We're of syncing Members for a specific Groups, so this cursor is to the next page of Members.
					Cursor:       testutil.GenPtr("https://graph.microsoft.com/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
					CollectionID: testutil.GenPtr("02bd9fd6-8f93-4758-87c3-1fb73740a315"),
					// GroupCursor to the next page of Groups.
					CollectionCursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/groups?$select=id&$top=1&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA"),
				},
			},
			wantErr: nil,
		},
		"group_page_1_of_2_member_page_2_of_2": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					// We're of syncing Members for a specific Groups, so this cursor is to the next page of Members.
					Cursor:       testutil.GenPtr(server.URL + "/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
					CollectionID: testutil.GenPtr("02bd9fd6-8f93-4758-87c3-1fb73740a315"),
					// GroupCursor to the next page of Groups.
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/groups?$select=id&$top=1&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "5bde3e51-d13b-4db1-9948-fe4b109d11a7-02bd9fd6-8f93-4758-87c3-1fb73740a315",
						"memberId": "5bde3e51-d13b-4db1-9948-fe4b109d11a7",
						"groupId":  "02bd9fd6-8f93-4758-87c3-1fb73740a315",
					},
					{
						"id":       "4782e723-f4f4-4af3-a76e-25e3bab0d896-02bd9fd6-8f93-4758-87c3-1fb73740a315",
						"memberId": "4782e723-f4f4-4af3-a76e-25e3bab0d896",
						"groupId":  "02bd9fd6-8f93-4758-87c3-1fb73740a315",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// There is no Cursor since we've finished all pages of Members for the current Group.
					CollectionID: testutil.GenPtr("02bd9fd6-8f93-4758-87c3-1fb73740a315"),
					// GroupCursor to the next page of Groups.
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/groups?$select=id&$top=1&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA"),
				},
			},
			wantErr: nil,
		},
		"group_page_2_of_2_member_page_1_of_2": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					// There is no Cursor since we've finished all pages of Members for the current Group.
					CollectionID: testutil.GenPtr("02bd9fd6-8f93-4758-87c3-1fb73740a315"),
					// GroupCursor to the next page of Groups.
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/groups?$select=id&$top=1&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "6e7b768e-07e2-4810-8459-485f84f8f204-06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
						"memberId": "6e7b768e-07e2-4810-8459-485f84f8f204",
						"groupId":  "06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
					},
					{
						"id":       "87d349ed-44d7-43e1-9a83-5f2406dee5bd-06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
						"memberId": "87d349ed-44d7-43e1-9a83-5f2406dee5bd",
						"groupId":  "06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// We're of syncing Members for a specific Groups, so this cursor is to the next page of Members.
					Cursor:       testutil.GenPtr("https://graph.microsoft.com/v1.0/groups/06f62f70-9827-4e6e-93ef-8e0f2d9b7b23/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
					CollectionID: testutil.GenPtr("06f62f70-9827-4e6e-93ef-8e0f2d9b7b23"),
					// There is no CollectionCursor since we're currently processing the last Group.
				},
			},
			wantErr: nil,
		},
		"group_page_2_of_2_member_page_2_of_2": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr(server.URL + "/v1.0/groups/06f62f70-9827-4e6e-93ef-8e0f2d9b7b23/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
					CollectionID: testutil.GenPtr("06f62f70-9827-4e6e-93ef-8e0f2d9b7b23"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "5bde3e51-d13b-4db1-9948-fe4b109d11a7-06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
						"memberId": "5bde3e51-d13b-4db1-9948-fe4b109d11a7",
						"groupId":  "06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
					},
					{
						"id":       "4782e723-f4f4-4af3-a76e-25e3bab0d896-06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
						"memberId": "4782e723-f4f4-4af3-a76e-25e3bab0d896",
						"groupId":  "06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
					},
				},
				// Cursor and CollectionCursor are both nil, so no cursor is set as this is the last page for the sync.
			},
			wantErr: nil,
		},
		"group_members_too_many_groups_returned": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/groups?$select=id&$top=1&$skiptoken=TOO-MANY-GROUPS"),
				},
			},
			wantErr: &framework.Error{
				Message: "Too many collection objects returned in response; expected 1, got 2.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		// If no groups are present in the current page and there is no next group link, exit early.
		"group_members_no_groups_returned": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/groups?$select=id&$top=1&$skiptoken=NOT-ENOUGH-GROUPS"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
			},
		},
		"group_member_invalid_group_cursor": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr(server.URL + "/invalid"),
				},
			},
			wantErr: &framework.Error{
				Message: "Datasource rejected request, returned status code: 404.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"group_page_1_of_1_member_page_1_of_1_with_filters": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Filter:                testutil.GenPtr("id eq '6e7b768e-07e2-4810-8459-485f84f8f204'"),
				ParentFilter:          testutil.GenPtr("id in ('02bd9fd6-8f93-4758-87c3-1fb73740a315','0a53828f-36c9-44c3-be3d-99a7fce977ac')"),
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "6e7b768e-07e2-4810-8459-485f84f8f204-02bd9fd6-8f93-4758-87c3-1fb73740a315",
						"memberId": "6e7b768e-07e2-4810-8459-485f84f8f204",
						"groupId":  "02bd9fd6-8f93-4758-87c3-1fb73740a315",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// We're of syncing Members for a specific Groups, so this cursor is to the next page of Members.
					Cursor:       testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA"),
					CollectionID: testutil.GenPtr("02bd9fd6-8f93-4758-87c3-1fb73740a315"),
					// GroupCursor to the next page of Groups.
					CollectionCursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/groups?$select=id&$top=1&$filter=id+in+%28%2702bd9fd6-8f93-4758-87c3-1fb73740a315%27%2C%270a53828f-36c9-44c3-be3d-99a7fce977ac%27%29&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA"),
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := azureadClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetRoleMembersPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	azureadClient := azuread.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *azuread.Request
		wantRes *azuread.Response
		wantErr *framework.Error
	}{
		"user_page_1_of_4_role_page_1_of_2": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "0fea7f0d-dea1-458d-9099-69fcc2e3cd42-65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
						"roleId":   "0fea7f0d-dea1-458d-9099-69fcc2e3cd42",
						"memberId": "65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
					},
					{
						"id":       "795326a8-6eef-410e-9604-649ca68e1241-65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
						"roleId":   "795326a8-6eef-410e-9604-649ca68e1241",
						"memberId": "65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// We're of syncing Roles for a specific Users, so this cursor is to the next page of Roles.
					Cursor:       testutil.GenPtr("https://graph.microsoft.com/v1.0/users/65bb46a4-7d3j-9302-8a21-4d90f7a0efdb/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4"),
					CollectionID: testutil.GenPtr("65bb46a4-7d3j-9302-8a21-4d90f7a0efdb"),
					// UserCursor to the next page of Users.
					CollectionCursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_1"),
				},
			},
			wantErr: nil,
		},
		"user_page_1_of_4_role_page_2_of_2": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					// We're of syncing Roles for a specific User, so this cursor is to the next page of Roles.
					Cursor:       testutil.GenPtr(server.URL + "/v1.0/users/65bb46a4-7d3j-9302-8a21-4d90f7a0efdb/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4"),
					CollectionID: testutil.GenPtr("65bb46a4-7d3j-9302-8a21-4d90f7a0efdb"),
					// UserCursor to the next page of Users.
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_1"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "62ceaa28-4794-48f9-9b54-f8ed6e9a7c84-65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
						"roleId":   "62ceaa28-4794-48f9-9b54-f8ed6e9a7c84",
						"memberId": "65bb46a4-7d3j-9302-8a21-4d90f7a0efdb",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// There is no Cursor since we've finished all pages of Roles for the current User.
					CollectionID: testutil.GenPtr("65bb46a4-7d3j-9302-8a21-4d90f7a0efdb"),
					// UserCursor to the next page of Users.
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_1"),
				},
			},
			wantErr: nil,
		},

		"user_page_2_of_4_role_page_1_of_2": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					// There is no Cursor since we've finished all pages of Roles for the current User.
					CollectionID: testutil.GenPtr("65bb46a4-7d3j-9302-8a21-4d90f7a0efdb"),
					// UserCursor to the next page of Users.
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_1"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "62ceaa28-4794-48f9-9b54-f8ed6e9a7c84-df102bb2-2365-235g-a2g6-edb774169548",
						"roleId":   "62ceaa28-4794-48f9-9b54-f8ed6e9a7c84",
						"memberId": "df102bb2-2365-235g-a2g6-edb774169548",
					},
					{
						"id":       "795326a8-6eef-410e-9604-649ca68e1241-df102bb2-2365-235g-a2g6-edb774169548",
						"roleId":   "795326a8-6eef-410e-9604-649ca68e1241",
						"memberId": "df102bb2-2365-235g-a2g6-edb774169548",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// We're of syncing Role for a specific Users, so this cursor is to the next page of Roles.
					Cursor:       testutil.GenPtr("https://graph.microsoft.com/v1.0/users/df102bb2-2365-235g-a2g6-edb774169548/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4"),
					CollectionID: testutil.GenPtr("df102bb2-2365-235g-a2g6-edb774169548"),
					// UserCursor to the next page of Users.
					CollectionCursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=1&$skiptoken=RFNwdAIAADA6YWFyb24uYXlhbGE4OTUwNzM0N0BzZ25sYWFkZGV2MS5vbm1pY3Jvc29mdC5jb20pVXNlcl9kZjEwMmJiMi0zNTMyLTQ1M2MtYTNiNC1lZGI3NzQxNjk1NDgAMDphYXJvbi5heWFsYTg5NTA3MzQ3QHNnbmxhYWRkZXYxLm9ubWljcm9zb2Z0LmNvbSlVc2VyX2RmMTAyYmIyLTM1MzItNDUzYy1hM2I0LWVkYjc3NDE2OTU0OLkAAAAAAAAAAAAA"),
				},
			},
			wantErr: nil,
		},

		"user_page_2_of_4_role_page_2_of_2": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr(server.URL + "/v1.0/users/df102bb2-2365-235g-a2g6-edb774169548/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4"),
					CollectionID:     testutil.GenPtr("df102bb2-2365-235g-a2g6-edb774169548"),
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=1&$skiptoken=RFNwdAIAADA6YWFyb24uYXlhbGE4OTUwNzM0N0BzZ25sYWFkZGV2MS5vbm1pY3Jvc29mdC5jb20pVXNlcl9kZjEwMmJiMi0zNTMyLTQ1M2MtYTNiNC1lZGI3NzQxNjk1NDgAMDphYXJvbi5heWFsYTg5NTA3MzQ3QHNnbmxhYWRkZXYxLm9ubWljcm9zb2Z0LmNvbSlVc2VyX2RmMTAyYmIyLTM1MzItNDUzYy1hM2I0LWVkYjc3NDE2OTU0OLkAAAAAAAAAAAAA"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "e8d9279e-6883-4add-96e8-5f7c8df5637f-df102bb2-2365-235g-a2g6-edb774169548",
						"roleId":   "e8d9279e-6883-4add-96e8-5f7c8df5637f",
						"memberId": "df102bb2-2365-235g-a2g6-edb774169548",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// cursor is nil, as this is the last page of roles
					CollectionID:     testutil.GenPtr("df102bb2-2365-235g-a2g6-edb774169548"),
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=1&$skiptoken=RFNwdAIAADA6YWFyb24uYXlhbGE4OTUwNzM0N0BzZ25sYWFkZGV2MS5vbm1pY3Jvc29mdC5jb20pVXNlcl9kZjEwMmJiMi0zNTMyLTQ1M2MtYTNiNC1lZGI3NzQxNjk1NDgAMDphYXJvbi5heWFsYTg5NTA3MzQ3QHNnbmxhYWRkZXYxLm9ubWljcm9zb2Z0LmNvbSlVc2VyX2RmMTAyYmIyLTM1MzItNDUzYy1hM2I0LWVkYjc3NDE2OTU0OLkAAAAAAAAAAAAA"),
				},
			},
			wantErr: nil,
		},

		"user_page_3_of_4_with_no_roles": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID:     testutil.GenPtr("df102bb2-2365-235g-a2g6-edb774169548"),
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=1&$skiptoken=RFNwdAIAADA6YWFyb24uYXlhbGE4OTUwNzM0N0BzZ25sYWFkZGV2MS5vbm1pY3Jvc29mdC5jb20pVXNlcl9kZjEwMmJiMi0zNTMyLTQ1M2MtYTNiNC1lZGI3NzQxNjk1NDgAMDphYXJvbi5heWFsYTg5NTA3MzQ3QHNnbmxhYWRkZXYxLm9ubWljcm9zb2Z0LmNvbSlVc2VyX2RmMTAyYmIyLTM1MzItNDUzYy1hM2I0LWVkYjc3NDE2OTU0OLkAAAAAAAAAAAAA"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				// User has no role assignments
				Objects: []map[string]any{},
				// Cursor and CollectionCursor are both nil, so no cursor is set as this is the last page for the sync.
				NextCursor: &pagination.CompositeCursor[string]{
					CollectionID:     testutil.GenPtr("201d31c0-653d-43a6-adf0-aee89a79c805"),
					CollectionCursor: testutil.GenPtr("https://graph.microsoft.com/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_2"),
				},
			},
			wantErr: nil,
		},

		// API: /v1.0/users/{user_id}/transitiveMemberOf/microsoft.graph.directoryRole
		// Issue: The API performs pagination before applying role filtering, leading to potential scenarios as follows:
		"user_page_4__role_page_1__two_roles_and_next_cursor": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID:     testutil.GenPtr("201d31c0-653d-43a6-adf0-aee89a79c805"),
					CollectionCursor: testutil.GenPtr(server.URL + "/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_2"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "0fea7f0d-dea1-458d-9099-69fcc2e3cd42-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy",
						"roleId":   "0fea7f0d-dea1-458d-9099-69fcc2e3cd42",
						"memberId": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy",
					},
					{
						"id":       "d973db57-eb50-4356-959e-f1ce19a22b98-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy",
						"roleId":   "d973db57-eb50-4356-959e-f1ce19a22b98",
						"memberId": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr("https://graph.microsoft.com/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$select=id%2cdisplayName&$top=2&$skiptoken=RFNwdAoAAQAAAAAAAAAAFAAAABkp8fswrv1Ls8cLjYDqBRABAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA"),
					CollectionID: testutil.GenPtr("uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"),
					// CollectionCursor is nil as this is the last page of users
				},
			},
			wantErr: nil,
		},

		"user_page_4__role_page_2__zero_roles_and_next_cursor": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr(server.URL + "/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$select=id%2cdisplayName&$top=2&$skiptoken=RFNwdAoAAQAAAAAAAAAAFAAAABkp8fswrv1Ls8cLjYDqBRABAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA"),
					CollectionID: testutil.GenPtr("uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr("https://graph.microsoft.com/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$top=2&$select=id%2cdisplayName&$skiptoken=RFNwdAoAAAAAAAAAAAAAFAAAAPWE8iLxC5NNtqCdf_NZ8bcCAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA"),
					CollectionID: testutil.GenPtr("uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"),
					// CollectionCursor is nil as this is the last page of users
				},
			},
			wantErr: nil,
		},

		"user_page_4__role_page_2__one_role_and_next_cursor": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr(server.URL + "/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$top=2&$select=id%2cdisplayName&$skiptoken=RFNwdAoAAAAAAAAAAAAAFAAAAPWE8iLxC5NNtqCdf_NZ8bcCAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA"),
					CollectionID: testutil.GenPtr("uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "540b4b34-c25b-437d-8eee-329463952334-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy",
						"roleId":   "540b4b34-c25b-437d-8eee-329463952334",
						"memberId": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr("https://graph.microsoft.com/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$top=2&$select=id%2cdisplayName&$skiptoken=RFNwdAoAAAAAAAAAAAAAFAAAABgFnxJuzI1NsFSV18Bt7PgCAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA"),
					CollectionID: testutil.GenPtr("uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"),
					// CollectionCursor is nil as this is the last page of users
				},
			},
			wantErr: nil,
		},

		"user_page_4__role_page_4__one_role_and_No_cursor": {
			context: context.Background(),
			request: &azuread.Request{
				Token:                 "Bearer Testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "RoleMember",
				PageSize:              2,
				APIVersion:            "v1.0",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr(server.URL + "/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$top=2&$select=id%2cdisplayName&$skiptoken=RFNwdAoAAAAAAAAAAAAAFAAAABgFnxJuzI1NsFSV18Bt7PgCAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA"),
					CollectionID: testutil.GenPtr("uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"),
				},
			},
			wantRes: &azuread.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "fc6c3c82-669c-4e24-b089-2a2847a43d14-uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy",
						"roleId":   "fc6c3c82-669c-4e24-b089-2a2847a43d14",
						"memberId": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy",
					},
				},
				// CollectionCursor and Cursor is nil as this is the last page of sync
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := azureadClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
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
