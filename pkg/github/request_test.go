// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package github_test

import (
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/github"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestPopulateRequestInfo(t *testing.T) {
	tests := map[string]struct {
		request         *github.Request
		wantRequestInfo *github.RequestInfo
		wantError       *framework.Error
	}{
		"nil_request": {
			request: nil,
			wantError: &framework.Error{
				Message: "Request is nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"invalid_entity": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "INVALID",
				PageSize:          100,
				Token:             "Bearer Testtoken",
			},
			wantError: &framework.Error{
				Message: "Invalid entity external ID: INVALID",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"enterprise_server_graphql_organization_entity": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Organization",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				EntityConfig: &framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "login",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "createdAt",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							ExternalId: "pinnedItemsRemaining",
							Type:       framework.AttributeTypeInt64,
							List:       false,
						},
						{
							ExternalId: "anyPinnableItems",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
					},
					ExternalId: "Organization",
				},
			},
			wantRequestInfo: &github.RequestInfo{
				Endpoint:   "https://ghe-test-server/api/graphql",
				HTTPMethod: "POST",
				Query: `query {
					enterprise (slug: "testID") {
						id
						organizations (first: 100) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								anyPinnableItems
								createdAt
								login
								pinnedItemsRemaining
							}
						}
					}
				}`,
			},
		},
		"enterprise_server_graphql_team_entity": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Team",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("testOrgAfter"), testutil.GenPtr("testTeamAfter")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultTeamEntityConfig(),
			},
			wantRequestInfo: &github.RequestInfo{
				Endpoint:   "https://ghe-test-server/api/graphql",
				HTTPMethod: "POST",
				Query: `query {
					enterprise (slug: "testID") {
						id
						organizations (first: 1, after: "testOrgAfter") {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								teams (first: 100, after: "testTeamAfter") {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										createdAt
										databaseId
										id
										members {
											edges {
												node {
													createdAt
													databaseId
													email
													id
													isViewer
													login
													updatedAt
												}
												role
											}
										}
										repositories {
											edges {
												node {
													allowUpdateBranch
													createdAt
													databaseId
													id
													name
													pushedAt
													url
												}
												permission
											}
										}
										slug
										updatedAt
										viewerCanAdminister
									}
								}
							}
						}
					}
				}`,
			},
		},
		"enterprise_cloud_graphql_issue_entity": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Issue",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("repoAfter1"), testutil.GenPtr("issueAfter1")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultIssueEntityConfig(),
			},
			wantRequestInfo: &github.RequestInfo{
				Endpoint:   "https://api.github.com/graphql",
				HTTPMethod: "POST",
				Query: `query {
					enterprise (slug: "testID") {
						id
						organizations (first: 1) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								repositories (first: 1, after: "repoAfter1") {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										issues (first: 100, after: "issueAfter1") {
											pageInfo {
												endCursor
												hasNextPage
											}
											nodes {
												author {
													login
												}
												createdAt
												id
												isPinned
												title
											}
										}
									}
								}
							}
						}
					}
				}`,
			},
		},
		"enterprise_server_rest_secretscanning_entity": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "SecretScanningAlert",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor:            nil,
			},
			wantRequestInfo: &github.RequestInfo{
				Endpoint:   "https://ghe-test-server/api/v3/enterprises/testID/secret-scanning/alerts?per_page=100",
				HTTPMethod: "GET",
				Query:      "",
			},
		},
		"enterprise_server_rest_missing_apiVersion": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        nil,
				EntityExternalID:  "SecretScanningAlert",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor:            nil,
			},
			wantError: &framework.Error{
				Message: "APIVersion is not set for an entity that is retrieved through the GitHub REST API.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			},
		},
		"enterprise_server_rest_secretscanning_entity_with_cursor": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "SecretScanningAlert",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://ghe-test-server/api/v3/enterprises/testID/secret-scanning/alerts?per_page=100&page=2"),
				},
			},
			wantRequestInfo: &github.RequestInfo{
				Endpoint:   "https://ghe-test-server/api/v3/enterprises/testID/secret-scanning/alerts?per_page=100&page=2",
				HTTPMethod: "GET",
				Query:      "",
			},
		},
		"enterprise_cloud_graphql_rest_secretscanning_entity": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "SecretScanningAlert",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				EntityConfig:      PopulateDefaultIssueEntityConfig(),
			},
			wantRequestInfo: &github.RequestInfo{
				Endpoint:   "https://api.github.com/enterprises/testID/secret-scanning/alerts?per_page=100",
				HTTPMethod: "GET",
				Query:      "",
			},
		},
		"enterprise_cloud_rest_secretscanning_entity_with_cursor": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "SecretScanningAlert",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://api.github.com/enterprises/testID/secret-scanning/alerts?per_page=100&page=2"),
				},
			},
			wantRequestInfo: &github.RequestInfo{
				Endpoint:   "https://api.github.com/enterprises/testID/secret-scanning/alerts?per_page=100&page=2",
				HTTPMethod: "GET",
				Query:      "",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRequestInfo, gotError := github.PopulateRequestInfo(tt.request)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}
			// These tests will mainly compare the endpointURL and HTTPMethod.
			// The query comparison will only be tested to see if both are empty or not.
			// query_test.go will be responsible for more robust query builder tests.
			if gotRequestInfo != nil && tt.wantRequestInfo != nil {
				if (gotRequestInfo.Query != "" && tt.wantRequestInfo.Query == "") ||
					(gotRequestInfo.Query == "" && tt.wantRequestInfo.Query != "") {
					t.Errorf("gotRequestInfo.Query: %v, wantRequestInfo.Query: %v", gotRequestInfo.Query, tt.wantRequestInfo.Query)
				}

				if gotRequestInfo.Endpoint != tt.wantRequestInfo.Endpoint {
					t.Errorf("gotRequestInfo.Endpoint: %v, wantRequestInfo.Endpoint: %v", gotRequestInfo.Endpoint, tt.wantRequestInfo.Endpoint)
				}

				if gotRequestInfo.HTTPMethod != tt.wantRequestInfo.HTTPMethod {
					t.Errorf("gotRequestInfo.HTTPMethod: %v, wantRequestInfo.HTTPMethod: %v", gotRequestInfo.HTTPMethod, tt.wantRequestInfo.HTTPMethod)
				}
			}
		})
	}
}
