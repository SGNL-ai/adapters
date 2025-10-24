// Copyright 2025 SGNL.ai, Inc.
package duo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// Duo API Metadata format used for pagination.
type ResponseMetadata struct {
	NextOffset *int64 `json:"next_offset,omitempty"`
}

// Duo API response format.
type DatasourceResponse struct {
	Stat     *string                  `json:"stat"`
	Metadata *ResponseMetadata        `json:"metadata,omitempty"`
	Objects  []map[string]interface{} `json:"response"`
}

// Entity contains entity specific information, such as the entity's unique ID attribute and the
// endpoint path to query that entity.
type Entity struct {
	// path is the endpoint to query the entity.
	path string
	// uniqueIDAttrExternalID is the external ID of the entity's uniqueId attribute.
	uniqueIDAttrExternalID string
}

const (
	User     = "User"
	Group    = "Group"
	Phone    = "Phone"
	Endpoint = "Endpoint"
	RFC2822  = "Mon, 02 Jan 2006 15:04:05 -0700"
)

var (
	// ValidEntityExternalIDs is a set of valid external IDs of entities that can be queried.
	// The map value is the Entity struct which contains the unique ID attribute.
	ValidEntityExternalIDs = map[string]Entity{
		User: {
			path:                   "users",
			uniqueIDAttrExternalID: "user_id",
		},
		Endpoint: {
			path:                   "endpoints",
			uniqueIDAttrExternalID: "epkey",
		},
		Group: {
			path:                   "groups",
			uniqueIDAttrExternalID: "group_id",
		},
		Phone: {
			path:                   "phones",
			uniqueIDAttrExternalID: "phone_id",
		},
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
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityExternalID,
		false,
	)
	if validationErr != nil {
		return nil, validationErr
	}

	endpointInfo, endpointErr := ConstructEndpoint(request)
	if endpointErr != nil {
		return nil, endpointErr
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointInfo.URL, nil)
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

	req.Header.Add("Authorization", endpointInfo.Auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Date", endpointInfo.Date)

	logger.Info("Sending request to datasource", fields.RequestURL(endpointInfo.URL))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(endpointInfo.URL),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute Duo request: %v.", err),
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
		logger.Error("Datasource responded with an error",
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
			Message: fmt.Sprintf("Failed to read Duo response body: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, nextCursor, frameworkErr := ParseResponse(body)
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

func ParseResponse(body []byte) (
	objects []map[string]any,
	nextCursor *pagination.CompositeCursor[int64],
	err *framework.Error,
) {
	var data *DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if data.Metadata != nil && data.Metadata.NextOffset != nil {
		nextCursor = &pagination.CompositeCursor[int64]{
			Cursor: data.Metadata.NextOffset,
		}
	}

	return data.Objects, nextCursor, nil
}
