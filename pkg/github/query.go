// Copyright 2026 SGNL.ai, Inc.

package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	CollectionPageSize = 1
)

// These are a list of attributes that will be added when post-processing the SoR response
// ex. InjectCommonFields() in datasource.go.
// These should be ignored when building the query.
var attributesToIgnore = map[string]struct{}{
	"enterpriseId":  {},
	"issueId":       {},
	"labelId":       {},
	"orgId":         {},
	"pullRequestId": {},
	"repositoryId":  {},
	"uniqueId":      {},
}

// GraphQLPayload is used as a wrapper to construct the query.
type GraphQLPayload struct {
	Query string `json:"query"`
}

// AttributeNode stores the metadata required to build the inner part of the query for an entity.
type AttributeNode struct {
	Name     string
	Children map[string]*AttributeNode
}

// QueryBuilder is an interface that defines the method for building a query.
// Each entity has its own builder struct that contains the query parameters required to retrieve the entity.
type QueryBuilder interface {
	Build(*Request) (string, *framework.Error)
}

// EnterpriseQueryInfo stores the metadata used when making a GitHub enterprise query.
type EnterpriseQueryInfo struct {
	EnterpriseSlug string
	PageSize       int64
}

type OrganizationQueryBuilder struct {
	EnterpriseQueryInfo *EnterpriseQueryInfo
	OrgAfter            *string
	Organizations       []string
	OrganizationOffset  int
}

// OrganizationUser is retrieved through the GitHub Organization query.
type OrganizationUserQueryBuilder struct {
	OrgLogin  string
	PageSize  int64
	UserAfter *string
}

type UserQueryBuilder struct {
	OrganizationQueryBuilder
	UserAfter *string
}

type TeamQueryBuilder struct {
	OrganizationQueryBuilder
	TeamAfter *string
}

type RepositoryQueryBuilder struct {
	OrganizationQueryBuilder
	RepoAfter *string
}

type CollaboratorQueryBuilder struct {
	RepositoryQueryBuilder
	CollabAfter *string
}

type LabelQueryBuilder struct {
	RepositoryQueryBuilder
	LabelAfter *string
}

type IssueLabelQueryBuilder struct {
	LabelQueryBuilder
	IssueAfter *string
}

type PullRequestLabelQueryBuilder struct {
	LabelQueryBuilder
	PullRequestAfter *string
}

type IssueQueryBuilder struct {
	RepositoryQueryBuilder
	IssueAfter *string
}

type IssueAssigneeQueryBuilder struct {
	IssueQueryBuilder
	AssigneeAfter *string
}

type IssueParticipantQueryBuilder struct {
	IssueQueryBuilder
	ParticipantAfter *string
}

type PullRequestQueryBuilder struct {
	RepositoryQueryBuilder
	PullRequestAfter *string
}

type PullRequestChangedFileQueryBuilder struct {
	PullRequestQueryBuilder
	ChangedFileAfter *string
}

type PullRequestAssigneeQueryBuilder struct {
	PullRequestQueryBuilder
	AssigneeAfter *string
}

type PullRequestParticipantQueryBuilder struct {
	PullRequestQueryBuilder
	ParticipantAfter *string
}

type PullRequestCommitQueryBuilder struct {
	PullRequestQueryBuilder
	CommitAfter *string
}

type PullRequestReviewQueryBuilder struct {
	PullRequestQueryBuilder
	ReviewAfter *string
}

func SetAfterParameter(value *string) string {
	if value == nil || *value == "" {
		return ""
	}

	return fmt.Sprintf(", after: \"%s\"", *value)
}

func SetFilterParameter(value *string) string {
	if value == nil || *value == "" {
		return ""
	}

	return fmt.Sprintf(", %s", *value)
}

func SetOrderByParameter(value *string) string {
	if value == nil || *value == "" {
		return ""
	}

	return fmt.Sprintf(", %s", *value)
}

func (b *OrganizationQueryBuilder) Build(request *Request) (string, *framework.Error) {
	if request.EnterpriseSlug != nil {
		orgAfterQuery := SetAfterParameter(b.OrgAfter)

		innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
		if err != nil {
			return "", err
		}

		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					%s
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, b.EnterpriseQueryInfo.PageSize, orgAfterQuery, innerNode.BuildQuery()), nil
	}

	// Return just the attributes of the organization.
	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "")
	if err != nil {
		return "", err
	}

	// If the Organizations slice is populated in the request, the query is simplified to
	// only return details of an organization.
	// This is how the query should be-
	// query {
	// 		organization (login: "org name") {
	// 			announcement,
	// 			announcementExpiresAt,
	// 			announcementUserDismissible,
	// 			anyPinnableItems,
	// 			archivedAt,
	// 			avatarUrl,
	// 			createdAt,
	// 			databaseId,
	// 			....
	// 		}
	// }
	// The `innerNode.BuildQuery()` returns a query that starts with `node  {` and ends with `}`.
	// Trim those off since the query does not require them.
	innerQuery := innerNode.BuildQuery()
	innerQuery = strings.TrimLeft(innerQuery, "node  {")
	innerQuery = strings.TrimRight(innerQuery, "}")

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			%s
		}
	}`, OrganizationName, innerQuery)

	return query, nil
}

func (b *OrganizationUserQueryBuilder) Build(request *Request) (string, *framework.Error) {
	userAfterQuery := SetAfterParameter(b.UserAfter)
	filterQuery := SetFilterParameter(request.Filter)
	orderByQuery := SetOrderByParameter(request.OrderBy)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, &b.OrgLogin, "edges")
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf(`query {
        organization (login: "%s") {
			id
            membersWithRole (first: %d%s%s%s) {
                pageInfo {
					endCursor
                    hasNextPage
                }
				%s
            }
        }
    }`, b.OrgLogin, b.PageSize, userAfterQuery, filterQuery, orderByQuery, innerNode.BuildQuery())

	return query, nil
}

func (b *UserQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	userAfterQuery := SetAfterParameter(b.UserAfter)
	filterQuery := SetFilterParameter(request.Filter)
	orderByQuery := SetOrderByParameter(request.OrderBy)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						membersWithRole (first: %d%s%s%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							%s
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery,
			b.EnterpriseQueryInfo.PageSize, userAfterQuery, filterQuery, orderByQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	return fmt.Sprintf(`query {
		organization (login: "%s") {
				id
				membersWithRole (first: %d%s%s%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					%s
				}
			}
		}`, OrganizationName, b.EnterpriseQueryInfo.PageSize, userAfterQuery, filterQuery, orderByQuery, innerNode.BuildQuery(),
	), nil
}

func (b *TeamQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	teamAfterQuery := SetAfterParameter(b.TeamAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						teams (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							%s
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, b.EnterpriseQueryInfo.PageSize,
			teamAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	return fmt.Sprintf(`query {
		organization (login: "%s") {
				id
				teams (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					%s
				}
			}
		}`, OrganizationName, b.EnterpriseQueryInfo.PageSize, teamAfterQuery, innerNode.BuildQuery(),
	), nil
}

func (b *RepositoryQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	filterQuery := SetFilterParameter(request.Filter)
	orderByQuery := SetOrderByParameter(request.OrderBy)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		query := fmt.Sprintf(`query {
		enterprise (slug: "%s") {
			id
			organizations (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					repositories (first: %d%s%s%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						%s
					}
				}
			}
		}
    }`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, b.EnterpriseQueryInfo.PageSize,
			repoAfterQuery, filterQuery, orderByQuery, innerNode.BuildQuery())

		return query, nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s%s%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				%s
			}
		}
    }`, OrganizationName, b.EnterpriseQueryInfo.PageSize, repoAfterQuery, filterQuery, orderByQuery, innerNode.BuildQuery())

	return query, nil
}

func (b *CollaboratorQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	collabAfterQuery := SetAfterParameter(b.CollabAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								collaborators (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									%s
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, b.EnterpriseQueryInfo.PageSize, collabAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					collaborators (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						%s
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		b.EnterpriseQueryInfo.PageSize, collabAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *LabelQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	labelAfterQuery := SetAfterParameter(b.LabelAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								labels (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									%s
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, b.EnterpriseQueryInfo.PageSize, labelAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					labels (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						%s
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		b.EnterpriseQueryInfo.PageSize, labelAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *IssueQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	issueAfterQuery := SetAfterParameter(b.IssueAfter)
	filterQuery := SetFilterParameter(request.Filter)
	orderByQuery := SetOrderByParameter(request.OrderBy)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								issues (first: %d%s%s%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									%s
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, b.EnterpriseQueryInfo.PageSize, issueAfterQuery, filterQuery, orderByQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					issues (first: %d%s%s%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						%s
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		b.EnterpriseQueryInfo.PageSize, issueAfterQuery, filterQuery, orderByQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *IssueLabelQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	labelAfterQuery := SetAfterParameter(b.LabelAfter)
	issueAfterQuery := SetAfterParameter(b.IssueAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								labels (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										issues (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, labelAfterQuery, b.EnterpriseQueryInfo.PageSize,
			issueAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					labels (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							issues (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		CollectionPageSize, labelAfterQuery,
		b.EnterpriseQueryInfo.PageSize, issueAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *PullRequestLabelQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	labelAfterQuery := SetAfterParameter(b.LabelAfter)
	pullRequestAfterQuery := SetAfterParameter(b.PullRequestAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								labels (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										pullRequests (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, labelAfterQuery, b.EnterpriseQueryInfo.PageSize,
			pullRequestAfterQuery, innerNode.BuildQuery()), nil
	}

	CollectionID := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
        organization (login: "%s") {
            id
            repositories (first: %d%s) {
                pageInfo {
                    endCursor
                    hasNextPage
                }
                nodes {
					id
					labels (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							pullRequests (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
            }
        }
    }`, CollectionID, CollectionPageSize, repoAfterQuery, CollectionPageSize,
		labelAfterQuery, b.EnterpriseQueryInfo.PageSize, pullRequestAfterQuery, innerNode.BuildQuery())

	return query, nil
}

func (b *IssueAssigneeQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	issueAfterQuery := SetAfterParameter(b.IssueAfter)
	assigneeAfterQuery := SetAfterParameter(b.AssigneeAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								issues (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										assignees (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, issueAfterQuery, b.EnterpriseQueryInfo.PageSize,
			assigneeAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					issues (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							assignees (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		CollectionPageSize, issueAfterQuery,
		b.EnterpriseQueryInfo.PageSize, assigneeAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *IssueParticipantQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	issueAfterQuery := SetAfterParameter(b.IssueAfter)
	participantAfterQuery := SetAfterParameter(b.ParticipantAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								issues (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										participants (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, issueAfterQuery, b.EnterpriseQueryInfo.PageSize,
			participantAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					issues (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							participants (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		CollectionPageSize, issueAfterQuery,
		b.EnterpriseQueryInfo.PageSize, participantAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *PullRequestQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	pullRequestAfterQuery := SetAfterParameter(b.PullRequestAfter)
	filterQuery := SetFilterParameter(request.Filter)
	orderByQuery := SetOrderByParameter(request.OrderBy)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								pullRequests (first: %d%s%s%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									%s
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, b.EnterpriseQueryInfo.PageSize, pullRequestAfterQuery, filterQuery, orderByQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					pullRequests (first: %d%s%s%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						%s
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		b.EnterpriseQueryInfo.PageSize, pullRequestAfterQuery, filterQuery, orderByQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *PullRequestChangedFileQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	pullRequestAfterQuery := SetAfterParameter(b.PullRequestAfter)
	changedFileAfterQuery := SetAfterParameter(b.ChangedFileAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								pullRequests (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										files (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, pullRequestAfterQuery, b.EnterpriseQueryInfo.PageSize,
			changedFileAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					pullRequests (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							files (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		CollectionPageSize, pullRequestAfterQuery,
		b.EnterpriseQueryInfo.PageSize, changedFileAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *PullRequestAssigneeQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	pullRequestAfterQuery := SetAfterParameter(b.PullRequestAfter)
	assigneeAfterQuery := SetAfterParameter(b.AssigneeAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								pullRequests (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										assignees (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, pullRequestAfterQuery, b.EnterpriseQueryInfo.PageSize,
			assigneeAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					pullRequests (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							assignees (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		CollectionPageSize, pullRequestAfterQuery,
		b.EnterpriseQueryInfo.PageSize, assigneeAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *PullRequestParticipantQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	pullRequestAfterQuery := SetAfterParameter(b.PullRequestAfter)
	participantAfterQuery := SetAfterParameter(b.ParticipantAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								pullRequests (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										participants (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, pullRequestAfterQuery, b.EnterpriseQueryInfo.PageSize,
			participantAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					pullRequests (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							participants (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		CollectionPageSize, pullRequestAfterQuery,
		b.EnterpriseQueryInfo.PageSize, participantAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *PullRequestCommitQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	pullRequestAfterQuery := SetAfterParameter(b.PullRequestAfter)
	commitAfterQuery := SetAfterParameter(b.CommitAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								pullRequests (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										commits (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, pullRequestAfterQuery, b.EnterpriseQueryInfo.PageSize,
			commitAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					pullRequests (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							commits (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		CollectionPageSize, pullRequestAfterQuery,
		b.EnterpriseQueryInfo.PageSize, commitAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

func (b *PullRequestReviewQueryBuilder) Build(request *Request) (string, *framework.Error) {
	orgAfterQuery := SetAfterParameter(b.OrgAfter)
	repoAfterQuery := SetAfterParameter(b.RepoAfter)
	pullRequestAfterQuery := SetAfterParameter(b.PullRequestAfter)
	reviewAfterQuery := SetAfterParameter(b.ReviewAfter)

	innerNode, err := AttributeQueryBuilder(request.EntityConfig, nil, "nodes")
	if err != nil {
		return "", err
	}

	if request.EnterpriseSlug != nil {
		return fmt.Sprintf(`query {
			enterprise (slug: "%s") {
				id
				organizations (first: %d%s) {
					pageInfo {
						endCursor
						hasNextPage
					}
					nodes {
						id
						repositories (first: %d%s) {
							pageInfo {
								endCursor
								hasNextPage
							}
							nodes {
								id
								pullRequests (first: %d%s) {
									pageInfo {
										endCursor
										hasNextPage
									}
									nodes {
										id
										latestOpinionatedReviews (first: %d%s) {
											pageInfo {
												endCursor
												hasNextPage
											}
											%s
										}
									}
								}
							}
						}
					}
				}
			}
		}`, b.EnterpriseQueryInfo.EnterpriseSlug, CollectionPageSize, orgAfterQuery, CollectionPageSize,
			repoAfterQuery, CollectionPageSize, pullRequestAfterQuery, b.EnterpriseQueryInfo.PageSize,
			reviewAfterQuery, innerNode.BuildQuery()), nil
	}

	OrganizationName := request.Organizations[b.OrganizationOffset]

	query := fmt.Sprintf(`query {
		organization (login: "%s") {
			id
			repositories (first: %d%s) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					id
					pullRequests (first: %d%s) {
						pageInfo {
							endCursor
							hasNextPage
						}
						nodes {
							id
							latestOpinionatedReviews (first: %d%s) {
								pageInfo {
									endCursor
									hasNextPage
								}
								%s
							}
						}
					}
				}
			}
		}
    }`,
		OrganizationName,
		CollectionPageSize, repoAfterQuery,
		CollectionPageSize, pullRequestAfterQuery,
		b.EnterpriseQueryInfo.PageSize, reviewAfterQuery,
		innerNode.BuildQuery())

	return query, nil
}

// AddChild adds a child to the current node and returns the child node.
func (node *AttributeNode) AddChild(path []string) *AttributeNode {
	if node.Children == nil {
		node.Children = make(map[string]*AttributeNode)
	}

	if len(path) == 0 {
		return node
	}

	childName := path[0]
	if _, exists := node.Children[childName]; !exists {
		node.Children[childName] = &AttributeNode{
			Name:     childName,
			Children: make(map[string]*AttributeNode),
		}
	}

	return node.Children[childName].AddChild(path[1:])
}

func AttributeQueryBuilder(
	entityConfig *framework.EntityConfig,
	login *string,
	rootName string,
) (*AttributeNode, *framework.Error) {
	rootParts := GetAttributePath(rootName)
	if len(rootParts) == 0 {
		return nil, &framework.Error{
			Message: "Root name is empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	rootNode := &AttributeNode{
		Name:     rootParts[0],
		Children: make(map[string]*AttributeNode),
	}
	// baseNode and rootNode are different when the rootName is a JSON path. This happens when
	// child entities with JSON path externalIDs are handled.
	baseNode := rootNode.AddChild(rootParts[1:])

	for _, attr := range entityConfig.Attributes {
		if _, found := attributesToIgnore[attr.ExternalId]; found {
			continue // Skip this attribute
		}

		attrParts := GetAttributePath(attr.ExternalId)
		baseNode.AddChild(attrParts)
	}

	for _, child := range entityConfig.ChildEntities {
		if child.ExternalId == OVDE {
			if login == nil {
				return nil, &framework.Error{
					Message: "login is nil for organizationVerifiedDomainEmails attribute.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
				}
			}

			baseNode.AddChild([]string{"node", fmt.Sprintf("organizationVerifiedDomainEmails (login: \"%s\")", *login)})

			continue
		}

		childParts := GetAttributePath(child.ExternalId)
		if len(childParts) == 0 {
			return nil, &framework.Error{
				Message: "Child name is empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}

		childNode, err := AttributeQueryBuilder(child, login, child.ExternalId)
		if err != nil {
			return nil, err
		}

		baseNode.Children[childParts[0]] = childNode
	}

	return rootNode, nil
}

func (node *AttributeNode) BuildQuery() string {
	if len(node.Children) == 0 {
		return node.Name
	}

	childrenQueries := make([]string, 0, len(node.Children))
	keys := make([]string, 0, len(node.Children))

	for key := range node.Children {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		child := node.Children[key]
		childrenQueries = append(childrenQueries, child.BuildQuery())
	}

	return fmt.Sprintf("%s { %s }", node.Name, strings.Join(childrenQueries, ", "))
}

// [sc-22880] TODO: Change this to use extractor package.
// GetAttributePath extracts the attribute name from a simple JSON path string.
// It is designed to work with basic JSON path expressions that start with '$' and are
// followed by direct attribute names (e.g., `$.name`).
// This function does not support complex JSON path queries (like `$..["$ref"]`), and
// it is important to note this limitation when extending or using this adapter.
//
// The function assumes the JSON path is simple and directly references an attribute without
// nested or recursive structures. For example, it can handle `$.name` but not `$..["$ref"]`.
// This limitation stems from the current design, where the function is used to dynamically
// construct GraphQL queries based on the requested attributes.
func GetAttributePath(input string) []string {
	if strings.HasPrefix(input, "$") {
		parts := strings.Split(input, ".")
		if len(parts) > 1 {
			return parts[1:]
		}
	}

	return []string{input}
}

func ConstructQuery(request *Request) (string, *framework.Error) {
	if request == nil {
		return "", &framework.Error{
			Message: "Request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var builder QueryBuilder

	var pageInfo *PageInfo

	if request.Cursor != nil && request.Cursor.Cursor != nil {
		var err *framework.Error

		pageInfo, err = DecodePageInfo(request.Cursor.Cursor)

		if err != nil {
			return "", err
		}
	}

	builder, builderErr := ConstructQueryBuilder(request, pageInfo)
	if builderErr != nil {
		return "", builderErr
	}

	if builder == nil {
		return "", &framework.Error{
			Message: fmt.Sprintf("Unsupported Query for provided entity ID: %s", request.EntityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	query, err := builder.Build(request)
	if err != nil {
		return "", err
	}

	// Marshal the payload to JSON
	jsonData, marshalErr := json.Marshal(GraphQLPayload{Query: query})
	if marshalErr != nil {
		return "", &framework.Error{
			Message: fmt.Sprintf("Failed to marshal query to retrieve Teams: %v.", marshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return string(jsonData), nil
}

func ConstructQueryBuilder(request *Request, pageInfo *PageInfo) (QueryBuilder, *framework.Error) {
	var builder QueryBuilder

	// The 2nd argument for GetPageInfoAfter() will vary depending on the request.
	// If the request is for an enterprise, keep the value of `n` as-is.
	// If the request is for an organization, decrement `n` by 1.
	// If an organization is passed in, the query does not need to worry about fetching the next
	// organization for the enterprise and will not need to populate the `OrgAfter` field.
	// This will reduce the number of `PageInfo` objects that need to be passed around by 1.
	orgListProvided := request.EnterpriseSlug == nil

	switch request.EntityExternalID {
	case OrganizationUser:
		if request.EnterpriseSlug != nil {
			if request.Cursor == nil || request.Cursor.CollectionID == nil {
				return nil, &framework.Error{
					Message: "Cursor or CollectionCursor is nil for OrganizationUser query.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				}
			}

			builder = &OrganizationUserQueryBuilder{
				OrgLogin:  *request.Cursor.CollectionID,
				PageSize:  request.PageSize,
				UserAfter: GetPageInfoAfter(pageInfo, 0, &orgListProvided),
			}
		} else {
			if len(request.Organizations) == 0 {
				return nil, &framework.Error{
					Message: "Organizations is nil or empty for OrganizationUser query.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
				}
			}

			orgLogin := request.Organizations[0]

			if pageInfo != nil {
				orgLogin = request.Organizations[pageInfo.OrganizationOffset]
			}

			builder = &OrganizationUserQueryBuilder{
				OrgLogin:  orgLogin,
				PageSize:  request.PageSize,
				UserAfter: GetPageInfoAfter(pageInfo, 0, &orgListProvided),
			}
		}
	default:
		var orgQueryBuilder OrganizationQueryBuilder
		if len(request.Organizations) > 0 {
			orgQueryBuilder = OrganizationQueryBuilder{
				Organizations: request.Organizations,
				EnterpriseQueryInfo: &EnterpriseQueryInfo{
					PageSize: request.PageSize,
				},
			}

			orgQueryBuilder.OrganizationOffset = 0
			if pageInfo != nil {
				orgQueryBuilder.OrganizationOffset = pageInfo.OrganizationOffset
			}
		} else {
			orgQueryBuilder = OrganizationQueryBuilder{
				EnterpriseQueryInfo: &EnterpriseQueryInfo{
					EnterpriseSlug: *request.EnterpriseSlug,
					PageSize:       request.PageSize,
				},
				OrgAfter: GetPageInfoAfter(pageInfo, 0, &orgListProvided),
			}
		}

		switch request.EntityExternalID {
		case Organization:
			builder = &orgQueryBuilder
		case Team:
			builder = &TeamQueryBuilder{
				OrganizationQueryBuilder: orgQueryBuilder,
				TeamAfter:                GetPageInfoAfter(pageInfo, 1, &orgListProvided),
			}
		case Repository:
			builder = &RepositoryQueryBuilder{
				OrganizationQueryBuilder: orgQueryBuilder,
				RepoAfter:                GetPageInfoAfter(pageInfo, 1, &orgListProvided),
			}
		case User:
			builder = &UserQueryBuilder{
				OrganizationQueryBuilder: orgQueryBuilder,
				UserAfter:                GetPageInfoAfter(pageInfo, 1, &orgListProvided),
			}
		default:
			repoQueryBuilder := RepositoryQueryBuilder{
				OrganizationQueryBuilder: orgQueryBuilder,
				RepoAfter:                GetPageInfoAfter(pageInfo, 1, &orgListProvided),
			}

			switch request.EntityExternalID {
			case Collaborator:
				builder = &CollaboratorQueryBuilder{
					RepositoryQueryBuilder: repoQueryBuilder,
					CollabAfter:            GetPageInfoAfter(pageInfo, 2, &orgListProvided),
				}
			case Label, IssueLabel, PullRequestLabel:
				labelQueryBuilder := LabelQueryBuilder{
					RepositoryQueryBuilder: repoQueryBuilder,
					LabelAfter:             GetPageInfoAfter(pageInfo, 2, &orgListProvided),
				}

				switch request.EntityExternalID {
				case Label:
					builder = &labelQueryBuilder
				case IssueLabel:
					builder = &IssueLabelQueryBuilder{
						LabelQueryBuilder: labelQueryBuilder,
						IssueAfter:        GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				case PullRequestLabel:
					builder = &PullRequestLabelQueryBuilder{
						LabelQueryBuilder: labelQueryBuilder,
						PullRequestAfter:  GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				}
			case Issue, IssueAssignee, IssueParticipant:
				issueQueryBuilder := IssueQueryBuilder{
					RepositoryQueryBuilder: repoQueryBuilder,
					IssueAfter:             GetPageInfoAfter(pageInfo, 2, &orgListProvided),
				}

				switch request.EntityExternalID {
				case Issue:
					builder = &issueQueryBuilder
				case IssueAssignee:
					builder = &IssueAssigneeQueryBuilder{
						IssueQueryBuilder: issueQueryBuilder,
						AssigneeAfter:     GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				case IssueParticipant:
					builder = &IssueParticipantQueryBuilder{
						IssueQueryBuilder: issueQueryBuilder,
						ParticipantAfter:  GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				}
			case PullRequest, PullRequestAssignee, PullRequestParticipant,
				PullRequestCommit, PullRequestChangedFile, PullRequestReview:
				pullRequestQueryBuilder := PullRequestQueryBuilder{
					RepositoryQueryBuilder: repoQueryBuilder,
					PullRequestAfter:       GetPageInfoAfter(pageInfo, 2, &orgListProvided),
				}

				switch request.EntityExternalID {
				case PullRequest:
					builder = &pullRequestQueryBuilder
				case PullRequestAssignee:
					builder = &PullRequestAssigneeQueryBuilder{
						PullRequestQueryBuilder: pullRequestQueryBuilder,
						AssigneeAfter:           GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				case PullRequestParticipant:
					builder = &PullRequestParticipantQueryBuilder{
						PullRequestQueryBuilder: pullRequestQueryBuilder,
						ParticipantAfter:        GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				case PullRequestCommit:
					builder = &PullRequestCommitQueryBuilder{
						PullRequestQueryBuilder: pullRequestQueryBuilder,
						CommitAfter:             GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				case PullRequestChangedFile:
					builder = &PullRequestChangedFileQueryBuilder{
						PullRequestQueryBuilder: pullRequestQueryBuilder,
						ChangedFileAfter:        GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				case PullRequestReview:
					builder = &PullRequestReviewQueryBuilder{
						PullRequestQueryBuilder: pullRequestQueryBuilder,
						ReviewAfter:             GetPageInfoAfter(pageInfo, 3, &orgListProvided),
					}
				}
			}
		}
	}

	return builder, nil
}

func DecodePageInfo(cursor *string) (*PageInfo, *framework.Error) {
	b, err := base64.StdEncoding.DecodeString(*cursor)
	if err != nil {
		return nil, &framework.Error{
			Message: "Cursor.Cursor base64 decoding failed.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var pageInfo PageInfo

	err = json.Unmarshal(b, &pageInfo)
	if err != nil {
		return nil, &framework.Error{
			Message: "PageInfo unmarshalling failed.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return &pageInfo, nil
}

// This returns the PageInfo struct of the n deep layer, where 0 is the outermost layer.
// If n > number of layers or the pageInfo is nil, the function returns nil.
// If n < 0, the function returns the PageInfo of the outermost layer.
// The 3rd argument if the value of `n` is decremented by 1 or not.
// If the request is for an enterprise, keep the value of `n` as-is.
// If the request is for an organization, decrement `n` by 1.
// If an organization is passed in, the query does not need to worry about fetching the next
// organization for the enterprise and will not need to populate the `OrgAfter` field.
// This will reduce the number of `PageInfo` objects that need to be passed around by 1.
func GetPageInfoAfter(pageInfo *PageInfo, n int, orgListProvided *bool) *string {
	if pageInfo == nil {
		return nil
	}

	if orgListProvided != nil && *orgListProvided {
		n--
	}

	if n <= 0 {
		return pageInfo.EndCursor
	}

	return GetPageInfoAfter(pageInfo.InnerPageInfo, n-1, nil)
}

func PopulateOrganizationCollectionConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: Organization,
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
