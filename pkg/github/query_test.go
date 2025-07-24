// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package github_test

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/github"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestConstructQuery(t *testing.T) {
	tests := map[string]struct {
		request   *github.Request
		wantQuery string
		wantError *framework.Error
	}{
		"nil_request": {
			request: nil,
			wantError: &framework.Error{
				Message: "Request is nil.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"invalid_entity": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server/api/graphql",
				EntityExternalID:  "invalid",
				PageSize:          100,
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				Token:             "Bearer Testtoken",
			},
			wantError: &framework.Error{
				Message: "Unsupported Query for provided entity ID: invalid",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"valid_without_cursor_for_org_entity": {
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
			wantQuery: `query {
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
		"valid_with_cursor_for_org_entity": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Organization",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("testcursor")},
					nil,
					nil,
				),
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
					},
					ExternalId: "Organization",
				},
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 100, after: "testcursor") {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							createdAt
							login
							pinnedItemsRemaining
						}
					}
				}
			}`,
		},
		"missing_cursor_for_member_entity": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "OrganizationUser",
				PageSize:          100,
				Token:             "Bearer Testtoken",
			},
			wantError: &framework.Error{
				Message: "Cursor or CollectionCursor is nil for OrganizationUser query.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"missing_collectionID_for_member_entity": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "OrganizationUser",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: nil,
				},
			},
			wantError: &framework.Error{
				Message: "Cursor or CollectionCursor is nil for OrganizationUser query.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			},
		},
		"user_entity_with_cursor": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "User",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("testOrgAfter"), testutil.GenPtr("testUserAfter")},
					nil,
					nil,
				),
				EntityConfig: &framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "login",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "email",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
					ExternalId: "User",
				},
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 1, after: "testOrgAfter") {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							membersWithRole (first: 100, after: "testUserAfter") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									email
									id
									login
									name
								}
							}
						}
					}
				}
			}`,
		},
		"team_builder_with_no_children": {
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
				EntityConfig: &framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
					ExternalId: "Team",
				},
			},
			wantQuery: `query {
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
									id
									name
								}
							}
						}
					}
				}
			}`,
		},
		"team_builder_with_children": {
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
				EntityConfig: &framework.EntityConfig{
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
					ExternalId: "Team",
					ChildEntities: []*framework.EntityConfig{
						{
							ExternalId: "$.members.edges",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "$.node.id",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
								{
									ExternalId: "$.node.name",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
								{
									ExternalId: "role",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
							},
						},
						{
							ExternalId: "$.repositories.edges",
							Attributes: []*framework.AttributeConfig{
								{
									ExternalId: "$.node.id",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
								{
									ExternalId: "$.node.name",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
								{
									ExternalId: "permission",
									Type:       framework.AttributeTypeString,
									List:       false,
								},
							},
						},
					},
				},
			},
			wantQuery: `query {
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
									id
									members {
										edges {
											node {
												id
												name
											}
											role
										}
									}
									name
									repositories {
										edges {
											node {
												id
												name
											}
											permission
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_team_builder_attributes": {
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
			wantQuery: `query {
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
		"default_organization_builder_attributes": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Organization",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("testOrgAfter")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultOrganizationEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 100, after: "testOrgAfter") {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							createdAt
							databaseId
							email
							id
							login
							updatedAt
							viewerCanCreateTeams
							viewerIsAMember
						}
					}
				}
			}`,
		},
		"default_repository_builder_attributes": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Repository",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testRepoAfter")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultRepositoryEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 100, after: "testRepoAfter") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									allowUpdateBranch
									collaborators {
										edges {
											node {
												id
											}
											permission
										}
									}
									createdAt
									databaseId
									id
									name
									pushedAt
								}
							}
						}
					}
				}
			}`,
		},
		"default_user_builder_attributes": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "User",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("userAfter")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultUserEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							membersWithRole (first: 100, after: "userAfter") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									createdAt
									databaseId
									email
									id
									isViewer
									login
									updatedAt
								}
							}
						}
					}
				}
			}`,
		},
		"default_organizationuser_builder_attributes": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "OrganizationUser",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("userAfter")},
					testutil.GenPtr("testOrgID"),
					nil,
				),
				EntityConfig: PopulateDefaultOrganizationUserEntityConfig(),
			},
			wantQuery: `query {
				organization (login: "testOrgID") {
					id
					membersWithRole (first: 100, after: "userAfter") {
						pageInfo {
							endCursor
							hasNextPage
						}
						edges {
							node {
								id
								organizationVerifiedDomainEmails (login: "testOrgID")
							}
							role
						}
					}
				}
			}`,
		},
		"default_collaborator_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Collaborator",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultCollaboratorEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									collaborators (first: 100) {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											createdAt
											databaseId
											email
											id
											isViewer
											login
											updatedAt
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_label_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Label",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1"), testutil.GenPtr("labelCursor")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultLabelEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									labels (first: 100, after: "labelCursor") {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											color
											createdAt
											id
											isDefault
											name
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_issuelabel_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "IssueLabel",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("repoAfter1"), testutil.GenPtr("labelAfter1"), testutil.GenPtr("issueAfter1")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultIssueLabelEntityConfig(),
			},
			wantQuery: `query {
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
										labels (first: 1, after: "labelAfter1") {
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
														id
														title
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}`,
		},
		"default_pullrequest_label_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("SGNL"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "PullRequestLabel",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1"), nil, testutil.GenPtr("PRCursor")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultPullRequestLabelEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "SGNL") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									labels (first: 1) {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											pullRequests (first: 100, after: "PRCursor") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													id
													title
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"entity_with_multiple_cursors": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "Collaborator",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("repoAfter"), testutil.GenPtr("testcursor1")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultCollaboratorEntityConfig(),
			},
			wantQuery: `query {
					enterprise (slug: "testID") {
						id
						organizations (first: 1) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								repositories (first: 1, after: "repoAfter") {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										collaborators (first: 100, after: "testcursor1") {
											pageInfo {
												endCursor
												hasNextPage
											}
											nodes {
												createdAt
												databaseId
												email
												id
												isViewer
												login
												updatedAt
											}
										}
									}
								}
							}
						}
					}
				}`,
		},
		"default_issue_entity_attributes": {
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
			wantQuery: `query {
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
		"default_issueassignee_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "IssueAssignee",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("issueAfter1"), testutil.GenPtr("assigneeAfter1")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultIssueAssigneeEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1) {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									issues (first: 1, after: "issueAfter1") {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											assignees (first: 100, after: "assigneeAfter1") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													id
													login
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"issueassignee_assignee_after_parameter_fix": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "IssueAssignee",
				PageSize:          50,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("issueAfter1"), testutil.GenPtr("assigneeAfter456")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultIssueAssigneeEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1) {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									issues (first: 1, after: "issueAfter1") {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											assignees (first: 50, after: "assigneeAfter456") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													id
													login
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_issueparticipant_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("testID"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "IssueParticipant",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("orgAfter1"), nil, testutil.GenPtr("issueAfter1"), testutil.GenPtr("participantAfter1")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultIssueParticipantEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "testID") {
					id
					organizations (first: 1, after: "orgAfter1") {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1) {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									issues (first: 1, after: "issueAfter1") {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											participants (first: 100, after: "participantAfter1") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													id
													login
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_pullrequest_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://api.github.com",
				EnterpriseSlug:    testutil.GenPtr("SGNL"),
				IsEnterpriseCloud: true,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "PullRequest",
				PageSize:          100,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1"), testutil.GenPtr("PRCursor")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultPullRequestEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "SGNL") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									pullRequests (first: 100, after: "PRCursor") {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											author {
												login
											}
											baseRepository {
												id
											}
											closed
											createdAt
											headRepository {
												id
											}
											id
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
		"default_pullrequest_changedfile_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("SGNL"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "PullRequestChangedFile",
				PageSize:          5,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1"), nil, testutil.GenPtr("fileCursor")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultPullRequestChangedFileEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "SGNL") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									pullRequests (first: 1) {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											files (first: 5, after: "fileCursor") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													changeType
													path
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_pullrequest_assignee_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("SGNL"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "PullRequestAssignee",
				PageSize:          5,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1"), nil, testutil.GenPtr("assigneeCursor")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultPullRequestAssigneeEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "SGNL") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									pullRequests (first: 1) {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											assignees (first: 5, after: "assigneeCursor") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													id
													login
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_pullrequest_participant_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("SGNL"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "PullRequestParticipant",
				PageSize:          5,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1"), nil, testutil.GenPtr("participantCursor")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultPullRequestParticipantEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "SGNL") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									pullRequests (first: 1) {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											participants (first: 5, after: "participantCursor") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													id
													login
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_pullrequest_review_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("SGNL"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "PullRequestReview",
				PageSize:          5,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1"), nil, testutil.GenPtr("reviewCursor")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultPullRequestReviewEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "SGNL") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									pullRequests (first: 1) {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											latestOpinionatedReviews (first: 5, after: "reviewCursor") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													author {
														login
													}
													authorCanPushToRepository
													createdAt
													id
													state
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"default_pullrequest_commit_entity_attributes": {
			request: &github.Request{
				BaseURL:           "https://ghe-test-server",
				EnterpriseSlug:    testutil.GenPtr("SGNL"),
				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),
				EntityExternalID:  "PullRequestCommit",
				PageSize:          5,
				Token:             "Bearer Testtoken",
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("testcursor1"), nil, testutil.GenPtr("commitCursor")},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultPullRequestCommitEntityConfig(),
			},
			wantQuery: `query {
				enterprise (slug: "SGNL") {
					id
					organizations (first: 1) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							repositories (first: 1, after: "testcursor1") {
								pageInfo {
									endCursor
									hasNextPage
								}
								nodes {
									id
									pullRequests (first: 1) {
										pageInfo {
											endCursor
											hasNextPage
										}
										nodes {
											id
											commits (first: 5, after: "commitCursor") {
												pageInfo {
													endCursor
													hasNextPage
												}
												nodes {
													commit {
														author {
															email
															user {
																id
																login
															}
														}
														committedDate
														id
													}
													id
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotQuery, gotError := github.ConstructQuery(tt.request)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("gotError: %v, wantError: %v", gotError, tt.wantError)
			}

			if tt.wantQuery != "" {
				wantQueryStr, err := LoadGraphQLPayload(tt.wantQuery)
				if err != nil {
					t.Fatalf("Error normalizing query: %v", err)
				}

				gotNormalizedQuery := NormalizeQuery(gotQuery)

				wantNormalizedQuery := NormalizeQuery(wantQueryStr)

				if diff := cmp.Diff(wantNormalizedQuery, gotNormalizedQuery); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotNormalizedQuery, wantNormalizedQuery) {
					t.Fatalf("gotQuery: %v, wantQuery: %v", gotNormalizedQuery, wantNormalizedQuery)
				}
			} else if gotQuery != "" {
				t.Errorf("gotQuery: %v, wantQuery: %v", gotQuery, tt.wantQuery)
			}
		})
	}
}

// Remove multiple spaces, commas, newlines, and tabs.
func NormalizeQuery(query string) string {
	normalized := strings.ReplaceAll(query, "\\n", " ")
	normalized = strings.ReplaceAll(normalized, "\\t", " ")
	normalized = strings.ReplaceAll(normalized, ",", " ")

	// Normalize multiple spaces to a single space
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")

	return normalized
}

// Unmarshal QueryBuilder output for query string comparison.
func LoadGraphQLPayload(query string) (string, error) {
	b, err := json.Marshal(github.GraphQLPayload{Query: query})
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func CreateGraphQLCompositeCursor(
	cursorPageInfoEndCursors []*string,
	collectionID *string,
	collectionPageInfoEndCursors []*string,
) (
	nextCursor *pagination.CompositeCursor[string],
) {
	nextCursor = &pagination.CompositeCursor[string]{}
	nextCursor.CollectionID = collectionID

	if cursorPageInfoEndCursors != nil {
		cursorPageInfo := BuildPageInfo(cursorPageInfoEndCursors)

		b1, marshalErr := json.Marshal(cursorPageInfo)
		if marshalErr != nil {
			return nil
		}

		encodedCursor := base64.StdEncoding.EncodeToString(b1)
		nextCursor.Cursor = &encodedCursor
	}

	if collectionPageInfoEndCursors != nil {
		collectionPageInfo := BuildPageInfo(collectionPageInfoEndCursors)

		b2, marshalErr := json.Marshal(collectionPageInfo)
		if marshalErr != nil {
			return nil
		}

		encodedCollectionCursor := base64.StdEncoding.EncodeToString(b2)
		nextCursor.CollectionCursor = &encodedCollectionCursor
	}

	return nextCursor
}

func ValidateGraphQLCompositeCursor(
	gotCompositeCursor *pagination.CompositeCursor[string],
	wantCompositeCursor *pagination.CompositeCursor[string],
) bool {
	if gotCompositeCursor == nil && wantCompositeCursor == nil {
		return true
	} else if gotCompositeCursor == nil || wantCompositeCursor == nil {
		return false
	}

	var err *framework.Error

	var gotCursorPageInfo *github.PageInfo

	if gotCompositeCursor.Cursor != nil {
		gotCursorPageInfo, err = github.DecodePageInfo(gotCompositeCursor.Cursor)
		if err != nil {
			return false
		}
	}

	var wantCursorPageInfo *github.PageInfo
	if wantCompositeCursor.Cursor != nil {
		wantCursorPageInfo, err = github.DecodePageInfo(wantCompositeCursor.Cursor)
		if err != nil {
			return false
		}
	}

	if !ValidatePageInfo(gotCursorPageInfo, wantCursorPageInfo) {
		return false
	}

	if (gotCompositeCursor.CollectionID != nil && wantCompositeCursor.CollectionID != nil) &&
		(*gotCompositeCursor.CollectionID != *wantCompositeCursor.CollectionID) {
		return false
	}

	if (gotCompositeCursor.CollectionID == nil && wantCompositeCursor.CollectionID != nil) ||
		(gotCompositeCursor.CollectionID != nil && wantCompositeCursor.CollectionID == nil) {
		return false
	}

	var gotCollectionPageInfo *github.PageInfo

	var wantCollectionPageInfo *github.PageInfo

	if gotCompositeCursor.CollectionCursor != nil {
		decodedCollectionCursor, err := base64.StdEncoding.DecodeString(*gotCompositeCursor.CollectionCursor)
		if err != nil {
			return false
		}

		json.Unmarshal(decodedCollectionCursor, &gotCollectionPageInfo)
	}

	if wantCompositeCursor.CollectionCursor != nil {
		decodedCollectionCursor, err := base64.StdEncoding.DecodeString(*wantCompositeCursor.CollectionCursor)
		if err != nil {
			return false
		}

		json.Unmarshal(decodedCollectionCursor, &wantCollectionPageInfo)
	}

	return ValidatePageInfo(gotCollectionPageInfo, wantCollectionPageInfo)
}

func ValidatePageInfo(gotPageInfo *github.PageInfo, wantPageInfo *github.PageInfo) bool {
	if gotPageInfo == nil && wantPageInfo == nil {
		return true
	}

	if gotPageInfo == nil || wantPageInfo == nil {
		return false
	}

	if (gotPageInfo.EndCursor == nil && wantPageInfo.EndCursor != nil) ||
		(gotPageInfo.EndCursor != nil && wantPageInfo.EndCursor == nil) {
		return false
	}

	if (gotPageInfo.EndCursor != nil && wantPageInfo.EndCursor != nil) &&
		(*gotPageInfo.EndCursor != *wantPageInfo.EndCursor) {
		return false
	}

	return ValidatePageInfo(gotPageInfo.InnerPageInfo, wantPageInfo.InnerPageInfo)
}

func BuildPageInfo(endCursors []*string) *github.PageInfo {
	if endCursors == nil {
		return nil
	}

	var outermostPageInfo *github.PageInfo

	// Iterate through the list of string pointers in reverse order
	for i := len(endCursors) - 1; i >= 0; i-- {
		pageInfo := &github.PageInfo{
			EndCursor:     endCursors[i],
			InnerPageInfo: outermostPageInfo,
		}
		outermostPageInfo = pageInfo
	}

	return outermostPageInfo
}

// CollectionID is the top level attribute used in the query. For most queries this is the EnterpriseID.
// However, for Collaborators and Users, collectionID is the Organization Login attribute.
func ValidationQueryBuilder(externalID string, collectionID string, pageSize int64, endCursors []*string, organizations ...string) string {
	req := PopulateDefaultRequest(externalID, collectionID, pageSize, organizations...)

	pageInfo := BuildPageInfo(endCursors)

	// Set the organization offset if the collectionID is an organization.
	for idx, org := range organizations {
		if org == collectionID {
			if pageInfo == nil {
				pageInfo = &github.PageInfo{}
			}

			pageInfo.OrganizationOffset = idx

			break
		}
	}

	builder, _ := github.ConstructQueryBuilder(req, pageInfo)

	b, _ := builder.Build(req)

	jsonData, _ := json.Marshal(github.GraphQLPayload{Query: b})

	return string(jsonData)
}

func PopulateDefaultRequest(externalID string, collectionID string, pageSize int64, organizations ...string) *github.Request {
	req := &github.Request{
		EntityExternalID: externalID,
		PageSize:         pageSize,
		Cursor: &pagination.CompositeCursor[string]{
			CollectionID: &collectionID,
		},
	}

	if len(organizations) > 0 {
		req.Organizations = organizations
	} else {
		req.EnterpriseSlug = &collectionID
	}

	switch externalID {
	case github.Organization:
		req.EntityConfig = PopulateDefaultOrganizationEntityConfig()
	case github.OrganizationUser:
		req.EntityConfig = PopulateDefaultOrganizationUserEntityConfig()
	case github.User:
		req.EntityConfig = PopulateDefaultUserEntityConfig()
	case github.Team:
		req.EntityConfig = PopulateDefaultTeamEntityConfig()
	case github.Repository:
		req.EntityConfig = PopulateDefaultRepositoryEntityConfig()
	case github.Collaborator:
		req.EntityConfig = PopulateDefaultCollaboratorEntityConfig()
	case github.Label:
		req.EntityConfig = PopulateDefaultLabelEntityConfig()
	case github.IssueLabel:
		req.EntityConfig = PopulateDefaultIssueLabelEntityConfig()
	case github.PullRequestLabel:
		req.EntityConfig = PopulateDefaultPullRequestLabelEntityConfig()
	case github.Issue:
		req.EntityConfig = PopulateDefaultIssueEntityConfig()
	case github.IssueAssignee:
		req.EntityConfig = PopulateDefaultIssueAssigneeEntityConfig()
	case github.IssueParticipant:
		req.EntityConfig = PopulateDefaultIssueParticipantEntityConfig()
	case github.PullRequest:
		req.EntityConfig = PopulateDefaultPullRequestEntityConfig()
	case github.PullRequestAssignee:
		req.EntityConfig = PopulateDefaultPullRequestAssigneeEntityConfig()
	case github.PullRequestParticipant:
		req.EntityConfig = PopulateDefaultPullRequestParticipantEntityConfig()
	case github.PullRequestCommit:
		req.EntityConfig = PopulateDefaultPullRequestCommitEntityConfig()
	case github.PullRequestChangedFile:
		req.EntityConfig = PopulateDefaultPullRequestChangedFileEntityConfig()
	case github.PullRequestReview:
		req.EntityConfig = PopulateDefaultPullRequestReviewEntityConfig()
	}

	return req
}

func PopulateDefaultOrganizationEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.Organization,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "enterpriseId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "databaseId",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
			{
				ExternalId: "email",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "viewerIsAMember",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "viewerCanCreateTeams",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "updatedAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
	}
}

func PopulateDefaultOrganizationUserEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.OrganizationUser,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "orgId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.node.id", // userId
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "role",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: "$.node.organizationVerifiedDomainEmails",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "email",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
		},
	}
}

func PopulateDefaultUserEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.User,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "databaseId",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
			{
				ExternalId: "email",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "isViewer",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "updatedAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
	}
}

func PopulateDefaultTeamEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.Team,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "enterpriseId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "orgId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "databaseId",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
			{
				ExternalId: "slug",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "viewerCanAdminister",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "updatedAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: "$.members.edges",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.node.id",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "$.node.databaseId",
						Type:       framework.AttributeTypeInt64,
						List:       false,
					},
					{
						ExternalId: "role",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "$.node.email",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "$.node.login",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "$.node.isViewer",
						Type:       framework.AttributeTypeBool,
						List:       false,
					},
					{
						ExternalId: "$.node.updatedAt",
						Type:       framework.AttributeTypeDateTime,
						List:       false,
					},
					{
						ExternalId: "$.node.createdAt",
						Type:       framework.AttributeTypeDateTime,
						List:       false,
					},
				},
			},
			{
				ExternalId: "$.repositories.edges",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "$.node.id",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "permission",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "$.node.name",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "$.node.databaseId",
						Type:       framework.AttributeTypeInt64,
						List:       false,
					},
					{
						ExternalId: "$.node.url",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "$.node.allowUpdateBranch",
						Type:       framework.AttributeTypeBool,
						List:       false,
					},
					{
						ExternalId: "$.node.pushedAt",
						Type:       framework.AttributeTypeDateTime,
						List:       false,
					},
					{
						ExternalId: "$.node.createdAt",
						Type:       framework.AttributeTypeDateTime,
						List:       false,
					},
				},
			},
		},
	}
}

func PopulateDefaultRepositoryEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.Repository,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "enterpriseId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "orgId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "name",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "databaseId",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
			{
				ExternalId: "allowUpdateBranch",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "pushedAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: "$.collaborators.edges",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "permission",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
					{
						ExternalId: "$.node.id",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
		},
	}
}

func PopulateDefaultCollaboratorEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.Collaborator,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "databaseId",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
			{
				ExternalId: "email",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "isViewer",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "updatedAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
	}
}

func PopulateDefaultLabelEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.Label,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "repositoryId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "name",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "color",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "isDefault",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
		},
	}
}

func PopulateDefaultPullRequestLabelEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.PullRequestLabel,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "labelId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "title",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultIssueLabelEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.IssueLabel,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "labelId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "title",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultIssueEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.Issue,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "title",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "repositoryId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.author.login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "isPinned",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
	}
}

func PopulateDefaultIssueAssigneeEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.IssueAssignee,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "issueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultIssueParticipantEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.IssueParticipant,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "issueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultPullRequestEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.PullRequest,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "title",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.author.login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.baseRepository.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.headRepository.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "closed",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
		},
	}
}

func PopulateDefaultPullRequestChangedFileEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.PullRequestChangedFile,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "pullRequestId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "path",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "changeType",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultPullRequestAssigneeEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.PullRequestAssignee,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "pullRequestId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultPullRequestParticipantEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.PullRequestParticipant,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "uniqueId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "pullRequestId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultPullRequestCommitEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.PullRequestCommit,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "pullRequestId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.commit.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.commit.committedDate",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
			{
				ExternalId: "$.commit.author.email",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.commit.author.user.id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.commit.author.user.login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

func PopulateDefaultPullRequestReviewEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.PullRequestReview,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "pullRequestId",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.author.login",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "state",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "authorCanPushToRepository",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "createdAt",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
	}
}

func PopulateDefaultSecretScanningAlertEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: github.SecretScanningAlert,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "number",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
			{
				ExternalId: "state",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "secret_type",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "$.repository.node_id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "created_at",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
	}
}
