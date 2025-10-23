// Copyright 2025 SGNL.ai, Inc.
package bamboohr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"go.uber.org/zap"
)

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client *http.Client
}

// BambooHR API response format.
type DatasourceResponse struct {
	Employees []map[string]any `json:"employees"`
}

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint path to query that entity.
type Entity struct {
	// path is the endpoint to query the entity.
	path string
}

const (
	Employee = "Employee"
)

var (
	// ValidEntityExternalIDs is a set of valid external IDs of entities that can be queried.
	ValidEntityExternalIDs = map[string]Entity{
		Employee: {
			path: "/reports/custom",
		},
	}

	nullValues = map[string]struct{}{
		"":                          {},
		"NULL":                      {},
		"0000-00-00":                {},
		"00/00/0000":                {},
		"00 00 0000":                {},
		"0000-00-00T00:00:00":       {},
		"0000-00-00T00:00:00Z":      {},
		"0000-00-00T00:00:00.000Z":  {},
		"0000-00-00T00:00:00+00:00": {},
		"0000-00-00T00:00:00+0000":  {},
		"0000-00-00T00:00:00+00":    {},
	}
)

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityConfig.ExternalId),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	if request.AttributeMappings == nil {
		return nil, &framework.Error{
			Message: "AttributeMappings is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityConfig.ExternalId,
		false,
	)
	if validationErr != nil {
		return nil, validationErr
	}

	endpointInfo, endpointErr := ConstructEndpoint(request)
	if endpointErr != nil {
		return nil, endpointErr
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodPost, endpointInfo.URL, strings.NewReader(endpointInfo.Body))
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	req.SetBasicAuth(request.APIKey, request.BasicAuthPassword)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	logger.Info("Sending HTTP request to datasource", fields.RequestURL(endpointInfo.URL))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("HTTP request to datasource failed",
			fields.RequestURL(endpointInfo.URL),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute BambooHR request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	defer res.Body.Close()

	response := &Response{
		StatusCode:       res.StatusCode,
		RetryAfterHeader: res.Header.Get("Retry-After"),
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("Datasource request failed",
			fields.RequestURL(endpointInfo.URL),
			fields.ResponseStatusCode(response.StatusCode),
			fields.ResponseRetryAfterHeader(response.RetryAfterHeader),
			fields.ResponseBody(res.Body),
			fields.SGNLEventTypeError(),
		)

		return response, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read BambooHR response body: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, nextCursor, frameworkErr := ParseResponse(body, request)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	response.NextCursor = nextCursor
	response.Objects = objects

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

func ParseResponse(body []byte, request *Request) (
	[]map[string]any,
	*pagination.CompositeCursor[int64],
	*framework.Error,
) {
	var response *DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &response); unmarshalErr != nil || response == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if response.Employees == nil {
		return nil, nil, &framework.Error{
			Message: "Failed to parse response body: employees field is missing.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, nextCursor, err := pagination.PaginateObjects(response.Employees, request.PageSize, request.Cursor)
	if err != nil {
		return nil, nil, err
	}

	if mappingErr := ProcessObjects(request, &objects); mappingErr != nil {
		return nil, nil, mappingErr
	}

	var nextCompositeCursor *pagination.CompositeCursor[int64]

	if nextCursor != nil {
		nextCompositeCursor = &pagination.CompositeCursor[int64]{Cursor: nextCursor}
	}

	return objects, nextCompositeCursor, nil
}

func ProcessObjects(request *Request, objects *[]map[string]any) *framework.Error {
	for idx := range *objects {
		err := ProcessObject(request, (*objects)[idx])
		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessObject(request *Request, object map[string]any) *framework.Error {
	for _, attr := range request.EntityConfig.Attributes {
		if value, found := object[attr.ExternalId]; found && value != nil {
			var err *framework.Error
			if attr.List {
				object[attr.ExternalId], err = ProcessAttributeList(value, attr, request.AttributeMappings)
			} else {
				object[attr.ExternalId], err = ProcessAttribute(value, attr, request.AttributeMappings)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ProcessAttributeList(
	value any,
	attr *framework.AttributeConfig,
	mappings *AttributeMappings,
) (any, *framework.Error) {
	attrList, ok := value.([]any)
	if !ok {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Unexpected Format: attribute %s is not a list.", attr.ExternalId),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
		}
	}

	for idx, currValue := range attrList {
		convertedVal, err := ProcessAttribute(currValue, attr, mappings)
		if err != nil {
			return nil, err
		}

		attrList[idx] = convertedVal
	}

	return attrList, nil
}

func ProcessAttribute(
	value any,
	attr *framework.AttributeConfig,
	mappings *AttributeMappings,
) (any, *framework.Error) {
	if value == nil {
		return nil, nil
	}

	attrStringValue, ok := value.(string)
	if !ok {
		return nil, &framework.Error{
			Message: "Unexpected Format: SoR Response includes non-string attribute types",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
		}
	}

	if _, found := nullValues[strings.ToUpper(attrStringValue)]; found {
		return nil, nil
	}

	var convertedVal any

	var err error

	switch attr.Type {
	case framework.AttributeTypeInt64, framework.AttributeTypeDouble:
		convertedVal, err = strconv.ParseFloat(attrStringValue, 64)
	case framework.AttributeTypeBool:
		convertedVal, err = ParseBool(attrStringValue, mappings.BoolMappings)
	default:
		return value, nil
	}

	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse attribute: %s, as type %s with value: %s.",
				attr.ExternalId, attrTypeToString(attr.Type), attrStringValue),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
		}
	}

	return convertedVal, nil
}

func ParseBool(value string, mappings *BoolAttributeMappings) (bool, error) {
	convertedVal, err := strconv.ParseBool(value)
	if err == nil {
		return convertedVal, nil
	}

	if mappings != nil {
		if slices.Contains(mappings.True, value) {
			return true, nil
		} else if slices.Contains(mappings.False, value) {
			return false, nil
		}
	}

	return false, fmt.Errorf("unable to parse boolean value: %s", value)
}

func attrTypeToString(t framework.AttributeType) string {
	switch t {
	case framework.AttributeTypeInt64, framework.AttributeTypeDouble:
		return "Float"
	case framework.AttributeTypeBool:
		return "Bool"
	default:
		return "Unsupported"
	}
}
