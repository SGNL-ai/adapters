// Copyright 2025 SGNL.ai, Inc.

package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

const (
	User             string = "User"
	Group            string = "Group"
	GroupMember      string = "GroupMember"
	Role             string = "Role"
	IdentityProvider string = "IdentityProvider"
	Policy           string = "Policy"
	RolePolicy       string = "RolePolicy"
	UserPolicy       string = "UserPolicy"
	GroupPolicy      string = "GroupPolicy"

	unhandledStatusCode int    = -1
	uniqueIDAttribute   string = "id"
	UserID              string = "UserId"
	PolicyArn           string = "PolicyArn"
	GroupID             string = "GroupId"
	AccountID           string = "AccountId"

	SessionName = "SGNLSession"
)

type EntityInfo struct {
	// MemberOf Specifies the entity name to which the member belong.
	MemberOf *string
	// CollectionAttribute is the attribute that contains the collection name.
	CollectionAttribute *string
	// Identifiers contains the attributes that can be used to uniquely identify the entity.
	Identifiers *Identifiers
}

type Identifiers struct {
	// ArnAttribute is the attribute that contains the Arn of the entity.
	// This is used to extract the AccountID.
	//
	// [Example: Arn, PolicyArn, etc.]
	ArnAttribute string

	// UniqueName of the AWS entity helps to get the entity.
	UniqueName string
}

type AccountCursor struct {
	Offset     int
	NextMarker *string
}

var (
	// ValidEntityExternalIDs is a set of valid external IDs of entities that can be queried.
	ValidEntityExternalIDs = map[string]EntityInfo{
		User: {
			Identifiers: &Identifiers{
				ArnAttribute: "Arn",
				UniqueName:   "UserName",
			},
		},
		Group: {
			Identifiers: &Identifiers{
				ArnAttribute: "Arn",
				UniqueName:   "GroupName",
			},
		},
		Role: {
			Identifiers: &Identifiers{
				ArnAttribute: "Arn",
				UniqueName:   "RoleName",
			},
		},
		Policy: {
			Identifiers: &Identifiers{
				ArnAttribute: "Arn",
				UniqueName:   "PolicyName",
			},
		},
		IdentityProvider: {
			Identifiers: &Identifiers{
				ArnAttribute: "Arn",
			},
		},
		GroupPolicy: {
			CollectionAttribute: func() *string {
				var s = "GroupName"

				return &s
			}(),
			MemberOf: func() *string {
				var s = Group

				return &s
			}(),
		},
		GroupMember: {
			CollectionAttribute: func() *string {
				s := "GroupName"

				return &s
			}(),
			MemberOf: func() *string {
				s := Group

				return &s
			}(),
		},
		RolePolicy: {
			CollectionAttribute: func() *string {
				s := "RoleName"

				return &s
			}(),
			MemberOf: func() *string {
				s := Role

				return &s
			}(),
		},
		UserPolicy: {
			CollectionAttribute: func() *string {
				var s = "UserName"

				return &s
			}(),
			MemberOf: func() *string {
				s := User

				return &s
			}(),
		},
	}
)

// InputParams is the input parameters for List{Entity} from AWS.
type InputParams struct {
	PathPrefix *string // An optional prefix to filter the entities based on their path.
	MaxItems   *int32  // An optional limit on the number of entities to fetch.
	Marker     *string // An optional marker to indicate the starting point for the next set of results.
}

// Options contains adapter level options for fetching entities.
type Options struct {
	InputParams

	EntityName         string  // The name of the entity for which the request is made.
	UniqueName         *string // Unique Name associated with the entity, used for identification.
	UniqueID           *string // Unique ID associated with the entity, used for identification.
	AccountIDRequested bool    // A flag indicating whether the Account ID is requested to be included in the response.
	IsMember           bool    // A flag indicating whether the entity is a member of another entity.
	MaxConcurrent      int     // The maximum number of concurrent requests to make.
}

type ClientConfig struct {
	AWSConfig *aws.Config
}

type Datasource struct {
	Client         *http.Client
	AWSConfig      *aws.Config
	MaxConcurrency int
}

// NewClient returns a Client to query the datasource.
func NewClient(
	client *http.Client,
	awsConfig *aws.Config,
	maxConcurrency int,
) (Client, error) {
	if awsConfig == nil {
		cfg, err := aws_config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}

		awsConfig = &cfg
	}

	return &Datasource{
		AWSConfig:      awsConfig,
		Client:         client,
		MaxConcurrency: maxConcurrency,
	}, nil
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(int64(request.MaxItems)),
	)

	logger.Info("Starting datasource request")

	entityName := request.EntityExternalID
	entityConfig := ValidEntityExternalIDs[request.EntityExternalID]

	// Timeout API calls that take longer than the configured timeout.
	ctx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	iamClient, accountCursor, err := d.GetIamClient(ctx, request)
	if err != nil {
		return nil, err
	}

	// [MemberEntities] For member entities, we need to set the `CollectionID` and `CollectionCursor`.
	memberOf := ValidEntityExternalIDs[entityName].MemberOf
	if memberOf != nil {
		memberReq := &Request{
			EntityExternalID:      *memberOf,
			Auth:                  request.Auth,
			MaxItems:              1,
			EntityConfig:          map[string]*EntityConfig{*memberOf: request.EntityConfig[*memberOf]},
			RequestTimeoutSeconds: request.RequestTimeoutSeconds,
			ResourceAccountRoles:  request.ResourceAccountRoles,
		}

		// If the CollectionCursor is set, use that as the Cursor for the next call to `GetPage`.
		if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
			memberReq.Cursor = &pagination.CompositeCursor[string]{
				Cursor: request.Cursor.CollectionCursor,
			}
		}

		if request.Cursor == nil {
			request.Cursor = &pagination.CompositeCursor[string]{}
		}

		isEmptyLastPage, cursorErr := pagination.UpdateNextCursorFromCollectionAPI(
			ctx,
			request.Cursor,
			func(ctx context.Context, _ *Request) (
				int, string, []map[string]any, *pagination.CompositeCursor[string], *framework.Error,
			) {
				resp, err := d.GetPage(ctx, memberReq)
				if err != nil {
					return 0, "", nil, nil, err
				}

				return resp.StatusCode, resp.RetryAfterHeader, resp.Objects, resp.NextCursor, nil
			},
			memberReq,
			*entityConfig.CollectionAttribute,
		)
		if cursorErr != nil {
			return nil, cursorErr
		}

		// If we hit a page with no `CollectionID` and no
		// `CollectionCursor` we should complete the sync at this point.
		if isEmptyLastPage {
			return &Response{
				StatusCode: http.StatusOK,
			}, nil
		}
	}

	var (
		statusCode int
		fetchErr   error
		nextMarker *string
		objects    []map[string]interface{}
	)

	opts := setOptions(request, entityName, memberOf, d.MaxConcurrency)

	response := &Response{}

	switch entityName {
	case User:
		handler := &UserHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.User](ctx, handler, opts)
	case Group:
		handler := &GroupHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.Group](ctx, handler, opts)
	case GroupMember:
		handler := &GroupMemberHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.User](ctx, handler, opts)
	case Role:
		handler := &RoleHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.Role](ctx, handler, opts)
	case Policy:
		handler := &PolicyHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.Policy](ctx, handler, opts)
	case IdentityProvider:
		handler := &IDPHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.SAMLProviderListEntry](ctx, handler, opts)
	case GroupPolicy:
		handler := &AttachedGroupPoliciesHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.AttachedPolicy](ctx, handler, opts)
	case RolePolicy:
		handler := &AttachedRolePoliciesHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.AttachedPolicy](ctx, handler, opts)
	case UserPolicy:
		handler := &AttachedUserPoliciesHandler{Client: iamClient}
		objects, statusCode, nextMarker, fetchErr = FetchEntities[types.AttachedPolicy](ctx, handler, opts)
	default:
		return nil, &framework.Error{
			Message: fmt.Sprintf("Unsupported entity type: %s", entityName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if fetchErr != nil {
		logger.Error("Datasource responded with an error", fields.ResponseStatusCode(statusCode), fields.SGNLEventTypeError())

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Unable to fetch AWS entity: %s, error: %v.", entityName, fetchErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(fetchErr, request.RequestTimeoutSeconds),
		)
	}

	if statusCode == unhandledStatusCode {
		logger.Error("Datasource responded with an error", fields.ResponseStatusCode(statusCode), fields.SGNLEventTypeError())

		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to fetch AWS entity: %s - Unhandled status code -1", entityName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	response.StatusCode = statusCode

	if nextMarker != nil {
		response.NextCursor = &pagination.CompositeCursor[string]{
			Cursor: nextMarker,
		}
	}

	// [MemberEntities] Set `id`, `memberUniqueIDAttribute` and `memberOfUniqueIDAttribute`.
	if memberOf != nil {
		var (
			memberOfUniqueIDValue                              any
			memberOfUniqueIDAttribute, memberUniqueIDAttribute string
		)

		if request.Cursor == nil || request.Cursor.CollectionID == nil {
			return nil, &framework.Error{
				Message: "Cursor or CollectionID is nil",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		memberOfUniqueIDAttribute = ValidEntityExternalIDs[*memberOf].Identifiers.UniqueName

		switch request.EntityExternalID {
		case GroupPolicy, RolePolicy, UserPolicy:
			memberUniqueIDAttribute = PolicyArn
		case GroupMember:
			memberUniqueIDAttribute = UserID
		default:
			return nil, &framework.Error{
				Message: fmt.Sprintf(
					"Failed during MemberEntity Post-Processing. Invalid Member Entity: %s",
					request.EntityExternalID,
				),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		if request.Cursor != nil && request.Cursor.CollectionID != nil {
			memberOfUniqueIDValue = *request.Cursor.CollectionID
		}

		for idx, member := range objects {
			memberUniqueIDValue, ok := member[memberUniqueIDAttribute].(string)
			if !ok {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse %s field in AWS Member response as string.", memberUniqueIDValue),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			objects[idx][uniqueIDAttribute] = fmt.Sprintf("%s-%s", memberUniqueIDValue, memberOfUniqueIDValue)
			objects[idx][memberUniqueIDAttribute] = memberUniqueIDValue
			objects[idx][memberOfUniqueIDAttribute] = memberOfUniqueIDValue
		}

		if response.NextCursor != nil && response.NextCursor.Cursor != nil {
			request.Cursor.Cursor = response.NextCursor.Cursor
		} else {
			request.Cursor.Cursor = nil
		}

		// If we have a next cursor for either the base collection (Groups, Roles, etc.) or members (Group/Role/etc. Members),
		// encode the cursor for the next page. Otherwise, don't set a cursor as this sync is complete.
		if request.Cursor.Cursor != nil || request.Cursor.CollectionCursor != nil {
			response.NextCursor = request.Cursor
		}
	}

	response.Objects = objects

	// Return the response if -
	// 1) There are no accounts to query
	// 2) All accounts have been queried
	// 3) The request is a `memberOf` entity and the `collectionCursor` has the next marker and account offset.
	if accountCursor == nil ||
		len(request.ResourceAccountRoles) == 0 ||
		(nextMarker == nil && accountCursor.Offset == len(request.ResourceAccountRoles)-1) ||
		memberOf != nil {
		logger.Info("Datasource request completed successfully",
			fields.ResponseStatusCode(response.StatusCode),
			fields.ResponseObjectCount(len(response.Objects)),
			fields.ResponseNextCursor(response.NextCursor),
		)

		return response, nil
	}

	encodedCursor, encodeErr := encodeAccountCursor(accountCursor.Offset, nextMarker, len(request.ResourceAccountRoles))
	if encodeErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error marshalling account cursor: %v", encodeErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if response.NextCursor != nil {
		response.NextCursor.Cursor = &encodedCursor
	} else {
		response.NextCursor = &pagination.CompositeCursor[string]{
			Cursor: &encodedCursor,
		}
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

// setOptions configures and returns an Options struct based on the provided Request,
// entity name, and membership information.
func setOptions(request *Request, entityName string, memberOf *string, maxConcurrent int) *Options {
	opts := &Options{
		InputParams: InputParams{
			MaxItems: func() *int32 {
				maxItems := request.MaxItems

				return &maxItems
			}(),
			PathPrefix: nil,
			Marker:     nil,
		},
		EntityName:         entityName,
		IsMember:           memberOf != nil,
		AccountIDRequested: request.AccountIDRequested,
		MaxConcurrent:      maxConcurrent,
	}

	// Set PathPrefix if it exists in the EntityConfig.
	if request.EntityConfig[request.EntityExternalID] != nil {
		opts.PathPrefix = request.EntityConfig[request.EntityExternalID].PathPrefix
	}

	// Initialize the Cursor if it is nil.
	if request.Cursor == nil {
		request.Cursor = &pagination.CompositeCursor[string]{}
	}

	// Set Marker if the Cursor exists and is not nil.
	if request.Cursor != nil {
		if request.Cursor.Cursor != nil {
			opts.Marker = request.Cursor.Cursor
		}

		opts.UniqueName = request.Cursor.CollectionID
	}

	return opts
}

// FetchEntities is a generic function to fetch AWS entities.
// It retrieves a list of entities using the provided handler,
// converts them to a map format, and returns them along with an HTTP status code
// and a marker for the next set of results if available.
func FetchEntities[T any](
	ctx context.Context,
	handler any,
	opts *Options,
) ([]map[string]interface{}, int, *string, error) {
	var (
		entities   []T
		err        error
		nextMarker *string
	)

	// Check the type of handler to determine if it supports the List operation.
	switch h := handler.(type) {
	case EntityLister[T]:
		entities, nextMarker, err = h.List(ctx, opts)
		if err != nil {
			return nil, statusCodeFromResponseError(err), nil, fmt.Errorf("Failed to list entities: %w", err)
		}
	default:
		return nil, http.StatusBadRequest, nil, fmt.Errorf("Handler does not support List operation")
	}

	objects := make([]map[string]interface{}, len(entities))
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	// Buffered channel to limit the number of concurrent goroutines.
	sem := make(chan struct{}, opts.MaxConcurrent)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Iterate over the entities and process each one concurrently.
	for i, entity := range entities {
		wg.Add(1)

		sem <- struct{}{} // Acquire a slot

		go func(i int, entity T) {
			defer wg.Done()
			defer func() { <-sem }() // Release the slot

			var detailedEntity T

			// Check if the handler supports the Get operation.
			if h, ok := handler.(EntityGetter[T]); ok {
				detailedEntity, err = h.Get(ctx, entity)
				if err != nil {
					errChan <- fmt.Errorf("Failed to get entity: %w", err)

					return
				}
			} else {
				// If the Get operation is not supported, use the original entity.
				detailedEntity = entity
			}

			object, err := EntityToObjects(detailedEntity)
			if err != nil {
				errChan <- fmt.Errorf("Failed to convert entity response to map: %w", err)

				return
			}

			// If AccountId is requested, insert it into the map object.
			if opts.AccountIDRequested {
				if err := ArnToAccountID(&object, opts.EntityName); err != nil {
					select {
					case errChan <- fmt.Errorf("Failed to add AccountID to entity: %w", err):
					default:
					}

					cancel()

					return
				}
			}

			objects[i] = object
		}(i, entity)
	}

	// Wait for all goroutines to finish and close the error channel.
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect the error from the error channel.
	err = <-errChan
	if err != nil {
		return nil, http.StatusInternalServerError, nil, err
	}

	// [IdentityProvider] No pagination support from the AWS side for this entity.
	if opts.EntityName == IdentityProvider {
		var paginationErr *framework.Error

		objects, nextMarker, paginationErr = pagination.PaginateObjects(
			objects, int64(*opts.MaxItems), &pagination.CompositeCursor[string]{
				Cursor: opts.Marker,
			})

		if paginationErr != nil {
			return nil, http.StatusInternalServerError, nil, fmt.Errorf("Failed to paginate objects: %v", paginationErr)
		}
	}

	return objects, http.StatusOK, nextMarker, nil
}

func EntityToObjects(entity interface{}) (map[string]interface{}, error) {
	entityBytes, err := json.Marshal(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entity to JSON: %w", err)
	}

	var object map[string]interface{}
	if err := json.Unmarshal(entityBytes, &object); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}

	return object, nil
}

// arnToAccountId adds the AccountID to the entity map using Entity Arn.
func ArnToAccountID(entity *map[string]interface{}, entityType string) error {
	if entity == nil {
		return fmt.Errorf("entity response is nil")
	}

	arnValue, exists := (*entity)[ValidEntityExternalIDs[entityType].Identifiers.ArnAttribute]
	if !exists {
		return fmt.Errorf("Unable to find Arn in entity")
	}

	arnString, ok := arnValue.(string)
	if !ok {
		return fmt.Errorf("failed to convert Arn to string")
	}

	parsedARN, err := arn.Parse(arnString)
	if err != nil {
		return fmt.Errorf("failed to parse Arn: %v", err)
	}

	(*entity)[AccountID] = parsedARN.AccountID

	return nil
}

// statusCodeFromResponseError returns the status code from an API error.
// If the status code cannot be determined, it returns a sentinel value of -1.
func statusCodeFromResponseError(err error) int {
	var httpResponseErr *awshttp.ResponseError
	if errors.As(err, &httpResponseErr) {
		return httpResponseErr.HTTPStatusCode()
	}

	// Return a sentinel value indicating the status code could not be determined
	return unhandledStatusCode
}

// GetIamClient returns an IAM client for the given request.
// If resource accounts are provided, it assumes the role for the account and returns the client.
// nolint:lll
func (d *Datasource) GetIamClient(ctx context.Context, request *Request) (*iam.Client, *AccountCursor, *framework.Error) {
	// Deep copy of the AWS configuration object ensures that each request operates with
	// its own independent configuration, preventing race conditions.
	awsConfig := d.AWSConfig.Copy()
	awsConfig.Credentials = credentials.NewStaticCredentialsProvider(request.AccessKey, request.SecretKey, "")
	awsConfig.Region = request.Region

	if len(request.ResourceAccountRoles) == 0 {
		return iam.NewFromConfig(awsConfig), nil, nil
	}

	var (
		roleARN   string
		decodeErr error
	)

	accountCursor := &AccountCursor{Offset: 0}

	if request.Cursor != nil {
		if request.Cursor.Cursor != nil {
			accountCursor, decodeErr = decodeAccountCursor(*request.Cursor.Cursor)
			if decodeErr != nil {
				return nil, nil, &framework.Error{
					Message: fmt.Sprintf("Error decoding cursor: %v", decodeErr),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			// Set the cursor to the next marker.
			request.Cursor.Cursor = accountCursor.NextMarker
		}

		// This condition is true for GroupMember, GroupPolicy, RolePolicy, and UserPolicy entities.
		if request.Cursor.CollectionCursor != nil {
			accountCursor, decodeErr = decodeAccountCursor(*request.Cursor.CollectionCursor)
			if decodeErr != nil {
				return nil, nil, &framework.Error{
					Message: fmt.Sprintf("Error decoding collection cursor: %v", decodeErr),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}
		}
	}

	roleARN = request.ResourceAccountRoles[accountCursor.Offset]
	stsClient := sts.NewFromConfig(awsConfig)
	// Assume the role.
	assumeRoleOutput, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn: aws.String(roleARN),
		// Each call to AssumeRole returns temporary security credentials, which are tied to the session name.
		// If you need to differentiate between credentials from different roles, using unique session names helps maintain clarity.
		RoleSessionName: aws.String(fmt.Sprintf("%s-%d", SessionName, accountCursor.Offset)),
	})
	if err != nil {
		return nil, nil, &framework.Error{
			Message: fmt.Sprintf("Failed to assume role: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	configWithAssumedRole := d.AWSConfig.Copy()
	configWithAssumedRole.Credentials = credentials.NewStaticCredentialsProvider(
		*assumeRoleOutput.Credentials.AccessKeyId,
		*assumeRoleOutput.Credentials.SecretAccessKey,
		*assumeRoleOutput.Credentials.SessionToken,
	)
	configWithAssumedRole.Region = request.Region

	return iam.NewFromConfig(configWithAssumedRole), accountCursor, nil
}

func decodeAccountCursor(cursor string) (*AccountCursor, error) {
	var accountCursor AccountCursor

	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(decodedCursor, &accountCursor); err != nil {
		return nil, err
	}

	return &accountCursor, nil
}

func encodeAccountCursor(accountOffset int, nextMarker *string, numOfResourceAccounts int) (string, error) {
	nextAccountCursor := AccountCursor{
		Offset:     accountOffset,
		NextMarker: nextMarker,
	}

	// There are no markers left to query from the current account
	// and there are more accounts to query from.
	// Increment the accountOffset by 1 to query from the next account.
	if nextMarker == nil && accountOffset+1 < numOfResourceAccounts {
		nextAccountCursor.Offset = accountOffset + 1
	}

	// Marshal the next account cursor and encode it.
	cursorBytes, marshalErr := json.Marshal(nextAccountCursor)
	if marshalErr != nil {
		return "", marshalErr
	}

	return base64.StdEncoding.EncodeToString(cursorBytes), nil
}
