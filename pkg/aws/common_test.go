// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package aws_test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsMiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/aws/smithy-go/middleware"
	framework "github.com/sgnl-ai/adapter-framework"
	aws_adapter "github.com/sgnl-ai/adapters/pkg/aws"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

var (
	defaultTimeout       = 120
	internalError        = "some-internal-error"
	validAuthCredentials = &framework.DatasourceAuthCredentials{
		Basic: &framework.BasicAuthCredentials{
			Username: "user",
			Password: "pass",
		},
	}

	validCommonConfig = &aws_adapter.Config{
		Region: "us-west-2",
		EntityConfig: map[string]*aws_adapter.EntityConfig{
			"User": {
				PathPrefix: testutil.GenPtr("/"),
			},
			"Role": {
				PathPrefix: testutil.GenPtr("/"),
			},
			"Group": {
				PathPrefix: testutil.GenPtr("/"),
			},
			"Policy": {
				PathPrefix: testutil.GenPtr("/"),
			},
		},
	}

	validCommonConfigWithAccounts = &aws_adapter.Config{
		Region: "us-west-2",
		ResourceAccountRoles: []string{
			"arn:aws:iam::123456789012:role/role-name",
			"arn:aws:iam::111111111111:role/role-name",
		},
		EntityConfig: map[string]*aws_adapter.EntityConfig{
			"User": {
				PathPrefix: testutil.GenPtr("/"),
			},
			"Role": {
				PathPrefix: testutil.GenPtr("/"),
			},
			"Group": {
				PathPrefix: testutil.GenPtr("/"),
			},
			"Policy": {
				PathPrefix: testutil.GenPtr("/"),
			},
		},
	}
	assumeRoleOperation = "AssumeRole"
)

func SetupTestConfig(ctx context.Context, withAPIOptionsFuncs ...func(*middleware.Stack) error) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("us-west-2"),
		config.WithAPIOptions(withAPIOptionsFuncs),
	)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ProvideAWSTestClient creates a new AWS client for testing purposes.
func ProvideAWSTestClient(cfg *aws.Config) (framework.Adapter[aws_adapter.Config], error) {
	client, err := aws_adapter.NewClient(http.DefaultClient, cfg, 1)
	if err != nil {
		return nil, fmt.Errorf("error creating client to query datasource: %v", err)
	}

	return aws_adapter.NewAdapter(client), nil
}

// Mocker is a middleware that helps to mock the AWS API calls.
//
// References:
//   - https://aws.github.io/aws-sdk-go-v2/docs/middleware/
//   - https://dev.to/aws-builders/testing-with-aws-sdk-for-go-v2-without-interface-mocks-55de
func Mocker(stack *middleware.Stack) error {
	// Initialize
	// Prepares the input, and sets any default parameters as needed.
	err := stack.Initialize.Add(
		middleware.InitializeMiddlewareFunc(
			"GetInputOptions",
			getInputOptions,
		), middleware.Before,
	)

	if err != nil {
		return err
	}

	// Finalize
	// Final message preparation, including retries and authentication (SigV4 signing).
	err = stack.Finalize.Add(
		middleware.FinalizeMiddlewareFunc(
			"AssumeRoleListerGetterMocker",
			func(ctx context.Context, _ middleware.FinalizeInput, _ middleware.FinalizeHandler,
			) (middleware.FinalizeOutput, middleware.Metadata, error) {
				operationName := awsMiddleware.GetOperationName(ctx)

				if operationName == assumeRoleOperation {
					return middleware.FinalizeOutput{
						Result: &sts.AssumeRoleOutput{
							Credentials: &types.Credentials{
								AccessKeyId:     aws.String("mockAccessKeyId"),
								SecretAccessKey: aws.String("mockSecretAccessKey"),
								SessionToken:    aws.String("mockSessionToken"),
							},
						},
					}, middleware.Metadata{}, nil
				}

				// Retrieve AWS SDK options from the middleware stack.
				options, ok := middleware.GetStackValue(ctx, aws_adapter.Options{}).(*aws_adapter.Options)
				if !ok {
					return middleware.FinalizeOutput{}, middleware.Metadata{}, nil
				}

				// Set up mocks based on the operation name and options.
				return setupMocks(operationName, *options)
			},
		),
		middleware.Before,
	)

	return err
}

func EmptyMocker(stack *middleware.Stack) error {
	// Finalize
	// Final message preparation, including retries and authentication (SigV4 signing).
	err := stack.Finalize.Add(
		middleware.FinalizeMiddlewareFunc(
			"AssumeRoleListerGetterMocker",
			func(ctx context.Context, _ middleware.FinalizeInput, _ middleware.FinalizeHandler,
			) (middleware.FinalizeOutput, middleware.Metadata, error) {
				operationName := awsMiddleware.GetOperationName(ctx)

				if operationName == assumeRoleOperation {
					return middleware.FinalizeOutput{
						Result: &sts.AssumeRoleOutput{
							Credentials: &types.Credentials{
								AccessKeyId:     aws.String("mockAccessKeyId"),
								SecretAccessKey: aws.String("mockSecretAccessKey"),
								SessionToken:    aws.String("mockSessionToken"),
							},
						},
					}, middleware.Metadata{}, nil
				}

				// Set up mocks based on the operation name and options.
				return setupEmptyMocks(operationName)
			},
		),
		middleware.Before,
	)

	return err
}

// PolicyDatasetMocker is a middleware that helps to mock the AWS API calls for benchmarking purposes.
func PolicyDatasetMocker(stack *middleware.Stack) error {
	// Initialize
	// Prepares the input, and sets any default parameters as needed.
	err := stack.Initialize.Add(
		middleware.InitializeMiddlewareFunc(
			"GetInputOptions",
			getInputOptions,
		), middleware.Before,
	)
	if err != nil {
		return err
	}

	// Finalize
	// Final message preparation, including retries and authentication (SigV4 signing).
	err = stack.Finalize.Add(
		middleware.FinalizeMiddlewareFunc(
			"AssumeRoleAndListerGetterMocker",
			func(ctx context.Context, _ middleware.FinalizeInput, _ middleware.FinalizeHandler,
			) (middleware.FinalizeOutput, middleware.Metadata, error) {
				operationName := awsMiddleware.GetOperationName(ctx)

				if operationName == assumeRoleOperation {
					return middleware.FinalizeOutput{
						Result: &sts.AssumeRoleOutput{
							Credentials: &types.Credentials{
								AccessKeyId:     aws.String("mockAccessKeyId"),
								SecretAccessKey: aws.String("mockSecretAccessKey"),
								SessionToken:    aws.String("mockSessionToken"),
							},
						},
					}, middleware.Metadata{}, nil
				}

				// Retrieve AWS SDK options from the middleware stack.
				options, ok := middleware.GetStackValue(ctx, aws_adapter.Options{}).(*aws_adapter.Options)
				if !ok {
					return middleware.FinalizeOutput{}, middleware.Metadata{}, nil
				}

				// Set up mocks based on the operation name and options.
				return setupPolicyDataset(operationName, *options)
			},
		),
		middleware.Before,
	)

	return err
}

// paginate paginates the given items based on the provided cursor and page size.
func paginate[T any](items []T, cursor *string, pageSize int) ([]T, *string) {
	if len(items) == 0 {
		return []T{}, nil
	}

	start := 0

	if cursor != nil {
		startIndex, _ := strconv.Atoi(*cursor)
		start = startIndex
	}

	end := start + pageSize

	if start >= len(items) {
		return []T{}, nil
	}

	// Ensure the end index does not exceed the bounds
	if end > len(items) {
		end = len(items)
	}

	// Determine the next cursor, if any
	var nextCursor *string

	if end < len(items) {
		nextCursorStr := strconv.Itoa(end)
		nextCursor = &nextCursorStr
	}

	return items[start:end], nextCursor
}
