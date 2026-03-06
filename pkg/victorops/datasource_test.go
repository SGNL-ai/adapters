// Copyright 2026 SGNL.ai, Inc.

// nolint: goconst
package victorops_test

import (
	"net/http"
)

const (
	mockAPIId  = "test-api-id"
	mockAPIKey = "test-api-key"
)

// TestServerHandler is a mock HTTP handler for VictorOps API endpoints.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Validate authentication headers.
	if r.Header.Get("X-VO-Api-Id") != mockAPIId || r.Header.Get("X-VO-Api-Key") != mockAPIKey {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))

		return
	}

	switch r.URL.RequestURI() {
	// Incident endpoints.
	case "/api-reporting/v2/incidents?offset=0&limit=2":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"incidents": [
				{"incidentNumber": "1", "currentPhase": "RESOLVED", "startTime": "2024-01-15T10:30:00Z"},
				{"incidentNumber": "2", "currentPhase": "ACKED", "startTime": "2024-01-16T11:00:00Z"}
			],
			"total": 5,
			"offset": 0,
			"limit": 2
		}`))
	case "/api-reporting/v2/incidents?offset=2&limit=2":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"incidents": [
				{"incidentNumber": "3", "currentPhase": "RESOLVED", "startTime": "2024-01-17T09:00:00Z"},
				{"incidentNumber": "4", "currentPhase": "RESOLVED", "startTime": "2024-01-18T14:00:00Z"}
			],
			"total": 5,
			"offset": 2,
			"limit": 2
		}`))
	case "/api-reporting/v2/incidents?offset=4&limit=2":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"incidents": [
				{"incidentNumber": "5", "currentPhase": "ACKED", "startTime": "2024-01-19T08:00:00Z"}
			],
			"total": 5,
			"offset": 4,
			"limit": 2
		}`))
	case "/api-reporting/v2/incidents?offset=0&limit=10":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"incidents": [
				{"incidentNumber": "1", "currentPhase": "RESOLVED"}
			],
			"total": 1,
			"offset": 0,
			"limit": 10
		}`))
	case "/api-reporting/v2/incidents?offset=0&limit=5":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"incidents": [
				{"incidentNumber": "10", "currentPhase": "RESOLVED", "pagedUsers": ["alice", "bob"], "pagedTeams": ["team-alpha"]}
			],
			"total": 1,
			"offset": 0,
			"limit": 5
		}`))
	case "/api-reporting/v2/incidents?offset=0&limit=2&currentPhase=RESOLVED":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"incidents": [
				{"incidentNumber": "1", "currentPhase": "RESOLVED"},
				{"incidentNumber": "3", "currentPhase": "RESOLVED"}
			],
			"total": 2,
			"offset": 0,
			"limit": 2
		}`))

	// User endpoints.
	case "/api-public/v2/user":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"users": [
				{"username": "user1", "firstName": "Alice", "lastName": "Smith", "email": "alice@example.com"},
				{"username": "user2", "firstName": "Bob", "lastName": "Jones", "email": "bob@example.com"}
			]
		}`))

	// Error endpoints.
	case "/api-reporting/v2/incidents?offset=0&limit=1":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}
})
