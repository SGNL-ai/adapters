// Copyright 2026 SGNL.ai, Inc.

package crowdstrike

import (
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// QueryBuilder is an interface that defines the method for building a query.
// Each entity has its own builder struct that contains the query parameters required to retrieve the entity.
type QueryBuilder interface {
	Build(*Request) (string, *framework.Error)
}

type UserQueryBuilder struct {
	Archived bool
	Enabled  bool
	PageSize int64
	First    int64
}

type IncidentQueryBuilder struct {
	First    int64
	PageSize int64
}

type EndpointQueryBuilder struct {
	Archived bool
	Enabled  bool
	PageSize int64
	First    int64
}

func SetAfterParameter(value *string) string {
	if value == nil {
		return ""
	}

	return fmt.Sprintf(`after: "%s"`, *value)
}

func (b *UserQueryBuilder) Build(request *Request) (string, *framework.Error) {
	var cursor string
	if request.GraphQLCursor != nil {
		cursor = SetAfterParameter(request.GraphQLCursor.Cursor)
	}

	query := fmt.Sprintf(
		`{
			entities(
				archived: %v
				enabled: %v
				types: [USER]
				sortKey: RISK_SCORE
				sortOrder: DESCENDING
				first: %d
				%s
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
		}`, b.Archived, b.Enabled, b.PageSize, cursor)

	return query, nil
}

func (b *IncidentQueryBuilder) Build(request *Request) (string, *framework.Error) {
	var cursor string
	if request.GraphQLCursor != nil {
		cursor = SetAfterParameter(request.GraphQLCursor.Cursor)
	}

	query := fmt.Sprintf(
		`{
		    incidents(
				first: %d
				sortKey: END_TIME
				sortOrder: DESCENDING
				%s
			) {
				pageInfo {
					hasNextPage
					endCursor
				}
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
		}`, b.PageSize, cursor)

	return query, nil
}

func (b *EndpointQueryBuilder) Build(request *Request) (string, *framework.Error) {
	var cursor string
	if request.GraphQLCursor != nil {
		cursor = SetAfterParameter(request.GraphQLCursor.Cursor)
	}

	query := fmt.Sprintf(
		`{
		    entities(
		        archived: %v
		        enabled: %v
		        types: [ENDPOINT]
		        sortKey: RISK_SCORE
		        sortOrder: DESCENDING
		        first: %d
		    	%s
			) {
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
		b.Archived, b.Enabled, b.PageSize, cursor)

	return query, nil
}

func GetQueryBuilder(request *Request, _ *PageInfo) (QueryBuilder, *framework.Error) {
	var builder QueryBuilder

	switch request.EntityExternalID {
	case User:
		builder = &UserQueryBuilder{
			Archived: request.Config.Archived,
			Enabled:  request.Config.Enabled,
			PageSize: request.PageSize,
		}
	case Incident:
		builder = &IncidentQueryBuilder{
			PageSize: request.PageSize,
		}
	case Endpoint:
		builder = &EndpointQueryBuilder{
			Archived: request.Config.Archived,
			Enabled:  request.Config.Enabled,
			PageSize: request.PageSize,
		}

	default:
		return nil, &framework.Error{
			Message: fmt.Sprintf("Unsupported Query for provided entity ID: %s", request.EntityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	return builder, nil
}
