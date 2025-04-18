// Copyright 2025 SGNL.ai, Inc.

package awss3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

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

	return &Datasource{
		AWSConfig: awsConfig,
		Client:    client,
	}, nil
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	entityName := request.EntityExternalID

	// Timeout API calls that take longer than the configured timeout.
	ctx, cancel := context.WithTimeout(ctx, time.Duration(request.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	// Deep copy of the AWS configuration object ensures that each request operates with
	// its own independent configuration, preventing race conditions.
	awsConfig := d.AWSConfig.Copy()
	awsConfig.Credentials = credentials.NewStaticCredentialsProvider(request.AccessKey, request.SecretKey, "")
	awsConfig.Region = request.Region

	handler := &S3Handler{Client: s3.NewFromConfig(awsConfig)}

	// Create the object key using entity name and path prefix
	objectKey := GetObjectKeyFromRequest(request)

	fileBytes, err := handler.GetFile(ctx, request.Bucket, objectKey)
	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to fetch entity from AWS S3: %s, error: %v.", entityName, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	// TODO: Converting bytes to csv reader would consume additional memory.
	// Optimize this to read s3 object body, which implements io.Reader interface, to csv.NewReader directly.

	if fileBytes == nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("The file for entity %s is empty.", entityName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	// Convert bytes to csv reader
	start := int64(1)
	if request.Cursor != nil && request.Cursor.Cursor != nil {
		start = *request.Cursor.Cursor
	}

	objects, hasNext, err := CSVBytesToPage(fileBytes, start, request.PageSize, request.AttributeConfig)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to get %s entity objects from the CSV file: %v.", entityName, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	response := &Response{
		StatusCode: 200,
		Objects:    objects,
	}

	if !hasNext {
		return response, nil
	}

	nextPage := start + request.PageSize
	response.NextCursor = &pagination.CompositeCursor[int64]{Cursor: &nextPage}

	return response, nil
}

// httpResponseFromError returns a awshttp.ResponseError from an SDK error.
// If the error cannot be parsed to an awshttp.ResponseError, it returns the original error object.
func httpResponseFromError(err error) (*awshttp.ResponseError, error) {
	var httpResponseErr *awshttp.ResponseError
	if !errors.As(err, &httpResponseErr) {
		return nil, err
	}

	return httpResponseErr, nil
}

func GetObjectKeyFromRequest(request *Request) string {
	return filepath.Join(
		filepath.Clean(request.PathPrefix),
		filepath.Clean(fmt.Sprintf("%s.%s", request.EntityExternalID, request.FileType)),
	)
}
