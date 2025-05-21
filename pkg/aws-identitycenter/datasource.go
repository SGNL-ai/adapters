package awsidentitycenter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	awsadapter "github.com/sgnl-ai/adapters/pkg/aws"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

const (
	PermissionSet   string = "PermissionSet"
	User            string = "User"
	Group           string = "Group"
	GroupMembership string = "GroupMembership"

	unhandledStatusCode int = -1
)

var validEntities = map[string]struct{}{
	PermissionSet:   {},
	User:            {},
	Group:           {},
	GroupMembership: {},
}

type Datasource struct {
	Client    *http.Client
	AWSConfig *aws.Config
}

// NewClient returns a Client to query the datasource.
func NewClient(client *http.Client, awsConfig *aws.Config) (Client, error) {
	if awsConfig == nil {
		cfg, err := aws_config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		awsConfig = &cfg
	}

	return &Datasource{AWSConfig: awsConfig, Client: client}, nil
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	if _, ok := validEntities[request.EntityExternalID]; !ok {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Unsupported entity type: %s", request.EntityExternalID),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	cfg := d.AWSConfig.Copy()
	cfg.Credentials = credentials.NewStaticCredentialsProvider(request.AccessKey, request.SecretKey, "")
	cfg.Region = request.Region

	resp := &Response{}
	var (
		objects    []map[string]any
		nextToken  *string
		statusCode int
		err        error
	)

	switch request.EntityExternalID {
	case PermissionSet:
		client := ssoadmin.NewFromConfig(cfg)
		input := &ssoadmin.ListPermissionSetsInput{
			InstanceArn: aws.String(request.InstanceARN),
			MaxResults:  aws.Int32(request.MaxResults),
		}
		if request.Cursor != nil && request.Cursor.Cursor != nil {
			input.NextToken = request.Cursor.Cursor
		}

		out, err2 := client.ListPermissionSets(ctx, input)
		if err2 != nil {
			err = err2
			statusCode = statusCodeFromResponseError(err2)
			break
		}
		objects = make([]map[string]any, len(out.PermissionSets))
		for i, arn := range out.PermissionSets {
			objects[i] = map[string]any{"Arn": arn}
		}
		nextToken = out.NextToken
		statusCode = http.StatusOK
	case User:
		client := identitystore.NewFromConfig(cfg)
		input := &identitystore.ListUsersInput{
			IdentityStoreId: aws.String(request.IdentityStoreID),
			MaxResults:      aws.Int32(request.MaxResults),
		}
		if request.Cursor != nil && request.Cursor.Cursor != nil {
			input.NextToken = request.Cursor.Cursor
		}

		// Add this log:
		fmt.Printf("DEBUG: ListUsersInput: %+v\n", input)

		out, err2 := client.ListUsers(ctx, input)
		if err2 != nil {
			err = err2
			statusCode = statusCodeFromResponseError(err2)
			// Add this log:
			fmt.Printf("ERROR: AWS ListUsers API call failed. Status: %d, Error: %s\n", statusCode, err2.Error())
			// For even more detail, you might try:
			// fmt.Printf("ERROR: AWS ListUsers API call failed. Detailed Error: %#v\n", err2)
			break
		}
		objects = make([]map[string]any, len(out.Users))
		for i, u := range out.Users {
			m, convErr := awsadapter.EntityToObjects(u)
			if convErr != nil {
				err = fmt.Errorf("failed to convert entity response to map: %w", convErr)
				statusCode = http.StatusInternalServerError
				break
			}
			objects[i] = m
		}
		nextToken = out.NextToken
		statusCode = http.StatusOK
	case Group:
		client := identitystore.NewFromConfig(cfg)
		input := &identitystore.ListGroupsInput{
			IdentityStoreId: aws.String(request.IdentityStoreID),
			MaxResults:      aws.Int32(request.MaxResults),
		}
		if request.Cursor != nil && request.Cursor.Cursor != nil {
			input.NextToken = request.Cursor.Cursor
		}

		out, err2 := client.ListGroups(ctx, input)
		if err2 != nil {
			err = err2
			statusCode = statusCodeFromResponseError(err2)
			break
		}
		objects = make([]map[string]any, len(out.Groups))
		for i, g := range out.Groups {
			m, convErr := awsadapter.EntityToObjects(g)
			if convErr != nil {
				err = fmt.Errorf("failed to convert entity response to map: %w", convErr)
				statusCode = http.StatusInternalServerError
				break
			}
			objects[i] = m
		}
		nextToken = out.NextToken
		statusCode = http.StatusOK
	case GroupMembership:
		client := identitystore.NewFromConfig(cfg)
		objects = []map[string]any{}

		// We need to first get all groups, then for each group get the memberships
		// If we have a cursor that includes a GroupId, use it
		var currentGroupId, currentGroupNextToken *string
		if request.Cursor != nil && request.Cursor.Cursor != nil {
			// Parse the composite cursor format: "groupId:nextToken"
			cursorStr := *request.Cursor.Cursor
			// Split by first colon only
			parts := strings.SplitN(cursorStr, ":", 2)
			if len(parts) >= 1 {
				currentGroupId = aws.String(parts[0])
				if len(parts) >= 2 && parts[1] != "" {
					currentGroupNextToken = aws.String(parts[1])
				}
			}
		}

		// If we don't have a currentGroupId, we need to list all groups first
		if currentGroupId == nil {
			groupsInput := &identitystore.ListGroupsInput{
				IdentityStoreId: aws.String(request.IdentityStoreID),
				MaxResults:      aws.Int32(100), // Get a reasonable batch of groups
			}

			groupsOut, err2 := client.ListGroups(ctx, groupsInput)
			if err2 != nil {
				err = err2
				statusCode = statusCodeFromResponseError(err2)
				break
			}

			if len(groupsOut.Groups) == 0 {
				// No groups to process
				statusCode = http.StatusOK
				break
			}

			// Get the first group and use it
			currentGroupId = groupsOut.Groups[0].GroupId

			// Store the remaining groups and next token for later processing
			if len(groupsOut.Groups) > 1 || groupsOut.NextToken != nil {
				// TODO: Store the list of remaining groups in a more persistent way
				// for now we'll just use the first group and rely on pagination
			}
		}

		// Now fetch memberships for the current group
		input := &identitystore.ListGroupMembershipsInput{
			IdentityStoreId: aws.String(request.IdentityStoreID),
			GroupId:         currentGroupId,
			MaxResults:      aws.Int32(request.MaxResults),
		}

		if currentGroupNextToken != nil {
			input.NextToken = currentGroupNextToken
		}

		out, err2 := client.ListGroupMemberships(ctx, input)
		if err2 != nil {
			err = err2
			statusCode = statusCodeFromResponseError(err2)
			break
		}

		// Process the memberships
		objects = make([]map[string]any, len(out.GroupMemberships))
		for i, mship := range out.GroupMemberships {
			m, convErr := awsadapter.EntityToObjects(mship)
			if convErr != nil {
				err = fmt.Errorf("failed to convert entity response to map: %w", convErr)
				statusCode = http.StatusInternalServerError
				break
			}

			// Add the GroupId to each membership since it's not in the response
			m["GroupId"] = *currentGroupId

			// Extract UserId from MemberId object if it exists
			if memberIdObj, ok := m["MemberId"].(map[string]interface{}); ok {
				if userId, ok := memberIdObj["Value"].(string); ok {
					// Replace the complex MemberId object with just the UserId string
					m["MemberId"] = userId
				}
			}

			objects[i] = m
		}

		// Set the next cursor using the composite format
		if out.NextToken != nil {
			// Continue with the same group but next page
			nextToken = aws.String(fmt.Sprintf("%s:%s", *currentGroupId, *out.NextToken))
		} else {
			// TODO: We need to implement proper pagination across groups
			// For now, we'll just return the memberships for the first group
		}

		statusCode = http.StatusOK
	}

	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to fetch AWS Identity Center entity: %s, error: %v.", request.EntityExternalID, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}, customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds))
	}

	resp.StatusCode = statusCode
	resp.Objects = objects
	if nextToken != nil {
		resp.NextCursor = &pagination.CompositeCursor[string]{Cursor: nextToken}
	}

	return resp, nil
}

func statusCodeFromResponseError(err error) int {
	var httpResponseErr *awshttp.ResponseError
	if errors.As(err, &httpResponseErr) {
		return httpResponseErr.HTTPStatusCode()
	}

	return unhandledStatusCode
}
