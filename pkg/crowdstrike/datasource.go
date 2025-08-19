// Copyright 2025 SGNL.ai, Inc.
package crowdstrike

// We ingest CrowdStrike data using both GraphQL and REST APIs.
// This file calls respective data fetchers based on the entity.
// See datasource_graphql.go and datasource_rest.go for more details.

import (
	"context"
	"net/http"

	framework "github.com/sgnl-ai/adapter-framework"
)

const (
	User             string = "user"
	Incident         string = "incident"
	Endpoint         string = "endpoint"
	Device           string = "endpoint_protection_device"
	EndpointIncident string = "endpoint_protection_incident"
	Detect           string = "endpoint_protection_detect"
	Alerts           string = "endpoint_protection_alert"
)

// Datasource directly implements a Client interface to allow querying
// an external datasource.
type Datasource struct {
	Client *http.Client
}

// Entity contains entity specific information, such as the entity's unique ID attribute.
type Entity struct {
	// UniqueIDAttrExternalID is the external ID of the entity's uniqueId attribute.
	UniqueIDAttrExternalID string

	// OrderByAttribute is the attribute to order the results by.
	OrderByAttribute string

	// UseIntCursor
	UseIntCursor bool
}

var (
	// ValidGraphQLEntityExternalIDs is a map of valid external IDs of entities that can be queried.
	// The map value is the Entity struct which contains the unique ID attribute.
	ValidGraphQLEntityExternalIDs = map[string]Entity{
		User: {
			UniqueIDAttrExternalID: "entityId",
			OrderByAttribute:       "RISK_SCORE",
		},
		Incident: {
			UniqueIDAttrExternalID: "incidentId",
			OrderByAttribute:       "END_TIME",
		},
		Endpoint: {
			UniqueIDAttrExternalID: "entityId",
			OrderByAttribute:       "RISK_SCORE",
		},
	}

	ValidRESTEntityExternalIDs = map[string]Entity{
		Device:           {},
		Detect:           {UseIntCursor: true},
		EndpointIncident: {UseIntCursor: true},
		Alerts:           {UseIntCursor: false},
	}
)

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	// Entities fetched using REST APIs
	if _, found := ValidRESTEntityExternalIDs[request.EntityExternalID]; found {
		return d.getRESTPage(ctx, request)
	}

	return d.getGraphQLPage(ctx, request)
}
