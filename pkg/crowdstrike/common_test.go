// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst,lll
package crowdstrike_test

import (
	"encoding/json"
	"io"
	"net/http"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/testutil"

	"github.com/sgnl-ai/adapters/pkg/crowdstrike"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

var (
	validAddress         = "api.us-2.crowdstrike.com"
	validAuthCredentials = &framework.DatasourceAuthCredentials{
		HTTPAuthorization: "Bearer testtoken",
	}
	validCommonConfig = &crowdstrike.Config{
		APIVersion: "v1",
	}
)

// GraphQLPayload is used as a wrapper to construct the query.
type GraphQLPayload struct {
	Query     string  `json:"query"`
	Variables *string `json:"variables"`
}

func ValidationQueryBuilder(externalID string, pageSize int64, cursor *string) string {
	req := PopulateDefaultRequest(externalID, pageSize, cursor)

	builder, _ := crowdstrike.GetQueryBuilder(req, nil)

	b, _ := builder.Build(req)

	jsonData, _ := json.Marshal(GraphQLPayload{Query: b, Variables: nil})

	return string(jsonData)
}

func PopulateDefaultRequest(externalID string, pageSize int64, cursor *string) *crowdstrike.Request {
	req := &crowdstrike.Request{
		EntityExternalID: externalID,
		PageSize:         pageSize,
		GraphQLCursor: &pagination.CompositeCursor[string]{
			Cursor: cursor,
		},
		Config: &crowdstrike.Config{
			APIVersion: "v1",
			Archived:   false,
			Enabled:    true,
		},
	}

	switch externalID {
	case crowdstrike.User:
		req.EntityConfig = PopulateUserEntityConfig()
	case crowdstrike.Endpoint:
		req.EntityConfig = PopulateEndpointEntityConfig()
	case crowdstrike.Incident:
		req.EntityConfig = PopulateIncidentEntityConfig()
	default:
		return nil
	}

	return req
}

func PopulateUserEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: crowdstrike.User,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "entityId",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "inactive",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "creationTime",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: `$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "objectGuid",
						Type:       framework.AttributeTypeString,
						List:       false,
						UniqueId:   true,
					},
					{
						ExternalId: "samAccountName",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
		},
	}
}

func PopulateEndpointEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: crowdstrike.Endpoint,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "entityId",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "inactive",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "creationTime",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: `$.accounts[?(@.__typename=="ActiveDirectoryAccountDescriptor")]`,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "objectGuid",
						Type:       framework.AttributeTypeString,
						List:       false,
						UniqueId:   true,
					},
					{
						ExternalId: "samAccountName",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
		},
	}
}

func PopulateIncidentEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: crowdstrike.Incident,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "incidentId",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "severity",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "startTime",
				Type:       framework.AttributeTypeDateTime,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: "$.compromisedEntities",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "entityId",
						Type:       framework.AttributeTypeString,
						List:       false,
						UniqueId:   true,
					},
					{
						ExternalId: "primaryDisplayName",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
		},
	}
}

func PopulateDetectionEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: crowdstrike.Detect,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "detection_id",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "email_sent",
				Type:       framework.AttributeTypeBool,
				List:       false,
			},
			{
				ExternalId: "status",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
	}
}

// Define the endpoints and responses for the mock server.
// This handler is intended to be re-used throughout the test package for GraphQL APIs.
var TestGraphQLServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer Testtoken" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`
		{
			"meta": {
				"query_time": 2.16e-7,
				"powered_by": "crowdstrike-api-gateway",
				"trace_id": "780a7c3b-ae2f-40f0-9cde-318f7faef0f5"
			},
			"errors": [
				{
					"code": 401,
					"message": "access denied, invalid bearer token"
				}
			]
		}`))
	}

	if r.URL.RequestURI() != "/identity-protection/combined/graphql/v1" {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	query, _ := io.ReadAll(r.Body)
	normalizedQueryStr := crowdstrike.NormalizeQuery(string(query))
	// Remove the extra whitespace added at the end of the query
	normalizedQueryStr = normalizedQueryStr[0 : len(normalizedQueryStr)-1]

	// GraphQL Endpoints
	switch normalizedQueryStr {

	// ****************** User Queries ******************
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.User, 2, nil)):
		w.Write([]byte(UserResponsePage1))
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.User, 2, testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNjQ1NDg3MTMzOTk5OTk5OSwiX2lkIjoiNDVkYzQwZTItN2I3Yi00ZjM4LTlhYzctOThmNGEzNWIyNGUxIn0="))):
		w.Write([]byte(UserResponsePage2))
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.User, 2, testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNjQwNDc5MTcxNzM1MjQ4OSwiX2lkIjoiODNhNDllZjEtMTdhNy00ZmE0LWI5MGYtOTE0MmRmYTQ5NTc3In0="))):
		w.Write([]byte(UserResponsePage3))

	// ****************** Endpoint Queries ******************
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.Endpoint, 2, nil)):
		w.Write([]byte(EndpointResponsePage1))
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.Endpoint, 2, testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuNDU5NDAwMDAwMDAwMDAwMDMsIl9pZCI6IjNjN2FlYmI5LTQxMWItNGVlOS1iNDgxLWU4ODFmMjlhZmNjOCJ9"))):
		w.Write([]byte(EndpointResponsePage2))
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.Endpoint, 2, testutil.GenPtr("eyJyaXNrU2NvcmUiOjAuMywiX2lkIjoiZmQxZTBmMGItZjFlMS00MjI0LThkNjAtNGYyOTdhYTkxYzI5In0="))):
		w.Write([]byte(EndpointResponsePage3))

	// ****************** Incident Queries ******************
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.Incident, 2, nil)):
		w.Write([]byte(IncidentResponsePage1))
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.Incident, 2, testutil.GenPtr("eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0yMFQwMTo1NToxMC4yNzRaIn0sInNlcXVlbmNlSWQiOjE1fQ=="))):
		w.Write([]byte(IncidentResponsePage2))
	case crowdstrike.NormalizeQuery(ValidationQueryBuilder(crowdstrike.Incident, 2, testutil.GenPtr("eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0wOVQxNDoyODowNC4wMDhaIn0sInNlcXVlbmNlSWQiOjEzfQ=="))):
		w.Write([]byte(IncidentResponsePage3))
	}
})

// Define the endpoints and responses for the mock server.
// This handler is intended to be re-used throughout the test package for REST APIs.
var TestRESTServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer Testtoken" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`
		{
			"meta": {
				"query_time": 2.16e-7,
				"powered_by": "crowdstrike-api-gateway",
				"trace_id": "780a7c3b-ae2f-40f0-9cde-318f7faef0f5"
			},
			"errors": [
				{
					"code": 401,
					"message": "access denied, invalid bearer token"
				}
			]
		}`))
	}
	switch r.URL.RequestURI() {
	// ************************ List detections ************************
	case "/detects/queries/detects/v1?limit=2":
		w.Write([]byte(DetectListResponseFirstPage))

	case "/detects/queries/detects/v1?limit=2&offset=2":
		w.Write([]byte(DetectListResponseMiddlePage))

	case "/detects/queries/detects/v1?limit=2&offset=4":
		w.Write([]byte(DetectListResponseLastPage))

	// Page 999 mimics a specialized error from CRWD REST APIs
	case "/detects/queries/detects/v1?limit=2&offset=999":
		w.Write([]byte(DetectResponseSpecializedErr))

	// ************************ Detailed detections ************************
	case "/detects/entities/summaries/GET/v1?limit=2":
		w.Write([]byte(DetectDetailedResponseFirstPage))

	case "/detects/entities/summaries/GET/v1?limit=2&offset=2":
		w.Write([]byte(DetectDetailedResponseMiddlePage))

	case "/detects/entities/summaries/GET/v1?limit=2&offset=4":
		w.Write([]byte(DetectDetailedResponseLastPage))

	// Page 999 mimics a specialized error from CRWD REST APIs
	case "/detects/entities/summaries/GET/v1?limit=2&offset=999":
		w.Write([]byte(DetectResponseSpecializedErr))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})
