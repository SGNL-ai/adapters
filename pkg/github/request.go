// Copyright 2026 SGNL.ai, Inc.
package github

import (
	"fmt"
	"net/http"
	"strconv"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

type RequestInfo struct {
	Endpoint           string
	HTTPMethod         string
	Query              string
	OrganizationOffset int
}

type DeploymentInfo struct {
	GraphQLBasePath string
	RESTBasePath    string
	RESTEndpoints   map[string]map[string]string
}

const (
	EnterpriseCloud  string = "EnterpriseCloud"
	EnterpriseServer string = "EnterpriseServer"
)

var (
	EndpointMappings = map[string]DeploymentInfo{
		EnterpriseCloud: {
			GraphQLBasePath: "/graphql",
			RESTBasePath:    "",
			RESTEndpoints: map[string]map[string]string{
				SecretScanningAlert: {
					"enterprise":   "/enterprises/%s/secret-scanning/alerts",
					"organization": "/orgs/%s/secret-scanning/alerts",
				},
			},
		},
		EnterpriseServer: {
			GraphQLBasePath: "/api/graphql",
			RESTBasePath:    "/api/%s", // %s is the APIVersion
			RESTEndpoints: map[string]map[string]string{
				SecretScanningAlert: {
					"enterprise":   "/enterprises/%s/secret-scanning/alerts",
					"organization": "/orgs/%s/secret-scanning/alerts",
				},
			},
		},
	}
)

// PopulateRequestInfo populates the RequestInfo struct with the necessary information to
// make a request to the datasource.
func PopulateRequestInfo(request *Request) (*RequestInfo, *framework.Error) {
	if request == nil {
		return nil, &framework.Error{
			Message: "Request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if _, found := ValidEntityExternalIDs[request.EntityExternalID]; !found {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Invalid entity external ID: %s", request.EntityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	var (
		deploymentInfo     DeploymentInfo
		organizationOffset int
		err                error
	)

	if request.IsEnterpriseCloud {
		deploymentInfo = EndpointMappings[EnterpriseCloud]
	} else {
		deploymentInfo = EndpointMappings[EnterpriseServer]
	}

	// REST Case
	// nolint:nestif
	if ValidEntityExternalIDs[request.EntityExternalID].isRestAPI {
		// The REST cursor is expected to be the endpoint for the next page.
		// If the cursor is not nil, use it as the endpoint.
		if request.Cursor != nil {
			if request.Cursor.CollectionID != nil {
				organizationOffset, err = strconv.Atoi(*request.Cursor.CollectionID)
				if err != nil {
					return nil, &framework.Error{
						Message: "Failed to convert the cursor's collectionID to an integer.",
						Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
					}
				}
			}

			if request.Cursor.Cursor != nil {
				return &RequestInfo{
					Endpoint:           *request.Cursor.Cursor,
					HTTPMethod:         http.MethodGet,
					OrganizationOffset: organizationOffset,
				}, nil
			}
		}

		if request.APIVersion == nil {
			return nil, &framework.Error{
				Message: "APIVersion is not set for an entity that is retrieved through the GitHub REST API.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}

		// Construct the REST endpoint for the first page.
		baseURI := deploymentInfo.RESTBasePath
		if !request.IsEnterpriseCloud {
			baseURI = fmt.Sprintf(baseURI, *request.APIVersion)
		}

		URI := ""
		// Certain endpoints require additional parameters in the URI.
		switch request.EntityExternalID {
		case SecretScanningAlert:
			if request.EnterpriseSlug != nil {
				URI = fmt.Sprintf(deploymentInfo.RESTEndpoints[request.EntityExternalID]["enterprise"], *request.EnterpriseSlug)
			} else if len(request.Organizations) > 0 {
				organizationOffset, frameworkErr := getOrganizationOffsetForRESTAPI(request)
				if frameworkErr != nil {
					return nil, frameworkErr
				}

				URI = fmt.Sprintf(
					deploymentInfo.RESTEndpoints[request.EntityExternalID]["organization"],
					request.Organizations[organizationOffset],
				)
			}
		}

		return &RequestInfo{
			Endpoint: fmt.Sprintf("%s%s%s?per_page=%d",
				request.BaseURL,
				baseURI,
				URI,
				request.PageSize),
			HTTPMethod:         http.MethodGet,
			OrganizationOffset: organizationOffset,
		}, nil
	}

	// GraphQL Case
	query, queryErr := ConstructQuery(request)
	if queryErr != nil {
		return nil, queryErr
	}

	return &RequestInfo{
		Endpoint:   fmt.Sprintf("%s%s", request.BaseURL, deploymentInfo.GraphQLBasePath),
		HTTPMethod: http.MethodPost,
		Query:      query,
	}, nil
}

// getOrganizationOffsetForRESTAPI returns the organization offset for the next page of the REST API.
func getOrganizationOffsetForRESTAPI(request *Request) (int, *framework.Error) {
	organizationOffset := 0

	if request.Cursor != nil && request.Cursor.CollectionID != nil {
		organizationOffset, err := strconv.Atoi(*request.Cursor.CollectionID)
		if err != nil {
			return 0, &framework.Error{
				Message: "Failed to convert the cursor's collectionID to an integer.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		return organizationOffset, nil
	}

	// If the cursor is nil, fetch first page for the first Organization.
	// If the cursor.cursor value is not nil, return the organization offset the from the cursor
	// because there are more pages to fetch.
	return organizationOffset, nil
}
