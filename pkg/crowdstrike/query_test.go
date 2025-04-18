// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst

package crowdstrike

import (
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestEndpointQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		builder  *EndpointQueryBuilder
		request  *Request
		expected string
		wantErr  error
	}{
		"with cursor": {
			builder: &EndpointQueryBuilder{
				Archived: false,
				Enabled:  true,
				PageSize: 100,
			},
			request: &Request{
				GraphQLCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("cursor123"),
				},
			},
			expected: `{
					entities(
						archived: false
						enabled: true
						types: [ENDPOINT]
						sortKey: RISK_SCORE
						sortOrder: DESCENDING
						first: 100
						after: "cursor123"
					)
					{
						pageInfo {
							hasNextPage
							endCursor
						}
						nodes {
							... on EndpointEntity {
								agentId
								agentVersion
								archived
								cid
								creationTime
								earliestSeenTraffic
								entityId
								guestAccountEnabled
								hasADDomainAdminRole
								hasRole
								hostName
								impactScore
								inactive
								lastIpAddress
								learned
								markTime
								mostRecentActivity
								primaryDisplayName
								riskScore
								riskScoreSeverity
								secondaryDisplayName
								shared
								stale
								staticIpAddresses
								type
								unmanaged
								watched
								ztaScore
								accounts {
									__typename
									... on ActiveDirectoryAccountDescriptor {
										archived
										cn
										consistencyGuid
										containingGroupIds
										creationTime
										dataSource
										department
										description
										dn
										domain
										enabled
										expirationTime
										flattenedContainingGroupIds
										lastUpdateTime
										lockoutTime
										mostRecentActivity
										objectGuid
										objectSid
										ou
										samAccountName
										servicePrincipalNames
										title
										upn
										userAccountControl
										userAccountControlFlags
									}
								}
								riskFactors {
									score
									severity
									type
								}
							}
						}
					}
				}`,
		},
		"no cursor": {
			builder: &EndpointQueryBuilder{
				Archived: false,
				Enabled:  true,
				PageSize: 100,
			},
			request: &Request{
				GraphQLCursor: nil,
			},
			expected: `{
					entities(
						archived: false
						enabled: true
						types: [ENDPOINT]
						sortKey: RISK_SCORE
						sortOrder: DESCENDING
						first: 100
					)
					{
						pageInfo {
							hasNextPage
							endCursor
						}
						nodes {
							... on EndpointEntity {
								agentId
								agentVersion
								archived
								cid
								creationTime
								earliestSeenTraffic
								entityId
								guestAccountEnabled
								hasADDomainAdminRole
								hasRole
								hostName
								impactScore
								inactive
								lastIpAddress
								learned
								markTime
								mostRecentActivity
								primaryDisplayName
								riskScore
								riskScoreSeverity
								secondaryDisplayName
								shared
								stale
								staticIpAddresses
								type
								unmanaged
								watched
								ztaScore
								accounts {
									__typename
									... on ActiveDirectoryAccountDescriptor {
										archived
										cn
										consistencyGuid
										containingGroupIds
										creationTime
										dataSource
										department
										description
										dn
										domain
										enabled
										expirationTime
										flattenedContainingGroupIds
										lastUpdateTime
										lockoutTime
										mostRecentActivity
										objectGuid
										objectSid
										ou
										samAccountName
										servicePrincipalNames
										title
										upn
										userAccountControl
										userAccountControlFlags
									}
								}
								riskFactors {
									score
									severity
									type
								}
							}
						}
					}
				}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query, err := tt.builder.Build(tt.request)
			if err != nil {
				if diff := cmp.Diff(tt.wantErr, err); diff != "" {
					t.Fatal(diff)
				}
			}

			if diff := cmp.Diff(NormalizeQuery(tt.expected), NormalizeQuery(query)); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestUserQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		builder  *UserQueryBuilder
		request  *Request
		expected string
		wantErr  error
	}{
		"with cursor": {
			builder: &UserQueryBuilder{
				Archived: false,
				Enabled:  true,
				PageSize: 100,
			},
			request: &Request{
				GraphQLCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("cursor123"),
				},
			},
			expected: `{
				entities(
					archived: false
					enabled: true
					types: [USER]
					sortKey: RISK_SCORE
					sortOrder: DESCENDING
					first: 100
					after: "cursor123"
					) {
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						... on UserEntity {
							archived
							creationTime
							earliestSeenTraffic
							emailAddresses
							entityId
							hasADDomainAdminRole
							impactScore
							inactive
							learned
							markTime
							mostRecentActivity
							riskScore
							riskScoreSeverity
							riskScoreWithoutLinkedAccounts
							secondaryDisplayName
							shared
							stale
							watched
							type
							riskFactors {
								score
								severity
								type
							}
							accounts {
								__typename
								... on ActiveDirectoryAccountDescriptor {
									archived
									cn
									consistencyGuid
									containingGroupIds
									creationTime
									dataSource
									department
									description
									dn
									domain
									enabled
									expirationTime
									flattenedContainingGroupIds
									lastUpdateTime
									lockoutTime
									mostRecentActivity
									objectGuid
									objectSid
									ou
									samAccountName
									servicePrincipalNames
									title
									upn
									userAccountControl
									userAccountControlFlags
								}
							}
							primaryDisplayName
						}
					}
				}
			}`,
		},
		"no cursor": {
			builder: &UserQueryBuilder{
				Archived: false,
				Enabled:  true,
				PageSize: 100,
			},
			request: &Request{
				GraphQLCursor: nil,
			},
			expected: `{
				entities(
					archived: false
					enabled: true
					types: [USER]
					sortKey: RISK_SCORE
					sortOrder: DESCENDING
					first: 100
					) {
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						... on UserEntity {
							archived
							creationTime
							earliestSeenTraffic
							emailAddresses
							entityId
							hasADDomainAdminRole
							impactScore
							inactive
							learned
							markTime
							mostRecentActivity
							riskScore
							riskScoreSeverity
							riskScoreWithoutLinkedAccounts
							secondaryDisplayName
							shared
							stale
							watched
							type
							riskFactors {
								score
								severity
								type
							}
							accounts {
								__typename
								... on ActiveDirectoryAccountDescriptor {
									archived
									cn
									consistencyGuid
									containingGroupIds
									creationTime
									dataSource
									department
									description
									dn
									domain
									enabled
									expirationTime
									flattenedContainingGroupIds
									lastUpdateTime
									lockoutTime
									mostRecentActivity
									objectGuid
									objectSid
									ou
									samAccountName
									servicePrincipalNames
									title
									upn
									userAccountControl
									userAccountControlFlags
								}
							}
							primaryDisplayName
						}
					}
				}
			}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query, err := tt.builder.Build(tt.request)
			if err != nil {
				if diff := cmp.Diff(tt.wantErr, err); diff != "" {
					t.Fatal(diff)
				}
			}

			if diff := cmp.Diff(NormalizeQuery(tt.expected), NormalizeQuery(query)); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestIncidentQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		builder  *IncidentQueryBuilder
		request  *Request
		expected string
		wantErr  error
	}{
		"with cursor": {
			builder: &IncidentQueryBuilder{
				PageSize: 100,
			},
			request: &Request{
				GraphQLCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("cursor123"),
				},
			},
			expected: `{
				incidents(
					first: 100
					sortKey: END_TIME
					sortOrder: DESCENDING
					after: "cursor123"
				)
				{
					pageInfo { hasNextPage endCursor }
					nodes {
						endTime
						incidentId
						lifeCycleStage
						markedAsRead
						severity
						startTime
						type
						compromisedEntities {
							archived
							creationTime
							entityId
							hasADDomainAdminRole
							hasRole
							learned
							markTime
							primaryDisplayName
							riskScore
							riskScoreSeverity
							secondaryDisplayName
							type
							watched
						}
						alertEvents {
							alertId
							alertType
							endTime
							eventId
							eventLabel
							eventSeverity
							eventType
							patternId
							resolved
							startTime
							timestamp
							entities {
								archived
								creationTime
								entityId
								hasADDomainAdminRole
								hasRole
								learned
								markTime
								primaryDisplayName
								riskScore
								riskScoreSeverity
								secondaryDisplayName
								type
								watched
							}
						}
					}
				}
			}`,
		},
		"no cursor": {
			builder: &IncidentQueryBuilder{
				PageSize: 100,
			},
			request: &Request{
				GraphQLCursor: nil,
			},
			expected: `{
				incidents(
					first: 100
					sortKey: END_TIME
					sortOrder: DESCENDING
				)
				{
					pageInfo { hasNextPage endCursor }
					nodes {
						endTime
						incidentId
						lifeCycleStage
						markedAsRead
						severity
						startTime
						type
						compromisedEntities {
							archived
							creationTime
							entityId
							hasADDomainAdminRole
							hasRole
							learned
							markTime
							primaryDisplayName
							riskScore
							riskScoreSeverity
							secondaryDisplayName
							type
							watched
						}
						alertEvents {
							alertId
							alertType
							endTime
							eventId
							eventLabel
							eventSeverity
							eventType
							patternId
							resolved
							startTime
							timestamp
							entities {
								archived
								creationTime
								entityId
								hasADDomainAdminRole
								hasRole
								learned
								markTime
								primaryDisplayName
								riskScore
								riskScoreSeverity
								secondaryDisplayName
								type
								watched
							}
						}
					}
				}
			}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			query, err := tt.builder.Build(tt.request)
			if err != nil {
				if diff := cmp.Diff(tt.wantErr, err); diff != "" {
					t.Fatal(diff)
				}
			}

			if diff := cmp.Diff(NormalizeQuery(tt.expected), NormalizeQuery(query)); diff != "" {
				t.Fatal(diff)
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
