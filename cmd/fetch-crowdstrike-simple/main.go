package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/crowdstrike"
)

func main() {
	var (
		address  = flag.String("address", "https://api.us-2.crowdstrike.com", "CrowdStrike API address")
		token    = flag.String("token", "", "Bearer token for authentication (required)")
		entity   = flag.String("entity", "endpoint_protection_combined_alerts", "Entity type")
		output   = flag.String("output", "", "Output file path (required)")
		pageSize = flag.Int("pagesize", 10, "Page size for requests")
		cursor   = flag.String("cursor", "", "Cursor for pagination")
	)
	flag.Parse()

	if *token == "" || *output == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -token <bearer-token> -output <output-file>\n", os.Args[0])
		os.Exit(1)
	}

	// Create HTTP client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create adapter (similar to smoketests)
	adapter := crowdstrike.NewAdapter(crowdstrike.NewClient(httpClient))

	// Create entity config
	entityConfig := &framework.EntityConfig{
		ExternalId: *entity,
		Attributes: []*framework.AttributeConfig{
			{
				ExternalId: "composite_id",
				Type:       framework.AttributeTypeString,
				UniqueId:   true,
			},
			{
				ExternalId: "aggregate_id",
				Type:       framework.AttributeTypeString,
			},
			{
				ExternalId: "status",
				Type:       framework.AttributeTypeString,
			},
		},
	}

	// Create request
	request := &framework.Request[crowdstrike.Config]{
		Address: *address,
		Auth: &framework.DatasourceAuthCredentials{
			HTTPAuthorization: "Bearer " + *token,
		},
		Config: &crowdstrike.Config{
			APIVersion: "v1",
			Archived:   false,
			Enabled:    true,
		},
		Entity:   *entityConfig,
		PageSize: int64(*pageSize),
		Cursor:   *cursor,
	}

	// Make the request
	fmt.Printf("Fetching %s from %s...\n", *entity, *address)
	response := adapter.GetPage(context.Background(), request)

	// Check for errors
	if response.Error != nil {
		fmt.Fprintf(os.Stderr, "Error from adapter: %s\n", response.Error.Message)
		os.Exit(1)
	}

	if response.Success == nil {
		fmt.Fprintf(os.Stderr, "No success response received\n")
		os.Exit(1)
	}

	// Create output structure
	outputData := map[string]interface{}{
		"request": map[string]interface{}{
			"address":   *address,
			"entity":    *entity,
			"pageSize":  *pageSize,
			"cursor":    *cursor,
			"timestamp": time.Now().Format(time.RFC3339),
		},
		"response": map[string]interface{}{
			"objectCount": len(response.Success.Objects),
			"nextCursor":  response.Success.NextCursor,
			"objects":     response.Success.Objects,
		},
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	err = os.WriteFile(*output, jsonData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response saved to: %s\n", *output)
	fmt.Printf("Objects retrieved: %d\n", len(response.Success.Objects))
	if response.Success.NextCursor != "" {
		fmt.Printf("Next cursor: %s\n", response.Success.NextCursor)
	} else {
		fmt.Printf("No more pages available\n")
	}
}