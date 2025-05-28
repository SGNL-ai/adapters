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

	// Get file metadata first
	fileSize, err := handler.GetFileSize(ctx, request.Bucket, objectKey)
	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Unable to access file for entity %s from AWS S3: %v.", entityName, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	}

	if fileSize == 0 {
		return nil, &framework.Error{
			Message: fmt.Sprintf("The CSV file for entity %s is empty.", entityName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	// Determine pagination start position
	start := int64(1) // Start at row 1 (skip header row)
	if request.Cursor != nil && request.Cursor.Cursor != nil {
		start = *request.Cursor.Cursor
	}

	// Decide whether to use streaming or traditional approach based on file size
	// For files larger than 10MB, use streaming to avoid memory issues
	const streamingThreshold = 10 * 1024 * 1024 // 10MB

	var objects []map[string]any
	var hasNext bool
	var processErr error

	if fileSize > streamingThreshold {
		// Use streaming approach for large files
		objects, hasNext, processErr = d.processLargeFileStreaming(
			ctx, handler, request.Bucket, objectKey, fileSize, start, request.PageSize, request.AttributeConfig,
		)
	} else {
		// Use traditional approach for smaller files (backward compatibility)
		objects, hasNext, processErr = d.processSmallFileTraditional(
			ctx, handler, request.Bucket, objectKey, start, request.PageSize, request.AttributeConfig,
		)
	}

	if processErr != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Unable to process CSV data for entity %s: %v.", entityName, processErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(processErr, request.RequestTimeoutSeconds),
		)
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

// processLargeFileStreaming handles large files using streaming approach
func (d *Datasource) processLargeFileStreaming(
	ctx context.Context,
	handler *S3Handler,
	bucket, key string,
	fileSize, start, pageSize int64,
	attrConfig []*framework.AttributeConfig,
) ([]map[string]any, bool, error) {

	// First, get the CSV headers by reading a small chunk from the beginning
	headerChunk, err := handler.GetHeaderChunk(ctx, bucket, key)
	if err != nil {
		return nil, false, fmt.Errorf("unable to read CSV file headers: %v", err)
	}

	headers, err := CSVHeaders(headerChunk)
	if err != nil {
		return nil, false, fmt.Errorf("unable to parse CSV file headers: %v", err)
	}

	if len(headers) == 0 {
		return nil, false, fmt.Errorf("CSV file does not contain valid column headers")
	}

	// Process the file using streaming
	objects, hasNext, err := StreamingCSVToPage(
		handler, ctx, bucket, key, fileSize, headers, start, pageSize, attrConfig,
	)
	if err != nil {
		return nil, false, fmt.Errorf("unable to process CSV file data: %v", err)
	}

	return objects, hasNext, nil
}

// processSmallFileTraditional handles small files using the original approach
func (d *Datasource) processSmallFileTraditional(
	ctx context.Context,
	handler *S3Handler,
	bucket, key string,
	start, pageSize int64,
	attrConfig []*framework.AttributeConfig,
) ([]map[string]any, bool, error) {

	// Use the original GetFile method for small files
	fileBytes, err := handler.GetFile(ctx, bucket, key)
	if err != nil {
		return nil, false, fmt.Errorf("unable to read CSV file: %v", err)
	}

	if fileBytes == nil {
		return nil, false, fmt.Errorf("CSV file is empty or corrupted")
	}

	// Use the original CSV processing method
	objects, hasNext, err := CSVBytesToPage(fileBytes, start, pageSize, attrConfig)
	if err != nil {
		return nil, false, fmt.Errorf("unable to process CSV file data: %v", err)
	}

	return objects, hasNext, nil
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
