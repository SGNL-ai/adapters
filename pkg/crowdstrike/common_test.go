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

func PopulateEndpointIncidentEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: crowdstrike.EndpointIncident,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "incident_id",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "state",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "status",
				Type:       framework.AttributeTypeInt64,
				List:       false,
			},
		},
	}
}

func PopulateAlertsEntityConfig() *framework.EntityConfig {
	return &framework.EntityConfig{
		ExternalId: crowdstrike.Alerts,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "composite_id",
				Type:       framework.AttributeTypeString,
				List:       false,
				UniqueId:   true,
			},
			{
				ExternalId: "aggregate_id",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
			{
				ExternalId: "status",
				Type:       framework.AttributeTypeString,
				List:       false,
			},
		},
		ChildEntities: []*framework.EntityConfig{
			{
				ExternalId: "$.files_accessed",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "filename",
						Type:       framework.AttributeTypeString,
						List:       false,
						UniqueId:   false,
					},
					{
						ExternalId: "filepath",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
			{
				ExternalId: "$.files_written",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "filename",
						Type:       framework.AttributeTypeString,
						List:       false,
						UniqueId:   false,
					},
					{
						ExternalId: "filepath",
						Type:       framework.AttributeTypeString,
						List:       false,
					},
				},
			},
			{
				ExternalId: "$.mitre_attack",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "pattern_id",
						Type:       framework.AttributeTypeInt64,
						List:       false,
						UniqueId:   true,
					},
				},
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

	// ************************ Endpoint Incidents ************************
	case "/incidents/queries/incidents/v1?limit=100":
		// Handle GET request for incident IDs - empty list
		if r.Method == http.MethodGet {
			w.Write([]byte(EndpointIncidentEmptyListResponse))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	case "/incidents/entities/incidents/GET/v1?limit=100":
		// Handle POST request for incident details - should error with empty IDs
		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			var reqBody map[string]any
			json.Unmarshal(body, &reqBody)

			// Check if ids array is empty or missing
			ids, hasIDs := reqBody["ids"]
			if !hasIDs || ids == nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(EndpointIncidentEmptyIDsErrorResponse))
			} else if idsArray, ok := ids.([]any); ok && len(idsArray) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(EndpointIncidentEmptyIDsErrorResponse))
			} else {
				// Valid IDs provided - return success
				w.Write([]byte(EndpointIncidentValidResponse))
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	// ************************ Alerts ************************
	case "/alerts/combined/alerts/v1?limit=2":
		// Handle POST request for alerts
		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			var reqBody map[string]any
			json.Unmarshal(body, &reqBody)

			if reqBody["after"] == nil || reqBody["after"] == "" {
				w.Write([]byte(AlertResponseFirstPage))
			} else if reqBody["after"] == "eyJ2ZXJzaW9uIjoidjEiLCJ0b3RhbF9oaXRzIjoyMywidG90YWxfcmVsYXRpb24iOiJlcSIsImNsdXN0ZXJfaWQiOiJ0ZXN0IiwiYWZ0ZXIiOlsxNzQ5NjExMTU3MjIxLCJ0ZXN0aWQ6aW5kOjUzODhjNTkyMTg5NDQ0YWQ5ZTg0ZGYwNzFjOGYzOTU0Ojk3ODI3ODI2MTQtMTAzMDMtMzE4MzE1NjgiXSwidG90YWxfZmV0Y2hlZCI6Mn0=" {
				w.Write([]byte(AlertResponseMiddlePage))
			} else if reqBody["after"] == "eyJ2ZXJzaW9uIjoidjEiLCJ0b3RhbF9oaXRzIjoyMywidG90YWxfcmVsYXRpb24iOiJlcSIsImNsdXN0ZXJfaWQiOiJ0ZXN0IiwiYWZ0ZXIiOlsxNzQ5NTEyMzQ1Njc4LCJ0ZXN0aWQ6aW5kOmU0NTY3ODkwMTIzNDU2Nzg5MGNkZWYxMjM0NTY3ODkwLTIwMTUzLTcwNTEiXSwidG90YWxfZmV0Y2hlZCI6NH0=" {
				w.Write([]byte(AlertResponseLastPage))
			} else if reqBody["after"] == "1000" { // Non existent page - triggers 404
				w.WriteHeader(http.StatusNotFound)
			} else if reqBody["after"] == "999" { // Non existent page - triggers specialized error
				w.Write([]byte(AlertResponseSpecializedErr))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})
