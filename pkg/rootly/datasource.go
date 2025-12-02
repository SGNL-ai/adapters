// Copyright 2025 SGNL.ai, Inc.
package rootly

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client *http.Client
}

type DatasourceResponse struct {
	Data     []map[string]any `json:"data"`
	Included []map[string]any `json:"included,omitempty"`
	Meta     struct {
		Page       int `json:"current_page"`
		Pages      int `json:"total_pages"`
		TotalCount int `json:"total_count"`
	} `json:"meta"`
}

type DatasourceErrorResponse struct {
	Errors []struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Status string `json:"status"`
	} `json:"errors"`
}

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ConstructEndpoint(request), nil)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req = req.WithContext(apiCtx)

	req.Header.Add("Authorization", request.HTTPAuthorization)
	req.Header.Add("Content-Type", "application/vnd.api+json")

	// Use the client from the datasource instead of the request
	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to query datasource: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read response body: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse DatasourceErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"Failed to parse error response: %v. Status: %d. Body: %s.", err, resp.StatusCode, string(body),
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		errorMsg := fmt.Sprintf("Received Http Error %d: %s", resp.StatusCode, resp.Status)

		return nil, &framework.Error{
			Message: errorMsg,
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var datasourceResponse DatasourceResponse
	if err := json.Unmarshal(body, &datasourceResponse); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse response: %v. Body: %s.", err, string(body)),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Check if the response has the required data field
	if datasourceResponse.Data == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Invalid response format: missing required data field. Body: %s.", string(body)),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Determine next cursor based on pagination
	var nextCursor *string

	currentPage := datasourceResponse.Meta.Page
	totalPages := datasourceResponse.Meta.Pages

	if currentPage < totalPages {
		nextPageStr := strconv.Itoa(currentPage + 1)
		nextCursor = &nextPageStr
	}

	// Process and add included data to each object
	processedData := make([]map[string]any, len(datasourceResponse.Data))

	for i, dataObject := range datasourceResponse.Data {
		// Create a copy of the original object
		processedObject := make(map[string]any)
		for key, value := range dataObject {
			processedObject[key] = value
		}

		// Get the id from the current data object
		dataObjectID, ok := dataObject["id"]
		if !ok {
			// If no id found, add empty included array
			processedObject["included"] = []any{}
		} else {
			// Filter included objects that match this data object's id
			var matchingIncluded []any

			// Collect selected_users and selected_groups as flat arrays with field_id embedded
			var allSelectedUsers []any
			var allSelectedGroups []any

			for _, includedItem := range datasourceResponse.Included {
				if attributes, hasAttrs := includedItem["attributes"].(map[string]any); hasAttrs {
					if incidentID, hasIncidentID := attributes["incident_id"]; hasIncidentID {
						if incidentID == dataObjectID {
							// Create a copy of the included item to modify
							modifiedItem := make(map[string]any)

							// Get form_field_id if it exists and flatten it to the top level
							var formFieldID string
							if ffi, hasFormFieldID := attributes["form_field_id"]; hasFormFieldID {
								if ffiStr, isString := ffi.(string); isString {
									formFieldID = ffiStr
									modifiedItem["form_field_id"] = ffi // Flatten to top level for JSONPath access
								}
							}

							// Process selected_users and selected_groups from attributes
							if attributes, hasAttrs := includedItem["attributes"].(map[string]any); hasAttrs {
								for key, value := range attributes {
									if key == "selected_users" || key == "selected_groups" {
										entityType := key // "selected_users" or "selected_groups"

										// Handle []any type
										if arrayValue, isArray := value.([]any); isArray {
											for _, item := range arrayValue {
												var userMap map[string]any
												if mapValue, isMap := item.(map[string]any); isMap {
													userMap = mapValue
												} else if mapValueInterface, isMapInterface := item.(map[string]interface{}); isMapInterface {
													// Convert map[string]interface{} to map[string]any
													userMap = make(map[string]any, len(mapValueInterface))
													for k, v := range mapValueInterface {
														userMap[k] = v
													}
												}

												if userMap != nil {
													// Add each user/group as a separate included item to avoid double-wrapping
													expandedItem := make(map[string]any, len(userMap)+2)
													for k, v := range userMap {
														expandedItem[k] = v
													}
													expandedItem["entity_type"] = entityType
													if formFieldID != "" {
														expandedItem["form_field_id"] = formFieldID

														// Add to flat arrays with field_id embedded for simple JSONPath access
														userWithFieldID := make(map[string]any, len(userMap)+1)
														for k, v := range userMap {
															userWithFieldID[k] = v
														}
														userWithFieldID["field_id"] = formFieldID

														if key == "selected_users" {
															allSelectedUsers = append(allSelectedUsers, userWithFieldID)
														} else {
															allSelectedGroups = append(allSelectedGroups, userWithFieldID)
														}
													}
													matchingIncluded = append(matchingIncluded, expandedItem)
												}
											}
										} else if arrayValueInterface, isArrayInterface := value.([]interface{}); isArrayInterface {
											// Handle []interface{} type
											for _, item := range arrayValueInterface {
												var userMap map[string]any
												if mapValue, isMap := item.(map[string]any); isMap {
													userMap = mapValue
												} else if mapValueInterface, isMapInterface := item.(map[string]interface{}); isMapInterface {
													userMap = make(map[string]any, len(mapValueInterface))
													for k, v := range mapValueInterface {
														userMap[k] = v
													}
												}

												if userMap != nil {
													// Add each user/group as a separate included item to avoid double-wrapping
													expandedItem := make(map[string]any, len(userMap)+2)
													for k, v := range userMap {
														expandedItem[k] = v
													}
													expandedItem["entity_type"] = entityType
													if formFieldID != "" {
														expandedItem["form_field_id"] = formFieldID

														// Add to flat arrays with field_id embedded for simple JSONPath access
														userWithFieldID := make(map[string]any, len(userMap)+1)
														for k, v := range userMap {
															userWithFieldID[k] = v
														}
														userWithFieldID["field_id"] = formFieldID

														if key == "selected_users" {
															allSelectedUsers = append(allSelectedUsers, userWithFieldID)
														} else {
															allSelectedGroups = append(allSelectedGroups, userWithFieldID)
														}
													}
													matchingIncluded = append(matchingIncluded, expandedItem)
												}
											}
										}
									}
								}
							}

							// Copy other fields from includedItem to modifiedItem
							for key, value := range includedItem {
								modifiedItem[key] = value
							}

							matchingIncluded = append(matchingIncluded, modifiedItem)
						}
					}
				}
			}
			processedObject["included"] = matchingIncluded

			// Add flattened selected_users and selected_groups arrays
			// Each item has a field_id property for filtering
			// Access with JSONPath like: $.all_selected_users[?field_id='0c696034-4c87-4a29-b7a1-8a5524a443ca']
			// Or with simpler syntax depending on your JSONPath library
			if len(allSelectedUsers) > 0 {
				processedObject["all_selected_users"] = allSelectedUsers
			}
			if len(allSelectedGroups) > 0 {
				processedObject["all_selected_groups"] = allSelectedGroups
			}
		}

		fmt.Println(processedObject)
		processedData[i] = processedObject
	}

	return &Response{
		Objects:    processedData,
		NextCursor: nextCursor,
	}, nil
}
