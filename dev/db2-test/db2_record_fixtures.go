//go:build db2
// +build db2

// Script to record real DB2 responses as test fixtures for contract testing.
// Usage: CGO_ENABLED=1 go run -tags db2 scripts/db2_record_fixtures.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/condexpr"
	"github.com/sgnl-ai/adapters/pkg/db2"
)

// FixtureRequest captures the request parameters for replay.
type FixtureRequest struct {
	Entity     string                 `json:"entity"`
	Schema     string                 `json:"schema"`
	Database   string                 `json:"database"`
	PageSize   int64                  `json:"pageSize"`
	Cursor     string                 `json:"cursor,omitempty"`
	Filter     *condexpr.Condition    `json:"filter,omitempty"`
	Attributes []string               `json:"attributes"`
}

// FixtureResponse captures the response for replay.
type FixtureResponse struct {
	Objects    []map[string]interface{} `json:"objects"`
	NextCursor string                   `json:"nextCursor,omitempty"`
	Error      *FixtureError            `json:"error,omitempty"`
}

// FixtureError captures error details.
type FixtureError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Fixture combines request and response for contract testing.
type Fixture struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	RecordedAt  string          `json:"recordedAt"`
	Request     FixtureRequest  `json:"request"`
	Response    FixtureResponse `json:"response"`
}

// TestCase defines a test scenario to record.
type TestCase struct {
	Name        string
	Description string
	Entity      string
	Schema      string
	PageSize    int64
	Cursor      string
	Filter      *condexpr.Condition
	Attributes  []string
}

func main() {
	password := os.Getenv("DB2_PASSWORD")
	certB64 := os.Getenv("DB2_CERT_BASE64")

	if password == "" {
		fmt.Println("ERROR: DB2_PASSWORD not set")
		os.Exit(1)
	}

	// Define test cases to record
	testCases := []TestCase{
		{
			Name:        "ekpo_with_filter",
			Description: "EKPO entity with BUKRS=US02 filter, page size 10",
			Entity:      "EKPO",
			Schema:      "SAP_ECC_MST_RMS",
			PageSize:    10,
			Filter: &condexpr.Condition{
				Field:    "BUKRS",
				Operator: "=",
				Value:    "US02",
			},
			Attributes: []string{"id", "MANDT", "EBELN", "EBELP", "NETPR", "BUKRS"},
		},
		{
			Name:        "ekpo_small_page",
			Description: "EKPO entity with small page size (3) for pagination testing - page 1",
			Entity:      "EKPO",
			Schema:      "SAP_ECC_MST_RMS",
			PageSize:    3,
			Filter: &condexpr.Condition{
				Field:    "BUKRS",
				Operator: "=",
				Value:    "US02",
			},
			Attributes: []string{"id", "MANDT", "EBELN", "EBELP", "NETPR"},
		},
		{
			Name:        "ekpo_small_page_2",
			Description: "EKPO entity with small page size (3) for pagination testing - page 2",
			Entity:      "EKPO",
			Schema:      "SAP_ECC_MST_RMS",
			PageSize:    3,
			Cursor:      "000|4500000003|00030", // Cursor from ekpo_small_page
			Filter: &condexpr.Condition{
				Field:    "BUKRS",
				Operator: "=",
				Value:    "US02",
			},
			Attributes: []string{"id", "MANDT", "EBELN", "EBELP", "NETPR"},
		},
		{
			Name:        "ekpo_no_filter",
			Description: "EKPO entity without filter, page size 5",
			Entity:      "EKPO",
			Schema:      "SAP_ECC_MST_RMS",
			PageSize:    5,
			Attributes: []string{"id", "MANDT", "EBELN", "EBELP"},
		},
	}

	// Create fixtures directory
	fixturesDir := "pkg/db2/testdata/fixtures"
	if err := os.MkdirAll(fixturesDir, 0755); err != nil {
		fmt.Printf("ERROR: Failed to create fixtures directory: %v\n", err)
		os.Exit(1)
	}

	// Create adapter
	sqlClient := db2.NewDefaultSQLClient()
	datasource := db2.NewClient(sqlClient)
	adapter := db2.NewAdapter(datasource)

	fmt.Println("=== Recording DB2 Test Fixtures ===")
	fmt.Printf("Output directory: %s\n\n", fixturesDir)

	var allFixtures []Fixture

	for _, tc := range testCases {
		fmt.Printf("Recording: %s\n", tc.Name)
		fmt.Printf("  Entity: %s, PageSize: %d\n", tc.Entity, tc.PageSize)

		fixture, err := recordFixture(adapter, tc, password, certB64)
		if err != nil {
			fmt.Printf("  ERROR: %v\n", err)
			continue
		}

		// Save individual fixture
		fixturePath := filepath.Join(fixturesDir, tc.Name+".json")
		if err := saveFixture(fixturePath, fixture); err != nil {
			fmt.Printf("  ERROR: Failed to save: %v\n", err)
			continue
		}

		fmt.Printf("  OK: Saved to %s (%d objects)\n", fixturePath, len(fixture.Response.Objects))
		allFixtures = append(allFixtures, *fixture)
	}

	// Save combined fixtures file
	allFixturesPath := filepath.Join(fixturesDir, "all_fixtures.json")
	if err := saveAllFixtures(allFixturesPath, allFixtures); err != nil {
		fmt.Printf("ERROR: Failed to save combined fixtures: %v\n", err)
	} else {
		fmt.Printf("\nOK: All fixtures saved to %s\n", allFixturesPath)
	}

	fmt.Printf("\n=== Recording Complete ===\n")
	fmt.Printf("Total fixtures recorded: %d\n", len(allFixtures))
}

func recordFixture(adapter framework.Adapter[db2.Config], tc TestCase, password, certB64 string) (*Fixture, error) {
	config := &db2.Config{
		Database: getEnvOrDefault("DB2_DATABASE", "LMTESTDB"),
		Schema:   tc.Schema,
	}

	if tc.Filter != nil {
		config.Filters = map[string]condexpr.Condition{
			tc.Entity: *tc.Filter,
		}
	}

	if certB64 != "" {
		config.CertificateChain = certB64
	}

	// Build attributes config
	var attrConfigs []*framework.AttributeConfig
	for _, attr := range tc.Attributes {
		attrType := framework.AttributeTypeString
		if attr == "NETPR" {
			attrType = framework.AttributeTypeDouble
		}
		attrConfigs = append(attrConfigs, &framework.AttributeConfig{
			ExternalId: attr,
			Type:       attrType,
			UniqueId:   attr == "id",
		})
	}

	request := &framework.Request[db2.Config]{
		Auth: &framework.DatasourceAuthCredentials{
			Basic: &framework.BasicAuthCredentials{
				Username: getEnvOrDefault("DB2_USER", "db2inst1"),
				Password: password,
			},
		},
		Address:  fmt.Sprintf("%s:%s", getEnvOrDefault("DB2_HOST", "ec2-18-191-143-175.us-east-2.compute.amazonaws.com"), getEnvOrDefault("DB2_PORT", "50001")),
		PageSize: tc.PageSize,
		Cursor:   tc.Cursor,
		Entity: framework.EntityConfig{
			ExternalId: tc.Entity,
			Attributes: attrConfigs,
		},
		Config: config,
	}

	response := adapter.GetPage(context.Background(), request)

	fixture := &Fixture{
		Name:        tc.Name,
		Description: tc.Description,
		RecordedAt:  time.Now().UTC().Format(time.RFC3339),
		Request: FixtureRequest{
			Entity:     tc.Entity,
			Schema:     tc.Schema,
			Database:   config.Database,
			PageSize:   tc.PageSize,
			Cursor:     tc.Cursor,
			Filter:     tc.Filter,
			Attributes: tc.Attributes,
		},
	}

	if response.Error != nil {
		fixture.Response.Error = &FixtureError{
			Message: response.Error.Message,
			Code:    response.Error.Code.String(),
		}
		return fixture, fmt.Errorf("adapter error: %s", response.Error.Message)
	}

	if response.Success != nil {
		// Convert framework.Object to map[string]interface{}
		for _, obj := range response.Success.Objects {
			fixture.Response.Objects = append(fixture.Response.Objects, map[string]interface{}(obj))
		}
		fixture.Response.NextCursor = response.Success.NextCursor
	}

	return fixture, nil
}

func saveFixture(path string, fixture *Fixture) error {
	data, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func saveAllFixtures(path string, fixtures []Fixture) error {
	data, err := json.MarshalIndent(fixtures, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// sanitizeForFilename removes characters that are problematic in filenames.
func sanitizeForFilename(s string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return replacer.Replace(s)
}
