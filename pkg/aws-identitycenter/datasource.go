package awsidentitycenter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
		input := &identitystore.ListGroupMembershipsInput{
			IdentityStoreId: aws.String(request.IdentityStoreID),
			MaxResults:      aws.Int32(request.MaxResults),
		}
		if request.Cursor != nil && request.Cursor.Cursor != nil {
			input.NextToken = request.Cursor.Cursor
		}

		out, err2 := client.ListGroupMemberships(ctx, input)
		if err2 != nil {
			err = err2
			statusCode = statusCodeFromResponseError(err2)
			break
		}
		objects = make([]map[string]any, len(out.GroupMemberships))
		for i, mship := range out.GroupMemberships {
			m, convErr := awsadapter.EntityToObjects(mship)
			if convErr != nil {
				err = fmt.Errorf("failed to convert entity response to map: %w", convErr)
				statusCode = http.StatusInternalServerError
				break
			}
			objects[i] = m
		}
		nextToken = out.NextToken
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
