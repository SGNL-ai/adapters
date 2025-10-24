// Copyright 2025 SGNL.ai, Inc.
package hashicorp

import (
	"bytes"
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

const (
	EntityTypeHosts               = "hosts"
	EntityTypeHostSets            = "host-sets"
	EntityTypeCredentials         = "credentials"
	EntityTypeCredentialLibraries = "credential-libraries"
	EntityTypeAccounts            = "accounts"
	EntityTypeHostCatalogs        = "host-catalogs"
	EntityTypeCredentialStores    = "credential-stores"
	EntityTypeAuthMethods         = "auth-methods"

	ParamHostCatalogID     = "host_catalog_id"
	ParamCredentialStoreID = "credential_store_id"
	ParamAuthMethodID      = "auth_method_id"

	APIVersion = "v1"
)

var noParentEntity = ""

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client *http.Client
}

type DatasourceResponse struct {
	Items        []map[string]any `json:"items"`
	ResponseType string           `json:"response_type"`
	ListToken    string           `json:"list_token"`
	SortBy       string           `json:"sort_by"`
	SortDir      string           `json:"sort_dir"`
	EstItemCount int              `json:"est_item_count"`
}

type AuthResponse struct {
	Attributes struct {
		Token string `json:"token"`
	} `json:"attributes"`
	Command string `json:"command"`
}

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client) Client {
	return &Datasource{
		Client: client,
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	token, frameworkErr := d.getAuthToken(&request.Auth, request.BaseURL)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	request.authToken = token

	// call "auth-methods" first to get the auth token
	switch request.EntityExternalID {
	case EntityTypeHosts, EntityTypeHostSets:
		return d.getCollectionResourcePage(ctx, EntityTypeHostCatalogs, ParamHostCatalogID, request)
	case EntityTypeCredentials, EntityTypeCredentialLibraries:
		return d.getCollectionResourcePage(ctx, EntityTypeCredentialStores, ParamCredentialStoreID, request)
	case EntityTypeAccounts:
		return d.getCollectionResourcePage(ctx, EntityTypeAuthMethods, ParamAuthMethodID, request)
	default:
		response, err := d.getResourcePage(ctx, request)
		if err != nil {
			return nil, err
		}

		if response.NextCursor != nil && *response.NextCursor != "" {
			response.NextCursor, err = d.getMarshalledCursor(response.NextCursor, nil, &noParentEntity)
			if err != nil {
				return nil, err
			}
		}

		return response, nil
	}
}

func ParseResponse(body []byte) (objects []map[string]any, nextCursor *string, err *framework.Error) {
	var data *DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if data == nil {
		return nil, nil, &framework.Error{
			Message: "Unmarshaled response data is nil",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	nextCursor = &data.ListToken
	if data.ResponseType == "complete" {
		nextCursor = nil
	}

	return data.Items, nextCursor, nil
}

func (d *Datasource) getAuthToken(auth *Auth, baseURL string) (string, *framework.Error) {
	if auth == nil {
		return "", nil
	}

	token, err := d.authenticate(auth, baseURL)
	if err != nil {
		return "", err
	}

	return *token, nil
}

func (d *Datasource) authenticate(auth *Auth, baseURL string) (*string, *framework.Error) {
	url := fmt.Sprintf("%s/%s/auth-methods/%s:authenticate", baseURL, APIVersion, auth.AuthMethodID)

	body := map[string]map[string]string{
		"attributes": {
			"login_name": auth.Username,
			"password":   auth.Password,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to marshal authentication request: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create authentication request: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to execute authentication request: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Authentication failed with status code: %d", resp.StatusCode),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to decode authentication response: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	token := fmt.Sprintf("Bearer %s", authResp.Attributes.Token)

	return &token, nil
}

// getResourcePage handles individual calls to API endpoints.
// Each endpoint is identified by the request.EntityExternalID.
// The request is constructed using the request.EntityConfig and request.AdditionalParams.
// The response is parsed into a Response object.
func (d *Datasource) getResourcePage(
	ctx context.Context,
	request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	url := ConstructEndpoint(request)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req = req.WithContext(apiCtx)
	req.Header.Add("Authorization", request.authToken)

	logger.Info("Sending request to datasource", fields.RequestURL(url))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(url),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	defer res.Body.Close()

	response := &Response{
		StatusCode: res.StatusCode,
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read response: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("Datasource responded with an error",
			fields.RequestURL(url),
			fields.ResponseStatusCode(response.StatusCode),
			fields.ResponseRetryAfterHeader(response.RetryAfterHeader),
			fields.ResponseBody(body),
			fields.SGNLEventTypeError(),
		)

		return response, nil
	}

	objects, nextCursor, parserErr := ParseResponse(body)
	if parserErr != nil {
		return nil, parserErr
	}

	response.Objects = objects
	response.NextCursor = nextCursor

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

func (d *Datasource) getMarshalledCursor(
	nextCursor,
	collectionID,
	collectionCursor *string) (*string, *framework.Error) {
	compositeCursor := pagination.CompositeCursor[string]{
		Cursor: nextCursor,
	}

	if collectionID != nil {
		// CollectionID is used for debugging purposes and is not used for pagination.
		compositeCursor.CollectionID = collectionID
		compositeCursor.CollectionCursor = collectionCursor
	}

	marshalledCursor, err := pagination.MarshalCursor(&compositeCursor)
	if err != nil {
		return nil, err
	}

	return &marshalledCursor, nil
}

// getCollectionResourcePage handles calls to API endpoints that are members of a collection.
// The parent collection is set by caller as parentEntity, example "host-catalogs".
// The child entity lookup ID is the mapping to query parameter used to lookup the child entity,
// example "host_catalog_id".
func (d *Datasource) getCollectionResourcePage(
	ctx context.Context,
	parentEntity,
	childEntityLookupID string,
	request *Request) (*Response, *framework.Error) {
	var (
		memberObjects = make([]map[string]any, 0, request.PageSize)
		parentID      string
		parentFilter,
		memberNextCursor,
		parentNextCursor *string
	)

	if filter, ok := request.EntityConfig[parentEntity]; ok {
		parentFilter = &filter.Filter
	}

	parentRequest := Request{
		BaseURL:          request.BaseURL,
		EntityExternalID: parentEntity,
		authToken:        request.authToken,
		// Iterate through each parent entity one at a time because
		// the API doesn't support filtering by parent entity Id.
		PageSize:              1,
		Filter:                parentFilter,
		RequestTimeoutSeconds: request.RequestTimeoutSeconds,
		EntityConfig:          request.EntityConfig,
	}

	if request.Cursor != nil {
		if request.Cursor.CollectionCursor != nil {
			parentRequest.Cursor = &pagination.CompositeCursor[string]{
				Cursor: request.Cursor.CollectionCursor,
			}
		}

		request.Cursor = &pagination.CompositeCursor[string]{
			Cursor: request.Cursor.Cursor,
		}
	}

	parentResponse, err := d.getResourcePage(ctx, &parentRequest)
	if err != nil {
		return nil, err
	}

	if parentResponse == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to get response for %s", parentEntity),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if len(parentResponse.Objects) == 0 {
		return &Response{
			StatusCode:       200,
			RetryAfterHeader: "",
			Objects:          memberObjects,
		}, nil
	}

	// Safely extract the ID from the parent object
	parentObj := parentResponse.Objects[0]
	if parentObj == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Parent object is nil for %s", parentEntity),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	idVal, exists := parentObj["id"]
	if !exists {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Parent object missing 'id' field for %s", parentEntity),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	parentID, ok := idVal.(string)
	if !ok {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Parent object 'id' field is not a string for %s", parentEntity),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	request.AdditionalParams = map[string]string{
		childEntityLookupID: parentID,
	}

	memberResponse, err := d.getResourcePage(ctx, request)
	if err != nil {
		return nil, err
	}

	if memberResponse == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to get response for %s", request.EntityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	memberObjects = append(memberObjects, memberResponse.Objects...)

	if memberResponse.NextCursor == nil {
		// Go to next parent collection.
		parentNextCursor = parentResponse.NextCursor
	} else {
		// Stay at current collection.
		memberNextCursor = memberResponse.NextCursor

		if request.Cursor != nil {
			parentNextCursor = request.Cursor.CollectionCursor
		}
	}

	var nextCursor *string

	if memberNextCursor != nil || parentNextCursor != nil {
		collectionID := fmt.Sprintf("%s-%s", parentID, parentEntity)

		var err *framework.Error

		nextCursor, err = d.getMarshalledCursor(memberNextCursor, &collectionID, parentNextCursor)

		if err != nil {
			return nil, err
		}
	}

	return &Response{
		StatusCode:       200,
		RetryAfterHeader: "",
		Objects:          memberObjects,
		NextCursor:       nextCursor,
	}, nil
}
