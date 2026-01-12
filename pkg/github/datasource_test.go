// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package github_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/github"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestParseGraphQLResponse(t *testing.T) {
	tests := map[string]struct {
		body             []byte
		wantObjects      []map[string]any
		entityExternalID string
		wantNextCursor   *pagination.CompositeCursor[string]
		wantErr          *framework.Error
	}{
		"single_page": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"login": "ArvindOrg2"
								},
								{
									"id": "MDEyOk9yZ2FuaXphdGlv333=",
									"login": "ArvindOrg1"
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "Organization",
			wantObjects: []map[string]any{
				{
					"id":           "MDEyOk9yZ2FuaXphdGlvbjk=",
					"login":        "ArvindOrg2",
					"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
				},
				{
					"id":           "MDEyOk9yZ2FuaXphdGlv333=",
					"login":        "ArvindOrg1",
					"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
				},
			},
			wantNextCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
		},
		"invalid_object_structure": {
			body: []byte(`{
				"enterprise": {
					"id": "MDEwOkVudGVycHJpc2Ux",
					"name": "SGNL",
					"organizations": {
						"pageInfo": {
							"hasNextPage": true,
							"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
						},
						"nodes": [
							{
								"login": "ArvindOrg2"
							},
							{
								"login": "ArvindOrg1"
							}
						]
					}
				}
			}`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: Data not found.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"error_parsing": {
			body: []byte(`{
				"errors": [
					{
						"path": [
							"query",
							"enterprise",
							"organizations",
							"nodes",
							"repositories",
							"nodes",
							"pullRequests",
							"nodes",
							"createdVEmail"
						],
						"extensions": {
							"code": "undefinedField",
							"typeName": "PullRequest",
							"fieldName": "createdVEmail"
						},
						"locations": [
							{
								"line": 33,
								"column": 33
							}
						],
						"message": "Field 'createdVEmail' doesn't exist on type 'PullRequest'"
					},
					{
						"path": [
							"query",
							"enterprise",
							"organizations",
							"nodes",
							"repositories",
							"nodes",
							"pullRequests",
							"nodes",
							"deletion"
						],
						"extensions": {
							"code": "undefinedField",
							"typeName": "PullRequest",
							"fieldName": "deletion"
						},
						"locations": [
							{
								"line": 34,
								"column": 33
							}
						],
						"message": "Field 'deletion' doesn't exist on type 'PullRequest'"
					}
				]
			}`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to get the datasource response: [Field 'createdVEmail' doesn't exist on type 'PullRequest' Field 'deletion' doesn't exist on type 'PullRequest'].",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_errors_structure": {
			body: []byte(`{
				"errors": {}
			}`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal object into Go struct field DatasourceResponse.errors of type []github.ErrorInfo.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"errors_present_with_length_0": {
			body: []byte(`{
				"errors": []
			}`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to get the datasource response. Unexpected error format: Errors array is empty.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"both_errors_and_data_present": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"login": "ArvindOrg2"
								},
								{
									"id": "MDEyOk9yZ2FuaXphdGlv333=",
									"login": "ArvindOrg1"
								}
							]
						}
					}
				},
				"errors": [
					{
						"path": [
							"query",
							"enterprise",
							"organizations",
							"nodes",
							"repositories",
							"nodes",
							"pullRequests",
							"nodes",
							"createdVEmail"
						],
						"extensions": {
							"code": "undefinedField",
							"typeName": "PullRequest",
							"fieldName": "createdVEmail"
						},
						"locations": [
							{
								"line": 33,
								"column": 33
							}
						],
						"message": "Field 'createdVEmail' doesn't exist on type 'PullRequest'"
					},
					{
						"path": [
							"query",
							"enterprise",
							"organizations",
							"nodes",
							"repositories",
							"nodes",
							"pullRequests",
							"nodes",
							"deletion"
						],
						"extensions": {
							"code": "undefinedField",
							"typeName": "PullRequest",
							"fieldName": "deletion"
						},
						"locations": [
							{
								"line": 34,
								"column": 33
							}
						],
						"message": "Field 'deletion' doesn't exist on type 'PullRequest'"
					}
				]
			}`),
			entityExternalID: "Organization",
			wantErr: &framework.Error{
				Message: "Failed to get the datasource response: [Field 'createdVEmail' doesn't exist on type 'PullRequest' Field 'deletion' doesn't exist on type 'PullRequest'].",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"invalid_errors_field_format": {
			body: []byte(`{
				"errors": [
					{
						"path": []
					}
				]
			}`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to get the datasource response. Unexpected error format: message is missing.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_error_message_field_type": {
			body: []byte(`{
				"errors": [
					{
						"path": [
							"query",
							"enterprise",
							"organizations",
							"nodes",
							"repositories",
							"nodes",
							"pullRequests",
							"nodes",
							"createdVEmail"
						],
						"extensions": {
							"code": "undefinedField",
							"typeName": "PullRequest",
							"fieldName": "createdVEmail"
						},
						"locations": [
							{
								"line": 33,
								"column": 33
							}
						],
						"message": 50
					}
				]
			}`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal number into Go struct field ErrorInfo.errors.message of type string.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"response_structure_missing_data": {
			body:             []byte(`{}`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: Data not found.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"response_structure_missing_enterprise": {
			body:             []byte(`{ "data": {} }`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to parse the datasource response: Enterprise not found.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"response_structure_missing_organizations": {
			body:             []byte(`{ "data": { "enterprise": {} } }`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to parse the datasource response: Organizations not found.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"response_structure_missing_nodes": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=",
								"hasNextPage": true
							}
						}
					}
				}
			}`),
			entityExternalID: "Repository",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response. Organization not found: unexpected end of JSON input.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"response_structure_missing_repositories": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEw"
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=",
								"hasNextPage": true
							}
						}
					}
				}
			}`),
			entityExternalID: "Repository",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to parse the datasource response: Repositories not found.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"response_structure_missing_pageinfo": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEw",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "Repository",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to validate LeafPageInfo: PageInfo is nil.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"response_structure_missing_issues": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx"
											}
										]
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "Issue",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to parse the datasource response: Issues not found.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"repositories_with_empty_nodes": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEw",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=",
								"hasNextPage": true
							}
						}
					}
				}
			}`),
			entityExternalID: "Repository",
			wantErr:          nil,
			wantObjects:      []map[string]any{},
			wantNextCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
		},
		"invalid_objects": {
			body: []byte(`{
				"enterprise": [
					"id": "MDEwOkVudGVycHJpc2Ux",
					"name": testutil.GenPtr("SGNL"),
					"organizations": {
						"pageInfo": [
							"hasNextPage": true,
							"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
						},
						"nodes": [
							{
								"login": "ArvindOrg2"
							},
							{
								"login": "ArvindOrg1"
							}
						]
					}
				}
			}`),
			entityExternalID: "Organization",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: invalid character ':' after array element.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"empty_organization_collection_for_sub_entities": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": []
						}
					}
				}
			}`),
			entityExternalID: "Team",
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to parse the datasource response: Organization collection is length 0, expected 1.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"empty_repository_collection_for_sub_entities": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"repositories": {
										"pageInfo": {
											"hasNextPage": false,
											"endCursor": null
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "Collaborator",
			wantObjects:      []map[string]any{},
			wantNextCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
			wantErr: nil,
		},
		"empty_issue_collection_for_sub_entities": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"repositories": {
										"pageInfo": {
											"hasNextPage": true,
											"endCursor": "repoCursor"
										},
										"nodes": [
											{
												"id": "repoId",
												"issues": {
													"pageInfo": {
														"hasNextPage": false,
														"endCursor": null
													},
													"nodes": []
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "IssueAssignee",
			wantObjects:      []map[string]any{},
			wantNextCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("repoCursor")},
				nil,
				nil,
			),
			wantErr: nil,
		},
		"empty_pullRequest_collection_for_sub_entities": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"repositories": {
										"pageInfo": {
											"hasNextPage": true,
											"endCursor": "repoCursor"
										},
										"nodes": [
											{
												"id": "repoId",
												"pullRequests": {
													"pageInfo": {
														"hasNextPage": false,
														"endCursor": null
													},
													"nodes": []
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "PullRequestParticipant",
			wantObjects:      []map[string]any{},
			wantNextCursor: CreateGraphQLCompositeCursor(
				[]*string{nil, testutil.GenPtr("repoCursor")},
				nil,
				nil,
			),
			wantErr: nil,
		},
		"multiple_repositories_in_collection_for_sub_entities": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"repositories": {
										"pageInfo": {
											"hasNextPage": false,
											"endCursor": null
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx"
											},
											{
												"id": "MDEwOlJlcG9zaXRvcnkx"
											}
										]
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "Collaborator",
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Repository container length is: 2, expected: 1.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"malformed_issue_collection_for_sub_entities": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"repositories": {
										"pageInfo": {
											"hasNextPage": true,
											"endCursor": "repoCursor"
										},
										"nodes": [
											{
												"id": "repoId",
												"issues": {
													"pageInfo": {
														"hasNextPage": false,
														"endCursor": null
													},
													"edges": []
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "IssueParticipant",
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response. Issue not found: unexpected end of JSON input.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"user_entity_with_edges_instead_of_nodes": {
			body: []byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"membersWithRole": {
										"pageInfo": {
											"hasNextPage": true,
											"endCursor": "repoCursor"
										},
										"edges": [
											{
												"role": "ADMIN",
												"node": {
													"id": "user1"
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "User",
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: unexpected end of JSON input.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"valid_organizationuser_entity": {
			body: []byte(`{
				"data": {
					"organization": {
						"id": "orgId1",
						"membersWithRole": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"edges": [
								{
									"role": "ADMIN",
									"node": {
										"id": "user1",
										"login": "isabella-sgnl"
									}
								},
								{
									"role": "MEMBER",
									"node": {
										"id": "user2",
										"login": "arooxa"
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "OrganizationUser",
			wantObjects: []map[string]any{
				{
					"role":     "ADMIN",
					"orgId":    "orgId1",
					"uniqueId": "orgId1-user1",
					"node": map[string]any{
						"id":    "user1",
						"login": "isabella-sgnl",
					},
				},
				{
					"role":     "MEMBER",
					"orgId":    "orgId1",
					"uniqueId": "orgId1-user2",
					"node": map[string]any{
						"id":    "user2",
						"login": "arooxa",
					},
				},
			},
			wantNextCursor: CreateGraphQLCompositeCursor(
				[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				nil,
				nil,
			),
		},
		"organizationuser_entity_missing_orgId": {
			body: []byte(`{
				"data": {
					"organization": {
						"membersWithRole": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"edges": [
								{
									"role": "ADMIN",
									"node": {
										"id": "user1",
										"login": "isabella-sgnl"
									}
								},
								{
									"role": "MEMBER",
									"node": {
										"id": "user2",
										"login": "arooxa"
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "OrganizationUser",
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Organization is nil or orgID is missing for the OrganizationUser entity.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"malformed_organizationuser_entity_with_nodes": {
			body: []byte(`{
				"data": {
					"organization": {
						"id": "orgId1",
						"membersWithRole": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"role": "ADMIN",
									"node": {
										"id": "user1",
										"login": "isabella-sgnl"
									}
								},
								{
									"role": "MEMBER",
									"node": {
										"id": "user2",
										"login": "arooxa"
									}
								}
							]
						}
					}
				}
			}`),
			entityExternalID: "OrganizationUser",
			wantObjects:      nil,
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: unexpected end of JSON input.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := github.ParseGraphQLResponse(tt.body, tt.entityExternalID, nil, 0)

			if diff := cmp.Diff(tt.wantObjects, gotObjects); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !ValidateGraphQLCompositeCursor(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestParseRESTResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		wantObjects    []map[string]any
		links          []string
		wantNextCursor *pagination.CompositeCursor[string]
		wantErr        *framework.Error
	}{
		"single_page": {
			body: []byte(`[
				{
					"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
					"alert": "ALERT1"
				},
				{
					"id": "MDEyOk9yZ2FuaXphdGlv333=",
					"alert": "ALERT2"
				}
			]`),
			links: []string{
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=1>; rel="prev"`,
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3>; rel="next"`,
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3>; rel="last"`,
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=1>; rel="first"`,
			},
			wantObjects: []map[string]any{
				{
					"id":    "MDEyOk9yZ2FuaXphdGlvbjk=",
					"alert": "ALERT1",
				},
				{
					"id":    "MDEyOk9yZ2FuaXphdGlv333=",
					"alert": "ALERT2",
				},
			},
			wantNextCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3"),
			},
		},
		"invalid_object_structure": {
			body: []byte(`{
				[
					{
						"id": "MDEw"
					},
					{
						"id": "MDEw"
					}
				]
			}`),
			links: []string{
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=1>; rel="prev"`,
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3>; rel="next"`,
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3>; rel="last"`,
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=1>; rel="first"`,
			},
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: invalid character '[' looking for beginning of object key string.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"missing_links_want_error": {
			body: []byte(`[
				{
					"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
					"alert": "ALERT1"
				},
				{
					"id": "MDEyOk9yZ2FuaXphdGlv333=",
					"alert": "ALERT2"
				}
			]`),
			links:       []string{},
			wantObjects: nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to parse the datasource response: Link header is empty or not found.",
				Code:    adapter_api_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"next_link_missing_want_next_cursor_with_collectionid": {
			body: []byte(`[
				{
					"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
					"alert": "ALERT1"
				},
				{
					"id": "MDEyOk9yZ2FuaXphdGlv333=",
					"alert": "ALERT2"
				}
			]`),
			links: []string{
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=1>; rel="prev"`,
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3>; rel="last"`,
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=1>; rel="first"`,
			},
			wantObjects: []map[string]any{
				{
					"id":    "MDEyOk9yZ2FuaXphdGlvbjk=",
					"alert": "ALERT1",
				},
				{
					"id":    "MDEyOk9yZ2FuaXphdGlv333=",
					"alert": "ALERT2",
				},
			},
			wantNextCursor: nil,
		},
		"only_next_link_present_want_next_cursor": {
			body: []byte(`[
				{
					"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
					"alert": "ALERT1"
				},
				{
					"id": "MDEyOk9yZ2FuaXphdGlv333=",
					"alert": "ALERT2"
				}
			]`),
			links: []string{
				`<https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3>; rel="next"`,
			},
			wantObjects: []map[string]any{
				{
					"id":    "MDEyOk9yZ2FuaXphdGlvbjk=",
					"alert": "ALERT1",
				},
				{
					"id":    "MDEyOk9yZ2FuaXphdGlv333=",
					"alert": "ALERT2",
				},
			},
			wantNextCursor: &pagination.CompositeCursor[string]{
				Cursor: testutil.GenPtr("https://test-instance.com/api/v3/repositories/1/issues?per_page=1&page=3"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotNextCursor, gotErr := github.ParseRESTResponse(tt.body, tt.links, 0, 0)

			if diff := cmp.Diff(tt.wantObjects, gotObjects); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
			}

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotNextCursor, tt.wantNextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotNextCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetOrganizationPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context      context.Context
		request      *github.Request
		wantRes      *github.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:          server.URL,
				Token:            "Bearer Testtoken",
				PageSize:         1,
				EntityExternalID: "Organization",
				EnterpriseSlug:   testutil.GenPtr("SGNL"),

				IsEnterpriseCloud: false,
				APIVersion:        testutil.GenPtr("v3"),

				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                   "MDEyOk9yZ2FuaXphdGlvbjk=",
						"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
						"databaseId":           float64(9),
						"email":                nil,
						"login":                "ArvindOrg1",
						"viewerIsAMember":      true,
						"viewerCanCreateTeams": true,
						"updatedAt":            "2024-02-02T23:20:22Z",
						"createdAt":            "2024-02-02T23:20:22Z",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "Organization",
					fields.FieldRequestPageSize:         int64(1),
				},
				{
					"level":                             "info",
					"msg":                               "Sending request to datasource",
					fields.FieldRequestEntityExternalID: "Organization",
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldRequestURL:              server.URL + "/api/graphql",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "Organization",
					fields.FieldRequestPageSize:         int64(1),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(1),
					fields.FieldResponseNextCursor: map[string]any{
						"cursor": "eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6IlkzVnljMjl5T25ZeU9wS3FRWEoyYVc1a1QzSm5NUWs9Iiwib3JnYW5pemF0aW9uT2Zmc2V0IjowLCJJbm5lclBhZ2VJbmZvIjpudWxsfQ==",
					},
				},
			},
		},
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "Organization",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                   "MDEyOk9yZ2FuaXphdGlvbjEw",
						"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
						"databaseId":           float64(10),
						"email":                nil,
						"login":                "ArvindOrg2",
						"viewerIsAMember":      true,
						"viewerCanCreateTeams": true,
						"updatedAt":            "2024-02-15T17:00:12Z",
						"createdAt":            "2024-02-15T17:00:12Z",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "Organization",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                   "MDEyOk9yZ2FuaXphdGlvbjU=",
						"enterpriseId":         "MDEwOkVudGVycHJpc2Ux",
						"databaseId":           float64(5),
						"email":                nil,
						"login":                "EnterpriseServerOrg",
						"viewerIsAMember":      true,
						"viewerCanCreateTeams": true,
						"updatedAt":            "2024-01-28T23:00:00Z",
						"createdAt":            "2024-01-28T22:59:59Z",
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.context)

			gotRes, gotErr := githubClient.GetPage(ctxWithLogger, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes.Objects, tt.wantRes.Objects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotRes.Objects, tt.wantRes.Objects)
			}

			if !ValidateGraphQLCompositeCursor(gotRes.NextCursor, tt.wantRes.NextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotRes.NextCursor, tt.wantRes.NextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}

func TestGetOrganizationUserPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "OrganizationUser",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationUserEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"role":     "ADMIN",
						"orgId":    "MDEyOk9yZ2FuaXphdGlvbjU=",
						"uniqueId": "MDEyOk9yZ2FuaXphdGlvbjU=-MDQ6VXNlcjQ=",
						"node": map[string]any{
							"id": "MDQ6VXNlcjQ=",
							"organizationVerifiedDomainEmails": []any{
								map[string]any{
									"email": "arvind@sgnldemos.com",
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
					testutil.GenPtr("ArvindOrg1"),
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				),
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "OrganizationUser",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationUserEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
					testutil.GenPtr("ArvindOrg1"),
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"role":     "MEMBER",
						"orgId":    "MDEyOk9yZ2FuaXphdGlvbjU=",
						"uniqueId": "MDEyOk9yZ2FuaXphdGlvbjU=-MDQ6VXNlcjk=",
						"node": map[string]any{
							"id": "MDQ6VXNlcjk=",
							"organizationVerifiedDomainEmails": []any{
								map[string]any{
									"email": "isabella@sgnldemos.com",
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					nil,
					nil,
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				),
			},
			wantErr: nil,
		},
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "OrganizationUser",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationUserEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					nil,
					nil,
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"role":     "ADMIN",
						"orgId":    "MDEyOk9yZ2FuaXphdGlvbjEy",
						"uniqueId": "MDEyOk9yZ2FuaXphdGlvbjEy-MDQ6VXNlcjQ=",
						"node": map[string]any{
							"id":                               "MDQ6VXNlcjQ=",
							"organizationVerifiedDomainEmails": []any{},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
					testutil.GenPtr("ArvindOrg2"),
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				),
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "OrganizationUser",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationUserEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
					testutil.GenPtr("ArvindOrg2"),
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: CreateGraphQLCompositeCursor(
					nil,
					nil,
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
				),
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes.Objects, tt.wantRes.Objects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotRes.Objects, tt.wantRes.Objects)
			}

			if !ValidateGraphQLCompositeCursor(gotRes.NextCursor, tt.wantRes.NextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotRes.NextCursor, tt.wantRes.NextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetTeamPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Team",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultTeamEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                  "MDQ6VGVhbTI=",
						"enterpriseId":        "MDEwOkVudGVycHJpc2Ux",
						"orgId":               "MDEyOk9yZ2FuaXphdGlvbjk=",
						"databaseId":          float64(2),
						"slug":                "secret-team-1",
						"viewerCanAdminister": true,
						"updatedAt":           "2024-02-02T23:21:54Z",
						"createdAt":           "2024-02-02T23:21:54Z",
						"members": map[string]any{
							"edges": []any{
								map[string]any{
									"role": "MAINTAINER",
									"node": map[string]any{
										"id":         "MDQ6VXNlcjY=",
										"databaseId": float64(6),
										"email":      "",
										"login":      "arvind",
										"isViewer":   true,
										"updatedAt":  "2024-01-31T05:09:26Z",
										"createdAt":  "2024-01-28T23:28:03Z",
									},
								},
							},
						},
						"repositories": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "ADMIN",
									"node": map[string]any{
										"id":                "MDEwOlJlcG9zaXRvcnk2",
										"name":              "arvindrepo2",
										"databaseId":        float64(6),
										"url":               "https://ghe-test-server/ArvindOrg1/arvindrepo2",
										"allowUpdateBranch": false,
										"pushedAt":          "2024-02-02T23:22:33Z",
										"createdAt":         "2024-02-02T23:22:32Z",
									},
								},
							},
						},
					},
					{
						"id":                  "MDQ6VGVhbTE=",
						"enterpriseId":        "MDEwOkVudGVycHJpc2Ux",
						"orgId":               "MDEyOk9yZ2FuaXphdGlvbjk=",
						"databaseId":          float64(1),
						"slug":                "team1",
						"viewerCanAdminister": true,
						"updatedAt":           "2024-02-02T23:21:02Z",
						"createdAt":           "2024-02-02T23:21:02Z",
						"members": map[string]any{
							"edges": []any{
								map[string]any{
									"role": "MEMBER",
									"node": map[string]any{
										"id":         "MDQ6VXNlcjQ=",
										"databaseId": float64(4),
										"email":      "",
										"login":      "isabella",
										"isViewer":   false,
										"updatedAt":  "2024-02-22T18:43:44Z",
										"createdAt":  "2024-01-28T22:02:26Z",
									},
								},
								map[string]any{
									"role": "MAINTAINER",
									"node": map[string]any{
										"id":         "MDQ6VXNlcjY=",
										"databaseId": float64(6),
										"email":      "",
										"login":      "arvind",
										"isViewer":   true,
										"updatedAt":  "2024-01-31T05:09:26Z",
										"createdAt":  "2024-01-28T23:28:03Z",
									},
								},
							},
						},
						"repositories": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "MAINTAIN",
									"node": map[string]any{
										"id":                "MDEwOlJlcG9zaXRvcnk1",
										"name":              "arvindrepo1",
										"databaseId":        float64(5),
										"url":               "https://ghe-test-server/ArvindOrg1/arvindrepo1",
										"allowUpdateBranch": false,
										"pushedAt":          "2024-02-02T23:22:20Z",
										"createdAt":         "2024-02-02T23:22:20Z",
									},
								},
								map[string]any{
									"permission": "WRITE",
									"node": map[string]any{
										"id":                "MDEwOlJlcG9zaXRvcnk2",
										"name":              "arvindrepo2",
										"databaseId":        float64(6),
										"url":               "https://ghe-test-server/ArvindOrg1/arvindrepo2",
										"allowUpdateBranch": false,
										"pushedAt":          "2024-02-02T23:22:33Z",
										"createdAt":         "2024-02-02T23:22:32Z",
									},
								},
								map[string]any{
									"permission": "READ",
									"node": map[string]any{
										"id":                "MDEwOlJlcG9zaXRvcnk3",
										"name":              "arvindrepo3",
										"databaseId":        float64(7),
										"url":               "https://ghe-test-server/ArvindOrg1/arvindrepo3",
										"allowUpdateBranch": false,
										"pushedAt":          "2024-02-02T23:22:45Z",
										"createdAt":         "2024-02-02T23:22:45Z",
									},
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Team",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultTeamEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Team",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultTeamEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                  "MDQ6VGVhbTM=",
						"enterpriseId":        "MDEwOkVudGVycHJpc2Ux",
						"orgId":               "MDEyOk9yZ2FuaXphdGlvbjU=",
						"databaseId":          float64(3),
						"slug":                "random-team-1",
						"viewerCanAdminister": true,
						"updatedAt":           "2024-02-16T04:26:06Z",
						"createdAt":           "2024-02-16T04:26:06Z",
						"members": map[string]any{
							"edges": []any{
								map[string]any{
									"role": "MAINTAINER",
									"node": map[string]any{
										"id":         "MDQ6VXNlcjY=",
										"databaseId": float64(6),
										"email":      "",
										"login":      "arvind",
										"isViewer":   true,
										"updatedAt":  "2024-01-31T05:09:26Z",
										"createdAt":  "2024-01-28T23:28:03Z",
									},
								},
							},
						},
						"repositories": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "MAINTAIN",
									"node": map[string]any{
										"id":                "MDEwOlJlcG9zaXRvcnkx",
										"name":              "enterprise_repo1",
										"databaseId":        float64(1),
										"url":               "https://ghe-test-server/EnterpriseServerOrg/enterprise_repo1",
										"allowUpdateBranch": false,
										"pushedAt":          "2024-02-02T23:17:27Z",
										"createdAt":         "2024-02-02T23:17:26Z",
									},
								},
							},
						},
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes.Objects, tt.wantRes.Objects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotRes.Objects, tt.wantRes.Objects)
			}

			if !ValidateGraphQLCompositeCursor(gotRes.NextCursor, tt.wantRes.NextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotRes.NextCursor, tt.wantRes.NextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetRepositoryPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Repository",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultRepositoryEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                "MDEwOlJlcG9zaXRvcnk1",
						"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
						"orgId":             "MDEyOk9yZ2FuaXphdGlvbjk=",
						"name":              "arvindrepo1",
						"databaseId":        float64(5),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:22:20Z",
						"createdAt":         "2024-02-02T23:22:20Z",
						"collaborators": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "ADMIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjQ=",
									},
								},
								map[string]any{
									"permission": "MAINTAIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjY=",
									},
								},
							},
						},
					},
					{
						"id":                "MDEwOlJlcG9zaXRvcnk2",
						"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
						"orgId":             "MDEyOk9yZ2FuaXphdGlvbjk=",
						"name":              "arvindrepo2",
						"databaseId":        float64(6),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:22:33Z",
						"createdAt":         "2024-02-02T23:22:32Z",
						"collaborators": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "ADMIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjQ=",
									},
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEG")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Repository",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultRepositoryEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEG")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                "MDEwOlJlcG9zaXRvcnk3",
						"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
						"orgId":             "MDEyOk9yZ2FuaXphdGlvbjk=",
						"name":              "arvindrepo3",
						"databaseId":        float64(7),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:22:45Z",
						"createdAt":         "2024-02-02T23:22:45Z",
						"collaborators":     map[string]any{"edges": []any{}},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Repository",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultRepositoryEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Repository",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultRepositoryEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
						"orgId":             "MDEyOk9yZ2FuaXphdGlvbjU=",
						"name":              "enterprise_repo1",
						"databaseId":        float64(1),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:17:27Z",
						"createdAt":         "2024-02-02T23:17:26Z",
						"collaborators":     map[string]any{"edges": []any{}},
					},
					{
						"id":                "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
						"orgId":             "MDEyOk9yZ2FuaXphdGlvbjU=",
						"name":              "enterprise_repo2",
						"databaseId":        float64(2),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:17:42Z",
						"createdAt":         "2024-02-02T23:17:41Z",
						"collaborators": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "MAINTAIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjY=",
									},
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{
						testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="),
						testutil.GenPtr("Y3Vyc29yOnYyOpEC"),
					},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Repository",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{
						testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="),
						testutil.GenPtr("Y3Vyc29yOnYyOpEC"),
					},
					nil,
					nil,
				),
				EntityConfig: PopulateDefaultRepositoryEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                "MDEwOlJlcG9zaXRvcnkz",
						"enterpriseId":      "MDEwOkVudGVycHJpc2Ux",
						"orgId":             "MDEyOk9yZ2FuaXphdGlvbjU=",
						"name":              "enterprise_repo3",
						"databaseId":        float64(3),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:18:01Z",
						"createdAt":         "2024-02-02T23:18:01Z",
						"collaborators":     map[string]any{"edges": []any{}},
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes.Objects, tt.wantRes.Objects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotRes.Objects, tt.wantRes.Objects)
			}

			if !ValidateGraphQLCompositeCursor(gotRes.NextCursor, tt.wantRes.NextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotRes.NextCursor, tt.wantRes.NextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetRepositoryPageWithOrganizations(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_page_of_first_organization": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Repository",
				Organizations:         []string{"arvindorg1", "arvindorg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 180,
				EntityConfig:          PopulateDefaultRepositoryEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                "MDEwOlJlcG9zaXRvcnk1",
						"name":              "arvindrepo1",
						"databaseId":        float64(5),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:22:20Z",
						"createdAt":         "2024-02-02T23:22:20Z",
						"orgId":             "O_kgDOCPwuWw",
						"collaborators": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "ADMIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjQ=",
									},
								},
								map[string]any{
									"permission": "MAINTAIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjY=",
									},
								},
							},
						},
					},
					{
						"id":                "MDEwOlJlcG9zaXRvcnk2",
						"name":              "arvindrepo2",
						"databaseId":        float64(6),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:22:33Z",
						"createdAt":         "2024-02-02T23:22:32Z",
						"orgId":             "O_kgDOCPwuWw",
						"collaborators": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "ADMIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjQ=",
									},
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpEG")},
					nil,
					nil,
				),
			},
		},
		"last_page_of_first_organization": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Repository",
				Organizations:         []string{"arvindorg1", "arvindorg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 180,
				EntityConfig:          PopulateDefaultRepositoryEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6IlkzVnljMjl5T25ZeU9wRUciLCJvcmdhbml6YXRpb25PZmZzZXQiOjAsIklubmVyUGFnZUluZm8iOm51bGx9"),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                "MDEwOlJlcG9zaXRvcnk1",
						"name":              "arvindrepo3",
						"databaseId":        float64(7),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:22:20Z",
						"createdAt":         "2024-02-02T23:22:20Z",
						"orgId":             "O_kgDOCPwuWw",
						"collaborators": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "ADMIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjQ=",
									},
								},
							},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6bnVsbCwib3JnYW5pemF0aW9uT2Zmc2V0IjoxLCJJbm5lclBhZ2VJbmZvIjpudWxsfQ=="),
				},
			},
		},
		"last_page_of_second_organization": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Repository",
				Organizations:         []string{"arvindorg1", "arvindorg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 180,
				EntityConfig:          PopulateDefaultRepositoryEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6bnVsbCwib3JnYW5pemF0aW9uT2Zmc2V0IjoxLCJJbm5lclBhZ2VJbmZvIjpudWxsfQ=="),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                "MDEwOlJlcG9zaXRvcnabc",
						"name":              "arvindrepo4",
						"databaseId":        float64(9),
						"allowUpdateBranch": false,
						"pushedAt":          "2024-02-02T23:22:20Z",
						"createdAt":         "2024-02-02T23:22:20Z",
						"orgId":             "O_kgDOCPwuXxyz",
						"collaborators": map[string]any{
							"edges": []any{
								map[string]any{
									"permission": "ADMIN",
									"node": map[string]any{
										"id": "MDQ6VXNlcjQ=",
									},
								},
							},
						},
					},
				},
				NextCursor: nil,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes.Objects, tt.wantRes.Objects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotRes.Objects, tt.wantRes.Objects)
			}

			if !ValidateGraphQLCompositeCursor(gotRes.NextCursor, tt.wantRes.NextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotRes.NextCursor, tt.wantRes.NextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetUserPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "User",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 90,
				EntityConfig:          PopulateDefaultUserEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"databaseId":   float64(4),
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"email":        "",
						"login":        "arooxa",
						"isViewer":     true,
						"updatedAt":    "2024-03-08T04:18:47Z",
						"createdAt":    "2024-03-08T04:18:47Z",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "User",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultUserEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjk=",
						"databaseId":   float64(9),
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"email":        "",
						"login":        "isabella-sgnl",
						"isViewer":     false,
						"updatedAt":    "2024-03-08T19:28:13Z",
						"createdAt":    "2024-03-08T17:52:21Z",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "User",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultUserEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"databaseId":   float64(4),
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"email":        "",
						"login":        "arooxa",
						"isViewer":     true,
						"updatedAt":    "2024-03-08T04:18:47Z",
						"createdAt":    "2024-03-08T04:18:47Z",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU="), testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "User",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultUserEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU="), testutil.GenPtr("Y3Vyc29yOnYyOpEE")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				// This is an empty list of objects because it is an organization with no users.
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes.Objects, tt.wantRes.Objects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotRes.Objects, tt.wantRes.Objects)
			}

			if !ValidateGraphQLCompositeCursor(gotRes.NextCursor, tt.wantRes.NextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotRes.NextCursor, tt.wantRes.NextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetCollaboratorPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Collaborator",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultCollaboratorEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"databaseId":   float64(4),
						"email":        "",
						"login":        "arooxa",
						"isViewer":     true,
						"updatedAt":    "2024-03-08T04:18:47Z",
						"createdAt":    "2024-03-08T04:18:47Z",
					},
					{
						"id":           "MDQ6VXNlcjk=",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"databaseId":   float64(9),
						"email":        "",
						"login":        "isabella-sgnl",
						"isViewer":     false,
						"updatedAt":    "2024-03-08T19:28:13Z",
						"createdAt":    "2024-03-08T17:52:21Z",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Collaborator",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultCollaboratorEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"databaseId":   float64(4),
						"email":        "",
						"login":        "arooxa",
						"isViewer":     true,
						"updatedAt":    "2024-03-08T04:18:47Z",
						"createdAt":    "2024-03-08T04:18:47Z",
					},
					{
						"id":           "MDQ6VXNlcjk=",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"databaseId":   float64(9),
						"email":        "",
						"login":        "isabella-sgnl",
						"isViewer":     false,
						"updatedAt":    "2024-03-08T19:28:13Z",
						"createdAt":    "2024-03-08T17:52:21Z",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEJ")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Collaborator",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultCollaboratorEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEJ")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjEw",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"databaseId":   float64(10),
						"email":        "",
						"login":        "r-rakshith",
						"isViewer":     false,
						"updatedAt":    "2024-03-08T17:53:47Z",
						"createdAt":    "2024-03-08T17:52:54Z",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "Collaborator",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultCollaboratorEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetLabelPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// This fetches labels [1, 8] of repo 1/2 for org 1/2.
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              8,
				EntityExternalID:      "Label",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultLabelEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDU6TGFiZWwx",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "bug",
						"color":        "d73a4a",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwy",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "documentation",
						"color":        "0075ca",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwz",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "duplicate",
						"color":        "cfd3d7",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWw0",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "enhancement",
						"color":        "a2eeef",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWw1",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "good first issue",
						"color":        "7057ff",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWw2",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "help wanted",
						"color":        "008672",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWw3",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "invalid",
						"color":        "e4e669",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWw4",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "question",
						"color":        "d876e3",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("OA")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// This fetches label 9 of repo 1/2 for org 1/2.
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              8,
				EntityExternalID:      "Label",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("OA")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDU6TGFiZWw5",
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "wontfix",
						"color":        "ffffff",
						"createdAt":    "2024-03-08T18:51:30Z",
						"isDefault":    true,
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// This fetches labels [1, 8] of repo 2/2 for org 1/2.
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              8,
				EntityExternalID:      "Label",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDU6TGFiZWwxMA==",
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "bug",
						"color":        "d73a4a",
						"createdAt":    "2024-03-08T18:51:44Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwxMQ==",
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "documentation",
						"color":        "0075ca",
						"createdAt":    "2024-03-08T18:51:44Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwxMg==",
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "duplicate",
						"color":        "cfd3d7",
						"createdAt":    "2024-03-08T18:51:44Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwxMw==",
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "enhancement",
						"color":        "a2eeef",
						"createdAt":    "2024-03-08T18:51:44Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwxNA==",
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "good first issue",
						"color":        "7057ff",
						"createdAt":    "2024-03-08T18:51:44Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwxNQ==",
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "help wanted",
						"color":        "008672",
						"createdAt":    "2024-03-08T18:51:44Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwxNg==",
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "invalid",
						"color":        "e4e669",
						"createdAt":    "2024-03-08T18:51:44Z",
						"isDefault":    true,
					},
					{
						"id":           "MDU6TGFiZWwxNw==",
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"name":         "question",
						"color":        "d876e3",
						"createdAt":    "2024-03-08T18:51:44Z",
						"isDefault":    true,
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("OA")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetIssueLabelPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// These is 1 issue on label 1/3 for repo 1/2 for org 1/2.
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueLabelEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDU6SXNzdWUz",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"labelId":      "MDU6TGFiZWwx",
						"uniqueId":     "MDU6TGFiZWwx-MDU6SXNzdWUz",
						"title":        "issue1",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("MQ")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// These are no issues on label 2/3 for repo 1/2 for org 1/2.
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("MQ")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("Mg")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// There are 2 issues for label 3/3 on repo 1/2 for org 1/2.
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("Mg")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDU6SXNzdWUz",
						"title":        "issue1",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"labelId":      "MDU6TGFiZWwy",
						"uniqueId":     "MDU6TGFiZWwy-MDU6SXNzdWUz",
					},
					{
						"id":           "MDU6SXNzdWU0",
						"title":        "issue2",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"labelId":      "MDU6TGFiZWwy",
						"uniqueId":     "MDU6TGFiZWwy-MDU6SXNzdWU0",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// There are no labels in repo 2/2 for org 1/2, so there are no issue labels.
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// There are no repositories in org 2, so there are no issues or issue labels.
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetPullRequestLabelPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// PullRequestLabels Page 1: Org 1/2, Repo 1/2, Label [1]/2, PullRequest [1]/1
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestLabelEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDExOlB1bGxSZXF1ZXN0MQ==",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"labelId":      "MDU6TGFiZWw0",
						"title":        "Create README.md",
						"uniqueId":     "MDU6TGFiZWw0-MDExOlB1bGxSZXF1ZXN0MQ==",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("NQ")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestLabels Page 2: Org 1/2, Repo 1/2, Label [2]/2, (has no pull requests)
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("NQ")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestLabels Page 3: Org 1/2, Repo 2/2, Label [1]/1, PullRequest [1, 2]/2
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDExOlB1bGxSZXF1ZXsdsd$S0=",
						"title":        "BRANCH4PR",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"labelId":      "MDU6TGFiZWw1",
						"uniqueId":     "MDU6TGFiZWw1-MDExOlB1bGxSZXF1ZXsdsd$S0=",
					},
					{
						"id":           "MDExOlB1bGxSZXFsssd@@",
						"title":        "BRANCH5PR UPDATE README",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"labelId":      "MDU6TGFiZWw1",
						"uniqueId":     "MDU6TGFiZWw1-MDExOlB1bGxSZXFsssd@@",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestLabels Page 4: Org 2/2, (has no repos)
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestLabel",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestLabelEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetPullRequestLabelPageWithOrganizations(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// PullRequestLabels Page 1: Org 1/2, Repo 1/2, Label [1]/2, PullRequest [1]/1
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestLabel",
				Organizations:         []string{"arvindorg1", "arvindorg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 300,
				EntityConfig:          PopulateDefaultPullRequestLabelEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "MDExOlB1bGxSZXF1ZXN0MQ==",
						"labelId":  "MDU6TGFiZWw0",
						"title":    "Create README.md",
						"uniqueId": "MDU6TGFiZWw0-MDExOlB1bGxSZXF1ZXN0MQ==",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6bnVsbCwib3JnYW5pemF0aW9uT2Zmc2V0IjowLCJJbm5lclBhZ2VJbmZvIjp7Imhhc05leHRQYWdlIjpmYWxzZSwiZW5kQ3Vyc29yIjoiTlEiLCJvcmdhbml6YXRpb25PZmZzZXQiOjAsIklubmVyUGFnZUluZm8iOm51bGx9fQ=="),
				},
			},
			wantErr: nil,
		},
		// PullRequestLabels Page 2: Org 1/2, Repo 1/2, Label [2]/2, (has no pull requests)
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestLabel",
				Organizations:         []string{"arvindorg1", "arvindorg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 300,
				EntityConfig:          PopulateDefaultPullRequestLabelEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6bnVsbCwib3JnYW5pemF0aW9uT2Zmc2V0IjowLCJJbm5lclBhZ2VJbmZvIjp7Imhhc05leHRQYWdlIjpmYWxzZSwiZW5kQ3Vyc29yIjoiTlEiLCJvcmdhbml6YXRpb25PZmZzZXQiOjAsIklubmVyUGFnZUluZm8iOm51bGx9fQ=="),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6IlkzVnljMjl5T25ZeU9wRUIiLCJvcmdhbml6YXRpb25PZmZzZXQiOjAsIklubmVyUGFnZUluZm8iOm51bGx9"),
				},
			},
			wantErr: nil,
		},
		// PullRequestLabels Page 3: Org 1/2, Repo 2/2, Label [1]/1, PullRequest [1, 2]/2
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestLabel",
				Organizations:         []string{"arvindorg1", "arvindorg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 300,
				EntityConfig:          PopulateDefaultPullRequestLabelEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6IlkzVnljMjl5T25ZeU9wRUIiLCJvcmdhbml6YXRpb25PZmZzZXQiOjAsIklubmVyUGFnZUluZm8iOm51bGx9"),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":       "MDExOlB1bGxSZXF1ZXsdsd$S0=",
						"title":    "BRANCH4PR",
						"labelId":  "MDU6TGFiZWw1",
						"uniqueId": "MDU6TGFiZWw1-MDExOlB1bGxSZXF1ZXsdsd$S0=",
					},
					{
						"id":       "MDExOlB1bGxSZXFsssd@@",
						"title":    "BRANCH5PR UPDATE README",
						"labelId":  "MDU6TGFiZWw1",
						"uniqueId": "MDU6TGFiZWw1-MDExOlB1bGxSZXFsssd@@",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6bnVsbCwib3JnYW5pemF0aW9uT2Zmc2V0IjoxLCJJbm5lclBhZ2VJbmZvIjpudWxsfQ=="),
				},
			},
			wantErr: nil,
		},
		// PullRequestLabels Page 4: Org 2/2, (has no repos)
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestLabel",
				Organizations:         []string{"arvindorg1", "arvindorg2"},
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestLabelEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6bnVsbCwib3JnYW5pemF0aW9uT2Zmc2V0IjoxLCJJbm5lclBhZ2VJbmZvIjpudWxsfQ=="),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetIssuePage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// These are issues 2/2 for repo 1/2 for org 1/2.
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              8,
				EntityExternalID:      "Issue",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDU6SXNzdWUz",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"title":        "issue1",
						"author": map[string]any{
							"login": "arooxa",
						},
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"createdAt":    "2024-03-15T18:40:52Z",
						"isPinned":     false,
					},
					{
						"id":           "MDU6SXNzdWU0",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"title":        "issue2",
						"author": map[string]any{
							"login": "arooxa",
						},
						"repositoryId": "MDEwOlJlcG9zaXRvcnkx",
						"createdAt":    "2024-03-15T18:41:04Z",
						"isPinned":     false,
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// These are issues 2/2 for repo 2/2 for org 1/2.
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              8,
				EntityExternalID:      "Issue",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDU6SXNzdWUy",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"title":        "issue3",
						"author": map[string]any{
							"login": "arooxa",
						},
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"createdAt":    "2024-03-14T17:43:03Z",
						"isPinned":     false,
					},
					{
						"id":           "MDU6SXNzdWU1",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"title":        "issue4",
						"author": map[string]any{
							"login": "arooxa",
						},
						"repositoryId": "MDEwOlJlcG9zaXRvcnky",
						"createdAt":    "2024-03-15T18:42:01Z",
						"isPinned":     false,
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// There are no repositories in org 2, so there are no issues.
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              8,
				EntityExternalID:      "Issue",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetIssueAssigneePage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// These are the 2 assignees on issue 1/2 for repo 1/2 for org 1/2.
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueAssigneeEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"issueId":      "MDU6SXNzdWUz",
						"uniqueId":     "MDU6SXNzdWUz-MDQ6VXNlcjQ=",
						"login":        "arooxa",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
					{
						"id":           "MDQ6VXNlcjk=",
						"issueId":      "MDU6SXNzdWUz",
						"uniqueId":     "MDU6SXNzdWUz-MDQ6VXNlcjk=",
						"login":        "isabella-sgnl",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// This is the 1 assignee on issue 2/2 for repo 1/2 for org 1/2.
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueAssigneeEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjk=",
						"issueId":      "MDU6SXNzdWU0",
						"uniqueId":     "MDU6SXNzdWU0-MDQ6VXNlcjk=",
						"login":        "isabella-sgnl",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// These are the 2 assignees on issue 1/2 for repo 2/2 for org 1/2.
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueAssigneeEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"issueId":      "MDU6SXNzdWUy",
						"uniqueId":     "MDU6SXNzdWUy-MDQ6VXNlcjQ=",
						"login":        "arooxa",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
					{
						"id":           "MDQ6VXNlcjk=",
						"issueId":      "MDU6SXNzdWUy",
						"uniqueId":     "MDU6SXNzdWUy-MDQ6VXNlcjk=",
						"login":        "isabella-sgnl",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// These are the 2 assignees on issue 2/2 for repo 2/2 for org 1/2.
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueAssigneeEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"issueId":      "MDU6SXNzdWU1",
						"uniqueId":     "MDU6SXNzdWU1-MDQ6VXNlcjQ=",
						"login":        "arooxa",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
					{
						"id":           "MDQ6VXNlcjEw",
						"issueId":      "MDU6SXNzdWU1",
						"uniqueId":     "MDU6SXNzdWU1-MDQ6VXNlcjEw",
						"login":        "r-rakshith",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// There are no repositories in org 2, so there are no issues or issue assignees.
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueAssigneeEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetIssueParticipantPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// This is the only participant on issue 1/2 for repo 1/2 for org 1/2.
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueParticipantEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"issueId":      "MDU6SXNzdWUz",
						"uniqueId":     "MDU6SXNzdWUz-MDQ6VXNlcjQ=",
						"login":        "arooxa",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// This is the only participant on issue 2/2 for repo 1/2 for org 1/2.
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueParticipantEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"issueId":      "MDU6SXNzdWU0",
						"uniqueId":     "MDU6SXNzdWU0-MDQ6VXNlcjQ=",
						"login":        "arooxa",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// These are the 2 participants on issue 1/2 for repo 2/2 for org 1/2.
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueParticipantEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"issueId":      "MDU6SXNzdWUy",
						"uniqueId":     "MDU6SXNzdWUy-MDQ6VXNlcjQ=",
						"login":        "arooxa",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
					{
						"id":           "MDQ6VXNlcjEw",
						"issueId":      "MDU6SXNzdWUy",
						"uniqueId":     "MDU6SXNzdWUy-MDQ6VXNlcjEw",
						"login":        "r-rakshith",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// These are the 2 participants on issue 2/2 for repo 2/2 for org 1/2.
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueParticipantEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDQ6VXNlcjQ=",
						"issueId":      "MDU6SXNzdWU1",
						"uniqueId":     "MDU6SXNzdWU1-MDQ6VXNlcjQ=",
						"login":        "arooxa",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
					{
						"id":           "MDQ6VXNlcjEw",
						"issueId":      "MDU6SXNzdWU1",
						"uniqueId":     "MDU6SXNzdWU1-MDQ6VXNlcjEw",
						"login":        "r-rakshith",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// There are no repositories in org 2, so there are no issues or issue participants.
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "IssueParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultIssueParticipantEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetPullRequestPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// PullRequests Page 1: Org 1/2, Repo 1/2, PR 1/1
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequest",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDExOlB1bGxSZXF1ZXN0MQ==",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"title":        "Create README.md",
						"closed":       false,
						"createdAt":    "2024-03-13T23:07:49Z",
						"author": map[string]any{
							"login": "arooxa",
						},
						"baseRepository": map[string]any{
							"id": "MDEwOlJlcG9zaXRvcnkx",
						},
						"headRepository": map[string]any{
							"id": "MDEwOlJlcG9zaXRvcnkx",
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequests Page 2: Org 1/2, Repo 2/2, PR [1, 2]/3
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequest",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDExOlB1bGxSZXF1ZXN0Mg==",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"title":        "[branch4PR] README",
						"closed":       false,
						"createdAt":    "2024-03-15T18:43:27Z",
						"author": map[string]any{
							"login": "arooxa",
						},
						"baseRepository": map[string]any{
							"id": "MDEwOlJlcG9zaXRvcnky",
						},
						"headRepository": map[string]any{
							"id": "MDEwOlJlcG9zaXRvcnky",
						},
					},
					{
						"id":           "MDExOlB1bGxSZXF1ZXN0Mw==",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"title":        "[branch5PR] README.md",
						"closed":       false,
						"createdAt":    "2024-03-15T18:46:54Z",
						"author": map[string]any{
							"login": "arooxa",
						},
						"baseRepository": map[string]any{
							"id": "MDEwOlJlcG9zaXRvcnky",
						},
						"headRepository": map[string]any{
							"id": "MDEwOlJlcG9zaXRvcnky",
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequests Page 3: Org 1/2, Repo 2/2, PR [3]/3
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequest",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":           "MDExOlB1bGxSZXF1ZXN0NA==",
						"enterpriseId": "MDEwOkVudGVycHJpc2Ux",
						"title":        "[branch6PR] readMe",
						"closed":       false,
						"createdAt":    "2024-03-15T22:40:43Z",
						"author": map[string]any{
							"login": "arooxa",
						},
						"baseRepository": map[string]any{
							"id": "MDEwOlJlcG9zaXRvcnky",
						},
						"headRepository": map[string]any{
							"id": "MDEwOlJlcG9zaXRvcnky",
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequests Page 4: Org 2/2 (has no repos)
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequest",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetPullRequestChangedFilePage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// PullRequestChangedFiles Page 1: Org 1/2, Repo 1/2, PR 1/1, Files [1]/1
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequestChangedFile",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestChangedFileEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"path":          "README.md",
						"changeType":    "ADDED",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0MQ==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0MQ==-README.md",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestChangedFiles Page 2: Org 1/2, Repo 2/2, PR 1/3, Files [1]/1
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequestChangedFile",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestChangedFileEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"path":          "random/file.txt",
						"changeType":    "DELETED",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-random/file.txt",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestChangedFiles Page 3: Org 1/2, Repo 2/2, PR 2/3, Files [1]/1
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequestChangedFile",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestChangedFileEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"path":          "README.md",
						"changeType":    "ADDED",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mw==-README.md",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestChangedFiles Page 4: Org 1/2, Repo 2/2, PR 3/3, Files [1]/1
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequestChangedFile",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestChangedFileEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"path":          "random/file.txt",
						"changeType":    "ADDED",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0NA==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0NA==-random/file.txt",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestChangedFiles Page 5: Org 2/2 (has no repos)
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              2,
				EntityExternalID:      "PullRequestChangedFile",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestChangedFileEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetPullRequestAssigneePage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// PullRequestAssignee Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Assignees [1]/1
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestAssigneeEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDQ6VXNlcjQ=",
						"login":         "arooxa",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0MQ==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0MQ==-MDQ6VXNlcjQ=",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestAssignee Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Assignees [1, 2]/2
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestAssigneeEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDQ6VXNlcjQ=",
						"login":         "arooxa",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-MDQ6VXNlcjQ=",
					},
					{
						"id":            "MDQ6VXNlcjk=",
						"login":         "isabella-sgnl",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-MDQ6VXNlcjk=",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestAssignee Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Assignees [1, 1]/1
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestAssigneeEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDQ6VXNlcjQ=",
						"login":         "arooxa",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mw==-MDQ6VXNlcjQ=",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestAssignee Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, (has no assignees)
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestAssigneeEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]interface{}{},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestAssignee Page 5: Org 2/2 (has no repos)
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestAssignee",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestAssigneeEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]interface{}{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetPullRequestParticipantPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// PullRequestParticipant Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Participants [1]/1
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestParticipantEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDQ6VXNlcjQ=",
						"login":         "arooxa",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0MQ==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0MQ==-MDQ6VXNlcjQ=",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestParticipant Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Participants [1, 2]/2
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestParticipantEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDQ6VXNlcjQ=",
						"login":         "arooxa",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-MDQ6VXNlcjQ=",
					},
					{
						"id":            "MDQ6VXNlcjEw",
						"login":         "r-rakshith",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mg==-MDQ6VXNlcjEw",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestParticipant Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Participants [1, 2]/2
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestParticipantEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDQ6VXNlcjQ=",
						"login":         "arooxa",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mw==-MDQ6VXNlcjQ=",
					},
					{
						"id":            "MDQ6VXNlcjEw",
						"login":         "r-rakshith",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0Mw==-MDQ6VXNlcjEw",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestParticipant Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Participants [1]/1
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestParticipantEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDQ6VXNlcjQ=",
						"login":         "arooxa",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0NA==",
						"uniqueId":      "MDExOlB1bGxSZXF1ZXN0NA==-MDQ6VXNlcjQ=",
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestParticipant Page 5: Org 2/2 (has no repos)
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestParticipant",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestParticipantEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]interface{}{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetPullRequestCommitPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// PullRequestCommit Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Commits [1]/1
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestCommit",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestCommitEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MTo0YWNkMDEzNTJkNTZjYTMzMTA1ZmMyMjU4ZDFmMTI4NzZmMzhlZjRh",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0MQ==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"commit": map[string]interface{}{
							"id":            "MDY6Q29tbWl0MTo0YWNkMDEzNTJkNTZjYTMzMTA1ZmMyMjU4ZDFmMTI4NzZmMzhlZjRh",
							"committedDate": "2024-03-13T23:07:39Z",
							"author": map[string]interface{}{
								"email": "arvind@sgnl.ai",
								"user": map[string]interface{}{
									"id":    "MDQ6VXNlcjQ=",
									"login": "arooxa",
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestCommit Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Commits [1, 3]/3
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestCommit",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestCommitEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MjozZjBiMmRiMDM3NmJjYTgwNjM0NDRmNjI4ZWI3ZWI5Y2U4NTk1ZGNj",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"commit": map[string]interface{}{
							"id":            "MDY6Q29tbWl0MjozZjBiMmRiMDM3NmJjYTgwNjM0NDRmNjI4ZWI3ZWI5Y2U4NTk1ZGNj",
							"committedDate": "2024-03-15T18:43:10Z",
							"author": map[string]interface{}{
								"email": "arvind@sgnl.ai",
								"user": map[string]interface{}{
									"id":    "MDQ6VXNlcjQ=",
									"login": "arooxa",
								},
							},
						},
					},
					{
						"id":            "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0Mjo2MTFlOTU3NGUzODNiNWQ2NmVjNjAwNDMxYTg4ODRkMzc4OGJiMTQx",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"commit": map[string]interface{}{
							"id":            "MDY6Q29tbWl0Mjo2MTFlOTU3NGUzODNiNWQ2NmVjNjAwNDMxYTg4ODRkMzc4OGJiMTQx",
							"committedDate": "2024-03-16T21:18:12Z",
							"author": map[string]interface{}{
								"email": "arvind@sgnl.ai",
								"user": map[string]interface{}{
									"id":    "MDQ6VXNlcjQ=",
									"login": "arooxa",
								},
							},
						},
					},
					{
						"id":            "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MjpjMWMzNmQ2ZWQ0M2U4ZmVmMjlhNGExNTc2ZWQxZTYxNGZkMGMzNDFi",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mg==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"commit": map[string]interface{}{
							"id":            "MDY6Q29tbWl0MjpjMWMzNmQ2ZWQ0M2U4ZmVmMjlhNGExNTc2ZWQxZTYxNGZkMGMzNDFi",
							"committedDate": "2024-03-22T21:48:21Z",
							"author": map[string]interface{}{
								"email": "rakshith@sgnl.ai",
								"user": map[string]interface{}{
									"id":    "MDQ6VXNlcjEw",
									"login": "r-rakshith",
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestCommit Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Commits [1]/1
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestCommit",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestCommitEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0Mzo1OTYxZGE3NDk1NmJhNWRiYTQ0YWEyYjQ4Mjc2MzM4MGNkNDhhMWZj",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0Mw==",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"commit": map[string]interface{}{
							"id":            "MDY6Q29tbWl0Mjo1OTYxZGE3NDk1NmJhNWRiYTQ0YWEyYjQ4Mjc2MzM4MGNkNDhhMWZj",
							"committedDate": "2024-03-15T18:45:03Z",
							"author": map[string]interface{}{
								"email": "arvind@sgnl.ai",
								"user": map[string]interface{}{
									"id":    "MDQ6VXNlcjQ=",
									"login": "arooxa",
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestCommit Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Commits [1, 2]/2
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestCommit",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestCommitEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"id":            "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0NDo1YTVlNzJmNWQwZjk0MjVlZTk3NDc4NzMxZTc2MDczYjBmMTYzY2Fi",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0NA==",
						"commit": map[string]interface{}{
							"id":            "MDY6Q29tbWl0Mjo1YTVlNzJmNWQwZjk0MjVlZTk3NDc4NzMxZTc2MDczYjBmMTYzY2Fi",
							"committedDate": "2024-03-15T22:39:33Z",
							"author": map[string]interface{}{
								"email": "arvind@sgnl.ai",
								"user": map[string]interface{}{
									"id":    "MDQ6VXNlcjQ=",
									"login": "arooxa",
								},
							},
						},
					},
					{
						"id":            "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0NDpkNjE2NmYwYTlmMmQwMGZlYmFjYzZhYTM3MTAwYWY0YzAxNzBlYzhk",
						"enterpriseId":  "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId": "MDExOlB1bGxSZXF1ZXN0NA==",
						"commit": map[string]interface{}{
							"id":            "MDY6Q29tbWl0MjpkNjE2NmYwYTlmMmQwMGZlYmFjYzZhYTM3MTAwYWY0YzAxNzBlYzhk",
							"committedDate": "2024-03-15T22:44:24Z",
							"author": map[string]interface{}{
								"email": "arvind@sgnl.ai",
								"user": map[string]interface{}{
									"id":    "MDQ6VXNlcjQ=",
									"login": "arooxa",
								},
							},
						},
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestCommit Page 5: Org 2/2 (has no repos)
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestCommit",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestCommitEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]interface{}{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetPullRequestReviewPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		// PullRequestReview Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, (has no reviews)
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestReview",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestReviewEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]interface{}{},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestReview Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Reviews [1]/1
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestReview",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestReviewEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"author": map[string]interface{}{
							"login": "r-rakshith",
						},
						"pullRequestId":             "MDExOlB1bGxSZXF1ZXN0Mg==",
						"enterpriseId":              "MDEwOkVudGVycHJpc2Ux",
						"state":                     "APPROVED",
						"id":                        "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3NQ==",
						"createdAt":                 "2024-03-15T21:05:52Z",
						"authorCanPushToRepository": true,
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestReview Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Reviews [1, 2]/2
		"third_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestReview",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestReviewEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"author": map[string]interface{}{
							"login": "r-rakshith",
						},
						"pullRequestId":             "MDExOlB1bGxSZXF1ZXN0Mw==",
						"enterpriseId":              "MDEwOkVudGVycHJpc2Ux",
						"state":                     "APPROVED",
						"id":                        "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Ng==",
						"createdAt":                 "2024-03-15T21:06:25Z",
						"authorCanPushToRepository": true,
					},
					{
						"author": map[string]interface{}{
							"login": "isabella-sgnl",
						},
						"pullRequestId":             "MDExOlB1bGxSZXF1ZXN0Mw==",
						"enterpriseId":              "MDEwOkVudGVycHJpc2Ux",
						"state":                     "APPROVED",
						"id":                        "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Mg==",
						"createdAt":                 "2024-03-15T19:45:09Z",
						"authorCanPushToRepository": true,
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestReview Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Reviews [1]/1
		"fourth_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestReview",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestReviewEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"author": map[string]interface{}{
							"login": "isabella-sgnl",
						},
						"enterpriseId":              "MDEwOkVudGVycHJpc2Ux",
						"pullRequestId":             "MDExOlB1bGxSZXF1ZXN0NA==",
						"state":                     "CHANGES_REQUESTED",
						"id":                        "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3OA==",
						"createdAt":                 "2024-03-15T22:46:20Z",
						"authorCanPushToRepository": true,
					},
				},
				NextCursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantErr: nil,
		},
		// PullRequestReview Page 5: Org 2/2 (has no repos)
		"last_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              5,
				EntityExternalID:      "PullRequestReview",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultPullRequestReviewEntityConfig(),
				Cursor: CreateGraphQLCompositeCursor(
					[]*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")},
					nil,
					nil,
				),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]interface{}{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetSecretScanningAlertPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "SecretScanningAlert",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultSecretScanningAlertEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"number":        float64(2),
						"created_at":    "2020-11-06T18:48:51Z",
						"url":           "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/2",
						"html_url":      "https://github.com/owner/private-repo/security/secret-scanning/2",
						"locations_url": "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/2/locations",
						"state":         "resolved",
						"resolution":    "false_positive",
						"resolved_at":   "2020-11-07T02:47:13Z",
						"resolved_by": map[string]interface{}{
							"login":               "monalisa",
							"id":                  float64(2),
							"node_id":             "MDQ6VXNlcjI=",
							"avatar_url":          "https://alambic.github.com/avatars/u/2?",
							"gravatar_id":         "",
							"url":                 "https://api.github.com/users/monalisa",
							"html_url":            "https://github.com/monalisa",
							"followers_url":       "https://api.github.com/users/monalisa/followers",
							"following_url":       "https://api.github.com/users/monalisa/following{/other_user}",
							"gists_url":           "https://api.github.com/users/monalisa/gists{/gist_id}",
							"starred_url":         "https://api.github.com/users/monalisa/starred{/owner}{/repo}",
							"subscriptions_url":   "https://api.github.com/users/monalisa/subscriptions",
							"organizations_url":   "https://api.github.com/users/monalisa/orgs",
							"repos_url":           "https://api.github.com/users/monalisa/repos",
							"events_url":          "https://api.github.com/users/monalisa/events{/privacy}",
							"received_events_url": "https://api.github.com/users/monalisa/received_events",
							"type":                "User",
							"site_admin":          true,
						},
						"secret_type":              "adafruit_io_key",
						"secret_type_display_name": "Adafruit IO Key",
						"secret":                   "aio_XXXXXXXXXXXXXXXXXXXXXXXXXXXX",
						"repository": map[string]interface{}{
							"id":        float64(1296269),
							"node_id":   "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
							"name":      "Hello-World",
							"full_name": "octocat/Hello-World",
							"owner": map[string]interface{}{
								"login":               "octocat",
								"id":                  float64(1),
								"node_id":             "MDQ6VXNlcjE=",
								"avatar_url":          "https://github.com/images/error/octocat_happy.gif",
								"gravatar_id":         "",
								"url":                 "https://api.github.com/users/octocat",
								"html_url":            "https://github.com/octocat",
								"followers_url":       "https://api.github.com/users/octocat/followers",
								"following_url":       "https://api.github.com/users/octocat/following{/other_user}",
								"gists_url":           "https://api.github.com/users/octocat/gists{/gist_id}",
								"starred_url":         "https://api.github.com/users/octocat/starred{/owner}{/repo}",
								"subscriptions_url":   "https://api.github.com/users/octocat/subscriptions",
								"organizations_url":   "https://api.github.com/users/octocat/orgs",
								"repos_url":           "https://api.github.com/users/octocat/repos",
								"events_url":          "https://api.github.com/users/octocat/events{/privacy}",
								"received_events_url": "https://api.github.com/users/octocat/received_events",
								"type":                "User",
								"site_admin":          false,
							},
							"private":           false,
							"html_url":          "https://github.com/octocat/Hello-World",
							"description":       "This your first repo!",
							"fork":              false,
							"url":               "https://api.github.com/repos/octocat/Hello-World",
							"archive_url":       "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
							"assignees_url":     "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
							"blobs_url":         "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
							"branches_url":      "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
							"collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
							"comments_url":      "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
							"commits_url":       "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
							"compare_url":       "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
							"contents_url":      "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
							"contributors_url":  "https://api.github.com/repos/octocat/Hello-World/contributors",
							"deployments_url":   "https://api.github.com/repos/octocat/Hello-World/deployments",
							"downloads_url":     "https://api.github.com/repos/octocat/Hello-World/downloads",
							"events_url":        "https://api.github.com/repos/octocat/Hello-World/events",
							"forks_url":         "https://api.github.com/repos/octocat/Hello-World/forks",
							"git_commits_url":   "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
							"git_refs_url":      "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
							"git_tags_url":      "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
							"issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
							"issue_events_url":  "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
							"issues_url":        "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
							"keys_url":          "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
							"labels_url":        "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
							"languages_url":     "https://api.github.com/repos/octocat/Hello-World/languages",
							"merges_url":        "https://api.github.com/repos/octocat/Hello-World/merges",
							"milestones_url":    "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
							"notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
							"pulls_url":         "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
							"releases_url":      "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
							"stargazers_url":    "https://api.github.com/repos/octocat/Hello-World/stargazers",
							"statuses_url":      "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
							"subscribers_url":   "https://api.github.com/repos/octocat/Hello-World/subscribers",
							"subscription_url":  "https://api.github.com/repos/octocat/Hello-World/subscription",
							"tags_url":          "https://api.github.com/repos/octocat/Hello-World/tags",
							"teams_url":         "https://api.github.com/repos/octocat/Hello-World/teams",
							"trees_url":         "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
							"hooks_url":         "https://api.github.com/repos/octocat/Hello-World/hooks",
						},
						"push_protection_bypassed_by": map[string]interface{}{
							"login":               "monalisa",
							"id":                  float64(2),
							"node_id":             "MDQ6VXNlcjI=",
							"avatar_url":          "https://alambic.github.com/avatars/u/2?",
							"gravatar_id":         "",
							"url":                 "https://api.github.com/users/monalisa",
							"html_url":            "https://github.com/monalisa",
							"followers_url":       "https://api.github.com/users/monalisa/followers",
							"following_url":       "https://api.github.com/users/monalisa/following{/other_user}",
							"gists_url":           "https://api.github.com/users/monalisa/gists{/gist_id}",
							"starred_url":         "https://api.github.com/users/monalisa/starred{/owner}{/repo}",
							"subscriptions_url":   "https://api.github.com/users/monalisa/subscriptions",
							"organizations_url":   "https://api.github.com/users/monalisa/orgs",
							"repos_url":           "https://api.github.com/users/monalisa/repos",
							"events_url":          "https://api.github.com/users/monalisa/events{/privacy}",
							"received_events_url": "https://api.github.com/users/monalisa/received_events",
							"type":                "User",
							"site_admin":          true,
						},
						"push_protection_bypassed":    true,
						"push_protection_bypassed_at": "2020-11-06T21:48:51Z",
						"resolution_comment":          "Example comment",
						"validity":                    "active",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://test-instance.com/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=2"),
				},
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "SecretScanningAlert",
				EnterpriseSlug:        testutil.GenPtr("SGNL"),
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultSecretScanningAlertEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr(server.URL + "/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=2"),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"number":                   float64(1),
						"created_at":               "2020-11-06T18:18:30Z",
						"url":                      "https://api.github.com/repos/owner/repo/secret-scanning/alerts/1",
						"html_url":                 "https://github.com/owner/repo/security/secret-scanning/1",
						"locations_url":            "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/1/locations",
						"state":                    "open",
						"resolution":               nil,
						"resolved_at":              nil,
						"resolved_by":              nil,
						"secret_type":              "mailchimp_api_key",
						"secret_type_display_name": "Mailchimp API Key",
						"secret":                   "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2",
						"repository": map[string]interface{}{
							"id":        float64(1296269),
							"node_id":   "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
							"name":      "Hello-World",
							"full_name": "octocat/Hello-World",
							"owner": map[string]interface{}{
								"login":               "octocat",
								"id":                  float64(1),
								"node_id":             "MDQ6VXNlcjE=",
								"avatar_url":          "https://github.com/images/error/octocat_happy.gif",
								"gravatar_id":         "",
								"url":                 "https://api.github.com/users/octocat",
								"html_url":            "https://github.com/octocat",
								"followers_url":       "https://api.github.com/users/octocat/followers",
								"following_url":       "https://api.github.com/users/octocat/following{/other_user}",
								"gists_url":           "https://api.github.com/users/octocat/gists{/gist_id}",
								"starred_url":         "https://api.github.com/users/octocat/starred{/owner}{/repo}",
								"subscriptions_url":   "https://api.github.com/users/octocat/subscriptions",
								"organizations_url":   "https://api.github.com/users/octocat/orgs",
								"repos_url":           "https://api.github.com/users/octocat/repos",
								"events_url":          "https://api.github.com/users/octocat/events{/privacy}",
								"received_events_url": "https://api.github.com/users/octocat/received_events",
								"type":                "User",
								"site_admin":          false,
							},
							"private":           false,
							"html_url":          "https://github.com/octocat/Hello-World",
							"description":       "This your first repo!",
							"fork":              false,
							"url":               "https://api.github.com/repos/octocat/Hello-World",
							"archive_url":       "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
							"assignees_url":     "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
							"blobs_url":         "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
							"branches_url":      "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
							"collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
							"comments_url":      "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
							"commits_url":       "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
							"compare_url":       "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
							"contents_url":      "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
							"contributors_url":  "https://api.github.com/repos/octocat/Hello-World/contributors",
							"deployments_url":   "https://api.github.com/repos/octocat/Hello-World/deployments",
							"downloads_url":     "https://api.github.com/repos/octocat/Hello-World/downloads",
							"events_url":        "https://api.github.com/repos/octocat/Hello-World/events",
							"forks_url":         "https://api.github.com/repos/octocat/Hello-World/forks",
							"git_commits_url":   "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
							"git_refs_url":      "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
							"git_tags_url":      "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
							"issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
							"issue_events_url":  "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
							"issues_url":        "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
							"keys_url":          "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
							"labels_url":        "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
							"languages_url":     "https://api.github.com/repos/octocat/Hello-World/languages",
							"merges_url":        "https://api.github.com/repos/octocat/Hello-World/merges",
							"milestones_url":    "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
							"notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
							"pulls_url":         "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
							"releases_url":      "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
							"stargazers_url":    "https://api.github.com/repos/octocat/Hello-World/stargazers",
							"statuses_url":      "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
							"subscribers_url":   "https://api.github.com/repos/octocat/Hello-World/subscribers",
							"subscription_url":  "https://api.github.com/repos/octocat/Hello-World/subscription",
							"tags_url":          "https://api.github.com/repos/octocat/Hello-World/tags",
							"teams_url":         "https://api.github.com/repos/octocat/Hello-World/teams",
							"trees_url":         "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
							"hooks_url":         "https://api.github.com/repos/octocat/Hello-World/hooks",
						},
						"push_protection_bypassed_by": nil,
						"push_protection_bypassed":    false,
						"push_protection_bypassed_at": nil,
						"resolution_comment":          nil,
						"validity":                    "unknown",
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetSecretScanningAlertPage_With_Organizations(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_org_only_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "SecretScanningAlert",
				Organizations:         []string{"org1", "org2"},
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultSecretScanningAlertEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"number":        float64(2),
						"created_at":    "2020-11-06T18:48:51Z",
						"url":           "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/2",
						"html_url":      "https://github.com/owner/private-repo/security/secret-scanning/2",
						"locations_url": "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/2/locations",
						"state":         "resolved",
						"resolution":    "false_positive",
						"resolved_at":   "2020-11-07T02:47:13Z",
						"resolved_by": map[string]interface{}{
							"login":               "monalisa",
							"id":                  float64(2),
							"node_id":             "MDQ6VXNlcjI=",
							"avatar_url":          "https://alambic.github.com/avatars/u/2?",
							"gravatar_id":         "",
							"url":                 "https://api.github.com/users/monalisa",
							"html_url":            "https://github.com/monalisa",
							"followers_url":       "https://api.github.com/users/monalisa/followers",
							"following_url":       "https://api.github.com/users/monalisa/following{/other_user}",
							"gists_url":           "https://api.github.com/users/monalisa/gists{/gist_id}",
							"starred_url":         "https://api.github.com/users/monalisa/starred{/owner}{/repo}",
							"subscriptions_url":   "https://api.github.com/users/monalisa/subscriptions",
							"organizations_url":   "https://api.github.com/users/monalisa/orgs",
							"repos_url":           "https://api.github.com/users/monalisa/repos",
							"events_url":          "https://api.github.com/users/monalisa/events{/privacy}",
							"received_events_url": "https://api.github.com/users/monalisa/received_events",
							"type":                "User",
							"site_admin":          true,
						},
						"secret_type":              "adafruit_io_key",
						"secret_type_display_name": "Adafruit IO Key",
						"secret":                   "aio_XXXXXXXXXXXXXXXXXXXXXXXXXXXX",
						"repository": map[string]interface{}{
							"id":        float64(1296269),
							"node_id":   "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
							"name":      "Hello-World",
							"full_name": "octocat/Hello-World",
							"owner": map[string]interface{}{
								"login":               "octocat",
								"id":                  float64(1),
								"node_id":             "MDQ6VXNlcjE=",
								"avatar_url":          "https://github.com/images/error/octocat_happy.gif",
								"gravatar_id":         "",
								"url":                 "https://api.github.com/users/octocat",
								"html_url":            "https://github.com/octocat",
								"followers_url":       "https://api.github.com/users/octocat/followers",
								"following_url":       "https://api.github.com/users/octocat/following{/other_user}",
								"gists_url":           "https://api.github.com/users/octocat/gists{/gist_id}",
								"starred_url":         "https://api.github.com/users/octocat/starred{/owner}{/repo}",
								"subscriptions_url":   "https://api.github.com/users/octocat/subscriptions",
								"organizations_url":   "https://api.github.com/users/octocat/orgs",
								"repos_url":           "https://api.github.com/users/octocat/repos",
								"events_url":          "https://api.github.com/users/octocat/events{/privacy}",
								"received_events_url": "https://api.github.com/users/octocat/received_events",
								"type":                "User",
								"site_admin":          false,
							},
							"private":           false,
							"html_url":          "https://github.com/octocat/Hello-World",
							"description":       "This your first repo!",
							"fork":              false,
							"url":               "https://api.github.com/repos/octocat/Hello-World",
							"archive_url":       "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
							"assignees_url":     "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
							"blobs_url":         "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
							"branches_url":      "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
							"collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
							"comments_url":      "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
							"commits_url":       "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
							"compare_url":       "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
							"contents_url":      "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
							"contributors_url":  "https://api.github.com/repos/octocat/Hello-World/contributors",
							"deployments_url":   "https://api.github.com/repos/octocat/Hello-World/deployments",
							"downloads_url":     "https://api.github.com/repos/octocat/Hello-World/downloads",
							"events_url":        "https://api.github.com/repos/octocat/Hello-World/events",
							"forks_url":         "https://api.github.com/repos/octocat/Hello-World/forks",
							"git_commits_url":   "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
							"git_refs_url":      "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
							"git_tags_url":      "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
							"issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
							"issue_events_url":  "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
							"issues_url":        "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
							"keys_url":          "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
							"labels_url":        "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
							"languages_url":     "https://api.github.com/repos/octocat/Hello-World/languages",
							"merges_url":        "https://api.github.com/repos/octocat/Hello-World/merges",
							"milestones_url":    "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
							"notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
							"pulls_url":         "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
							"releases_url":      "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
							"stargazers_url":    "https://api.github.com/repos/octocat/Hello-World/stargazers",
							"statuses_url":      "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
							"subscribers_url":   "https://api.github.com/repos/octocat/Hello-World/subscribers",
							"subscription_url":  "https://api.github.com/repos/octocat/Hello-World/subscription",
							"tags_url":          "https://api.github.com/repos/octocat/Hello-World/tags",
							"teams_url":         "https://api.github.com/repos/octocat/Hello-World/teams",
							"trees_url":         "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
							"hooks_url":         "https://api.github.com/repos/octocat/Hello-World/hooks",
						},
						"push_protection_bypassed_by": map[string]interface{}{
							"login":               "monalisa",
							"id":                  float64(2),
							"node_id":             "MDQ6VXNlcjI=",
							"avatar_url":          "https://alambic.github.com/avatars/u/2?",
							"gravatar_id":         "",
							"url":                 "https://api.github.com/users/monalisa",
							"html_url":            "https://github.com/monalisa",
							"followers_url":       "https://api.github.com/users/monalisa/followers",
							"following_url":       "https://api.github.com/users/monalisa/following{/other_user}",
							"gists_url":           "https://api.github.com/users/monalisa/gists{/gist_id}",
							"starred_url":         "https://api.github.com/users/monalisa/starred{/owner}{/repo}",
							"subscriptions_url":   "https://api.github.com/users/monalisa/subscriptions",
							"organizations_url":   "https://api.github.com/users/monalisa/orgs",
							"repos_url":           "https://api.github.com/users/monalisa/repos",
							"events_url":          "https://api.github.com/users/monalisa/events{/privacy}",
							"received_events_url": "https://api.github.com/users/monalisa/received_events",
							"type":                "User",
							"site_admin":          true,
						},
						"push_protection_bypassed":    true,
						"push_protection_bypassed_at": "2020-11-06T21:48:51Z",
						"resolution_comment":          "Example comment",
						"validity":                    "active",
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("1"),
				},
			},
			wantErr: nil,
		},
		"only_page_org2": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "SecretScanningAlert",
				Organizations:         []string{"org1", "org2"},
				IsEnterpriseCloud:     false,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultSecretScanningAlertEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: testutil.GenPtr("1"),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"number":                   float64(1),
						"created_at":               "2020-11-06T18:18:30Z",
						"url":                      "https://api.github.com/repos/owner/repo/secret-scanning/alerts/1",
						"html_url":                 "https://github.com/owner/repo/security/secret-scanning/1",
						"locations_url":            "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/1/locations",
						"state":                    "open",
						"resolution":               nil,
						"resolved_at":              nil,
						"resolved_by":              nil,
						"secret_type":              "mailchimp_api_key",
						"secret_type_display_name": "Mailchimp API Key",
						"secret":                   "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2",
						"repository": map[string]interface{}{
							"id":        float64(1296269),
							"node_id":   "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
							"name":      "Hello-World",
							"full_name": "octocat/Hello-World",
							"owner": map[string]interface{}{
								"login":               "octocat",
								"id":                  float64(1),
								"node_id":             "MDQ6VXNlcjE=",
								"avatar_url":          "https://github.com/images/error/octocat_happy.gif",
								"gravatar_id":         "",
								"url":                 "https://api.github.com/users/octocat",
								"html_url":            "https://github.com/octocat",
								"followers_url":       "https://api.github.com/users/octocat/followers",
								"following_url":       "https://api.github.com/users/octocat/following{/other_user}",
								"gists_url":           "https://api.github.com/users/octocat/gists{/gist_id}",
								"starred_url":         "https://api.github.com/users/octocat/starred{/owner}{/repo}",
								"subscriptions_url":   "https://api.github.com/users/octocat/subscriptions",
								"organizations_url":   "https://api.github.com/users/octocat/orgs",
								"repos_url":           "https://api.github.com/users/octocat/repos",
								"events_url":          "https://api.github.com/users/octocat/events{/privacy}",
								"received_events_url": "https://api.github.com/users/octocat/received_events",
								"type":                "User",
								"site_admin":          false,
							},
							"private":           false,
							"html_url":          "https://github.com/octocat/Hello-World",
							"description":       "This your first repo!",
							"fork":              false,
							"url":               "https://api.github.com/repos/octocat/Hello-World",
							"archive_url":       "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
							"assignees_url":     "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
							"blobs_url":         "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
							"branches_url":      "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
							"collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
							"comments_url":      "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
							"commits_url":       "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
							"compare_url":       "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
							"contents_url":      "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
							"contributors_url":  "https://api.github.com/repos/octocat/Hello-World/contributors",
							"deployments_url":   "https://api.github.com/repos/octocat/Hello-World/deployments",
							"downloads_url":     "https://api.github.com/repos/octocat/Hello-World/downloads",
							"events_url":        "https://api.github.com/repos/octocat/Hello-World/events",
							"forks_url":         "https://api.github.com/repos/octocat/Hello-World/forks",
							"git_commits_url":   "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
							"git_refs_url":      "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
							"git_tags_url":      "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
							"issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
							"issue_events_url":  "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
							"issues_url":        "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
							"keys_url":          "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
							"labels_url":        "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
							"languages_url":     "https://api.github.com/repos/octocat/Hello-World/languages",
							"merges_url":        "https://api.github.com/repos/octocat/Hello-World/merges",
							"milestones_url":    "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
							"notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
							"pulls_url":         "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
							"releases_url":      "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
							"stargazers_url":    "https://api.github.com/repos/octocat/Hello-World/stargazers",
							"statuses_url":      "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
							"subscribers_url":   "https://api.github.com/repos/octocat/Hello-World/subscribers",
							"subscription_url":  "https://api.github.com/repos/octocat/Hello-World/subscription",
							"tags_url":          "https://api.github.com/repos/octocat/Hello-World/tags",
							"teams_url":         "https://api.github.com/repos/octocat/Hello-World/teams",
							"trees_url":         "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
							"hooks_url":         "https://api.github.com/repos/octocat/Hello-World/hooks",
						},
						"push_protection_bypassed_by": nil,
						"push_protection_bypassed":    false,
						"push_protection_bypassed_at": nil,
						"resolution_comment":          nil,
						"validity":                    "unknown",
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(tt.wantRes, gotRes); diff != "" {
				t.Errorf("Response mismatch (-want +got):\n%s", diff)
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

func TestGetOrganizationUserPage_WithOrganizations(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	githubClient := github.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *github.Request
		wantRes *github.Response
		wantErr *framework.Error
	}{
		"first_org_first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "OrganizationUser",
				Organizations:         []string{"ArvindOrg1", "ArvindOrg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 500,
				EntityConfig:          PopulateDefaultOrganizationUserEntityConfig(),
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"role":     "ADMIN",
						"orgId":    "MDEyOk9yZ2FuaXphdGlvbjU=",
						"uniqueId": "MDEyOk9yZ2FuaXphdGlvbjU=-MDQ6VXNlcjQ=",
						"node": map[string]any{
							"id": "MDQ6VXNlcjQ=",
							"organizationVerifiedDomainEmails": []any{
								map[string]any{
									"email": "arvind@sgnldemos.com",
								},
							},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6IlkzVnljMjl5T25ZeU9wRUUiLCJvcmdhbml6YXRpb25PZmZzZXQiOjAsIklubmVyUGFnZUluZm8iOm51bGx9"),
				},
			},
			wantErr: nil,
		},
		"first_org_second_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "OrganizationUser",
				Organizations:         []string{"ArvindOrg1", "ArvindOrg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationUserEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6IlkzVnljMjl5T25ZeU9wRUUiLCJvcmdhbml6YXRpb25PZmZzZXQiOjAsIklubmVyUGFnZUluZm8iOm51bGx9"),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"role":     "MEMBER",
						"orgId":    "MDEyOk9yZ2FuaXphdGlvbjU=",
						"uniqueId": "MDEyOk9yZ2FuaXphdGlvbjU=-MDQ6VXNlcjk=",
						"node": map[string]any{
							"id": "MDQ6VXNlcjk=",
							"organizationVerifiedDomainEmails": []any{
								map[string]any{
									"email": "isabella@sgnldemos.com",
								},
							},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6bnVsbCwib3JnYW5pemF0aW9uT2Zmc2V0IjoxLCJJbm5lclBhZ2VJbmZvIjpudWxsfQ=="),
				},
			},
			wantErr: nil,
		},
		"second_org_first_page": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "OrganizationUser",
				Organizations:         []string{"ArvindOrg1", "ArvindOrg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationUserEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6bnVsbCwib3JnYW5pemF0aW9uT2Zmc2V0IjoxLCJJbm5lclBhZ2VJbmZvIjpudWxsfQ=="),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"role":     "ADMIN",
						"orgId":    "MDEyOk9yZ2FuaXphdGlvbjEy",
						"uniqueId": "MDEyOk9yZ2FuaXphdGlvbjEy-MDQ6VXNlcjQ=",
						"node": map[string]any{
							"id":                               "MDQ6VXNlcjQ=",
							"organizationVerifiedDomainEmails": []any{},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6IlkzVnljMjl5T25ZeU9wRUUiLCJvcmdhbml6YXRpb25PZmZzZXQiOjEsIklubmVyUGFnZUluZm8iOm51bGx9"),
				},
			},
			wantErr: nil,
		},
		"second_org_no_more_users": {
			context: context.Background(),
			request: &github.Request{
				BaseURL:               server.URL,
				Token:                 "Bearer Testtoken",
				PageSize:              1,
				EntityExternalID:      "OrganizationUser",
				Organizations:         []string{"ArvindOrg1", "ArvindOrg2"},
				IsEnterpriseCloud:     true,
				APIVersion:            testutil.GenPtr("v3"),
				RequestTimeoutSeconds: 5,
				EntityConfig:          PopulateDefaultOrganizationUserEntityConfig(),
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("eyJoYXNOZXh0UGFnZSI6ZmFsc2UsImVuZEN1cnNvciI6IlkzVnljMjl5T25ZeU9wRUUiLCJvcmdhbml6YXRpb25PZmZzZXQiOjEsIklubmVyUGFnZUluZm8iOm51bGx9"),
				},
			},
			wantRes: &github.Response{
				StatusCode: http.StatusOK,
				Objects:    []map[string]any{},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := githubClient.GetPage(tt.context, tt.request)

			if diff := cmp.Diff(gotRes.Objects, tt.wantRes.Objects); diff != "" {
				t.Errorf("Differences found: (-got +want)\n%s", diff)
			}

			if !reflect.DeepEqual(gotRes.Objects, tt.wantRes.Objects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotRes.Objects, tt.wantRes.Objects)
			}

			if !ValidateGraphQLCompositeCursor(gotRes.NextCursor, tt.wantRes.NextCursor) {
				t.Errorf("gotNextCursor: %v, wantNextCursor: %v", gotRes.NextCursor, tt.wantRes.NextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
