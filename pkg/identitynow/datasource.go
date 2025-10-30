// Copyright 2025 SGNL.ai, Inc.
package identitynow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

const (
	AccountEntitlements              = "accountEntitlements"
	Accounts                         = "accounts"
	DefaultAccountCollectionPageSize = 100

	// TODO [sc-19214]: Remove after POC complete.
	Delimiter string = " | "
)

var (
	// TODO [sc-19214]: Remove after POC complete.
	AttributesToConcatenate = []string{"groups", "Groups", "memberOf"}
	DefaultAccountSorter    = "id"
	Validate                = validator.New(validator.WithRequiredStructEnabled())
)

type AccountObject struct {
	AccountID       string `mapstructure:"id" validate:"required"`
	HasEntitlements bool   `mapstructure:"hasEntitlements" validate:"omitempty"`
}

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client                    *http.Client
	AccountCollectionPageSize int
}

// DatasouceResponse represents the API response from the datasource, which is an array of objects.
type DatasourceResponse = []map[string]any

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client, pageSize int) Client {
	return &Datasource{
		Client:                    client,
		AccountCollectionPageSize: pageSize,
	}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	if request.EntityExternalID == AccountEntitlements {
		return d.handleAccountEntitlements(ctx, request)
	}

	// [Accounts/Entitlements] This verifies that `CollectionID` and `CollectionCursor` are not set.
	validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityExternalID,
		// Send a bool indicating if the entity is a member of a collection.
		false,
	)
	if validationErr != nil {
		return nil, validationErr
	}

	// Construct the first page cursor if request.Cursor is nil.
	if request.Cursor == nil || request.Cursor.Cursor == nil {
		var zero int64

		request.Cursor = &pagination.CompositeCursor[int64]{
			Cursor: &zero,
		}
	}

	response, objects, err := d.ConstructEndpointAndGetResponse(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return response, nil
	}

	// TODO [sc-19214]: Remove after POC complete.
	// Before returning objects, create the "groupMembership" attribute by concatenating each element (i.e. groupId)
	// in the "attributes.groups" array into a delimited string.
	// e.g.
	// "attributes": {
	//	"groups": ["group1", "group2"]
	// }
	// --> "groupMembership": "group1 | group2"
	// Other attributes to concatenate are defined in AttributesToConcatenate.
	for _, object := range objects {
		if object["attributes"] == nil {
			continue
		}

		attributes, ok := object["attributes"].(map[string]any)
		if !ok {
			continue
		}

		for _, attributeToConcatenate := range AttributesToConcatenate {
			if attributes[attributeToConcatenate] == nil {
				continue
			}

			attribute, ok := attributes[attributeToConcatenate].([]any)
			if !ok {
				continue
			}

			if concatenated := concatenateArrayElements(attribute, Delimiter); concatenated != nil {
				// The new concatenated attribute has the following name format:
				// e.g. "groups" --> "groupsMembership"
				// And it's a top-level key of the object and no longer nested under "attributes".
				object[fmt.Sprintf("%sMembership", attributeToConcatenate)] = *concatenated
			}
		}
	}

	// Setup the next cursor.
	var nextCursor *pagination.CompositeCursor[int64]

	cursor := pagination.GetNextCursorFromPageSize(len(objects), request.PageSize, *request.Cursor.Cursor)

	if cursor != nil {
		nextCursor = &pagination.CompositeCursor[int64]{
			Cursor: cursor,
		}
	}

	return &Response{
		StatusCode:       response.StatusCode,
		RetryAfterHeader: response.RetryAfterHeader,
		Objects:          objects,
		NextCursor:       nextCursor,
	}, nil
}

// TODO: This is identical to the ParseResponse in the Okta package. Refactor into a shared "parser" package.
func ParseResponse(body []byte) (objects []map[string]any, err *framework.Error) {
	var data DatasourceResponse

	if unmarshalErr := json.Unmarshal(body, &data); unmarshalErr != nil || data == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to unmarshal the datasource response: %v.", unmarshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return data, nil
}

// TODO [sc-19214]: Remove after POC complete.
// concatenateArrayElements concatenates the elements of an array into a delimited string.
// It first asserts that the entire array is of type string. If assertion fails, nil is returned.
func concatenateArrayElements(rawArray []any, delimiter string) *string {
	array := make([]string, 0, len(rawArray))

	for _, element := range rawArray {
		elementAsStr, ok := element.(string)
		if !ok {
			return nil
		}

		array = append(array, elementAsStr)
	}

	concatenated := strings.Join(array, delimiter)

	return &concatenated
}

// handleAccountEntitlements retrieves the entitlements for an account.
// It retrieves multiple account IDs and then retrieves the entitlements for each account per request.
// The response is paginated by the page size.
//
// The composite cursor in the request consists of:
// Collection ID: Account ID.
// Cursor: The offset of the entitlements to retrieve for a given account ID.
// Collection Cursor: The offset of the account IDs to retrieve.
func (d *Datasource) handleAccountEntitlements(ctx context.Context, request *Request) (*Response, *framework.Error) {
	var (
		// Total number of entitlements processed in this page.
		totalNumberOfEntitlements int
		nextCollectionCursor      *int64
		nextCursor                *pagination.CompositeCursor[int64]
		retryAfterHeader          string
		zero                      int64
	)

	accountResponse, frameworkErr := d.getAccountsForEntitlements(ctx, request)
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	// There are no more accounts. Return an empty last page.
	if len(accountResponse.Objects) == 0 {
		return &Response{
			StatusCode: http.StatusNoContent,
		}, nil
	}

	// If there are more accounts, then set the next collection cursor.
	if accountResponse.NextCursor != nil {
		nextCollectionCursor = accountResponse.NextCursor.Cursor
	}

	// `currentCollectionCursor` keeps track of the collection cursor.
	// This value is used to construct the next collection cursor if
	// the page fills up before all accounts are processed.
	currentCollectionCursor := zero
	if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
		currentCollectionCursor = *request.Cursor.CollectionCursor
	}

	accountEntitlementObjects := make([]map[string]any, 0, request.PageSize)

	// Loop through the account objects and call `account/{accountId/entitlements` endpoint if the account has
	// entitlements.
	for i, account := range accountResponse.Objects {
		// Get the account `id` and `hasEntitlements` fields from the account object to construct the request.
		accountObject := new(AccountObject)

		if err := mapstructure.Decode(account, &accountObject); err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to decode account data when fetching its entitlements.: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		if err := Validate.Struct(accountObject); err != nil {
			return nil, &framework.Error{
				// nolint:lll
				Message: fmt.Sprintf("Failed to validate required fields for an account when fetching its entitlements.: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		// If the account has no entitlements, skip to the next account.
		if !accountObject.HasEntitlements {
			continue
		}

		// If the request cursor is not nil, set it to the cursor value.
		// This indicates that the next page of entitlements should be retrieved.
		newCursor := &zero
		if request.Cursor != nil && request.Cursor.Cursor != nil {
			newCursor = request.Cursor.Cursor
		}

		// Set the Cursor values so that the request will fetch a page of entitlements for the current account.
		request.Cursor = &pagination.CompositeCursor[int64]{
			Cursor:           newCursor,
			CollectionID:     &accountObject.AccountID,
			CollectionCursor: nextCollectionCursor,
		}

		// Get the entitlements for the account.
		response, entitlementObjects, frameworkErr := d.ConstructEndpointAndGetResponse(ctx, request)
		if frameworkErr != nil {
			return nil, frameworkErr
		}

		if response.StatusCode != http.StatusOK {
			return response, nil
		}

		retryAfterHeader = response.RetryAfterHeader

		totalNumberOfEntitlements += len(entitlementObjects)

		// Construct and initialize the next cursor for the next request.
		// The `nextCursor` will track state between iterations but it is not used on every iteration.
		// It is returned with the final response.
		nextCursor = &pagination.CompositeCursor[int64]{
			CollectionID:     request.Cursor.CollectionID,
			CollectionCursor: request.Cursor.CollectionCursor,
		}

		// nolint:lll
		nextEntitlementsCursor := pagination.GetNextCursorFromPageSize(len(entitlementObjects), request.PageSize, *request.Cursor.Cursor)

		if nextEntitlementsCursor != nil {
			nextCursor.Cursor = nextEntitlementsCursor
		} else {
			// There are no more entitlements for the account.
			// Unset the request cursor if it was set.
			// This is to ensure that the all the entitlements are retrieved for the next account
			// from offset 0.
			request.Cursor.Cursor = nil
			nextCursor.Cursor = nil
		}

		// The total number of entitlements retrieved is greater than the page size.
		// Change the cursor values and trim the entitlements to the page size.
		if totalNumberOfEntitlements > int(request.PageSize) {
			// Calculate the number of entitlements to trim from the end of the slice.
			excess := totalNumberOfEntitlements - int(request.PageSize)

			// Calculate the cursor so that the next request will fetch the remaining
			// entitlements for the account.
			entitlementCursor := int64(len(entitlementObjects) - excess)
			nextCursor.Cursor = &entitlementCursor

			// Set the collection cursor so that the next request will fetch the remaining
			// entitlements for the account.
			nextCursor.CollectionCursor = &currentCollectionCursor

			// Trim the entitlements to the page size and copy the sub-slice to a new slice.
			// This will help in garbage collection of the original slice.
			entitlementObjects = append([]map[string]any{}, entitlementObjects[:len(entitlementObjects)-excess]...)
		}

		// Loop through the entitlement objects, construct the necessary output and
		// append them to the `accountEntitlementObjects` slice.
		for _, entitlementObject := range entitlementObjects {
			entitlementID, ok := entitlementObject["id"].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to convert IdentityNow account entitlement object id field to string: %v.",
						entitlementObject["id"],
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			accountEntitlementObjects = append(accountEntitlementObjects, map[string]any{
				"id":            fmt.Sprintf("%s-%s", *request.Cursor.CollectionID, entitlementID),
				"accountId":     *request.Cursor.CollectionID,
				"entitlementId": entitlementID,
			})
		}

		// The number of entitlements retrieved is equal to the page size.
		// Check if there are more accounts and/or entitlements to process.
		if len(accountEntitlementObjects) == int(request.PageSize) {
			// There are more accounts to process and the current account has no more entitlements.
			// Set the collection cursor to the next account.
			if i != len(accountResponse.Objects)-1 && nextCursor.Cursor == nil {
				currentCollectionCursor++
				nextCursor.CollectionCursor = &currentCollectionCursor
			} else if i == len(accountResponse.Objects)-1 {
				// There are no more accounts to process for the current collection or batch of accounts.
				// Set the collection cursor to the `nextCollectionCursor`.
				// This value is same as `accountResponse.NextCursor.Cursor`.
				nextCursor.CollectionCursor = nextCollectionCursor

				// There are more entitlements to process for the current account.
				// Set the collection cursor to the current account.
				if nextCursor.Cursor != nil {
					nextCursor.CollectionCursor = &currentCollectionCursor
				}
			}

			// The page is full. Break out of the loop.
			break
		}

		// Move to the next account.
		currentCollectionCursor++
	}

	if len(accountEntitlementObjects) == 0 {
		return &Response{
			StatusCode: http.StatusNoContent,
			NextCursor: &pagination.CompositeCursor[int64]{
				CollectionCursor: nextCollectionCursor,
			},
		}, nil
	}

	response := &Response{
		StatusCode:       accountResponse.StatusCode,
		RetryAfterHeader: retryAfterHeader,
		Objects:          accountEntitlementObjects,
	}

	// If there are more entitlements or accounts to process, set the next cursor.
	if nextCursor.Cursor != nil || nextCursor.CollectionCursor != nil {
		response.NextCursor = nextCursor
	}

	return response, nil
}

// getAccountsForEntitlements retrieves the accounts for fetching the entitlements.
// The number of accounts returned is determined by the value in AccountCollectionPageSize.
func (d *Datasource) getAccountsForEntitlements(ctx context.Context, request *Request) (*Response, *framework.Error) {
	var (
		accountCursor *pagination.CompositeCursor[int64]
		zero          int64
	)

	if request.Cursor == nil {
		accountCursor = &pagination.CompositeCursor[int64]{
			Cursor: &zero,
		}
	} else {
		accountCursor = &pagination.CompositeCursor[int64]{
			Cursor: request.Cursor.CollectionCursor,
		}
	}

	// Setup an Account request to retrieve account IDs.
	accountReq := &Request{
		BaseURL:               request.BaseURL,
		Token:                 request.Token,
		PageSize:              int64(d.AccountCollectionPageSize),
		Cursor:                accountCursor,
		EntityExternalID:      Accounts,
		APIVersion:            request.APIVersion,
		RequestTimeoutSeconds: request.RequestTimeoutSeconds,
		// Use the filter on the AccountEntitlements config to filter
		// the accounts to retrieve entitlements from.
		// TODO [sc-19213]: Remove this hack once POC complete to use a proper filter config.
		Filter: request.Filter,
	}

	accountReq.Sorters = &DefaultAccountSorter
	if request.Sorters != nil {
		accountReq.Sorters = request.Sorters
	}

	accountRes, err := d.GetPage(ctx, accountReq)
	if err != nil {
		return nil, err
	}

	return accountRes, nil
}

// ConstructEndpointAndGetResponse constructs the endpoint and returns the response from the datasource.
func (d *Datasource) ConstructEndpointAndGetResponse(
	ctx context.Context,
	request *Request,
) (*Response, []map[string]any, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Starting datasource request")

	endpoint, errFramework := ConstructEndpoint(request)
	if errFramework != nil {
		return nil, nil, errFramework
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to create request to datasource: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Timeout API calls that take longer than the configured timeout.
	apiCtx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	req = req.WithContext(apiCtx)

	req.Header.Add("Authorization", request.Token)
	req.Header.Add("Content-Type", "application/json")

	logger.Info("Sending request to datasource", fields.RequestURL(endpoint))

	res, err := d.Client.Do(req)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.RequestURL(endpoint),
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		return nil, nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to execute IdentityNow request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Error("Datasource responded with an error",
			fields.RequestURL(endpoint),
			fields.ResponseStatusCode(res.StatusCode),
			fields.ResponseRetryAfterHeader(res.Header.Get("Retry-After")),
			fields.ResponseBody(res.Body),
			fields.SGNLEventTypeError(),
		)

		return &Response{
			StatusCode:       res.StatusCode,
			RetryAfterHeader: res.Header.Get("Retry-After"),
		}, nil, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to read IdentityNow response: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	objects, errFramework := ParseResponse(body)
	if errFramework != nil {
		return nil, nil, errFramework
	}

	response := &Response{
		StatusCode:       res.StatusCode,
		RetryAfterHeader: res.Header.Get("Retry-After"),
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(objects)),
		fields.ResponseNextCursor(nil),
	)

	return response, objects, nil
}
