// Copyright 2025 SGNL.ai, Inc.

package awss3

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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

const MaxBytesToProcessPerPage = 200 * MaxCSVRowSizeBytes // 200MB

// BOM (Byte Order Mark) patterns for different encodings.
var (
	UTF8BOM    = []byte{0xEF, 0xBB, 0xBF}
	UTF16LEBOM = []byte{0xFF, 0xFE}
	UTF16BEBOM = []byte{0xFE, 0xFF}
	UTF32LEBOM = []byte{0xFF, 0xFE, 0x00, 0x00}
	UTF32BEBOM = []byte{0x00, 0x00, 0xFE, 0xFF}
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

func stripBOM(reader *bufio.Reader) (bomLength int, err error) {
	peekedBytes, peekErr := reader.Peek(len(UTF32LEBOM))
	if peekErr != nil && peekErr != io.EOF && peekErr != bufio.ErrBufferFull {
		return 0, fmt.Errorf("error peeking for BOM: %w", peekErr)
	}

	identifiedBomLength := 0
	if bytes.HasPrefix(peekedBytes, UTF32LEBOM) {
		identifiedBomLength = len(UTF32LEBOM)
	} else if bytes.HasPrefix(peekedBytes, UTF32BEBOM) {
		identifiedBomLength = len(UTF32BEBOM)
	} else if bytes.HasPrefix(peekedBytes, UTF8BOM) {
		identifiedBomLength = len(UTF8BOM)
	} else if bytes.HasPrefix(peekedBytes, UTF16LEBOM) {
		identifiedBomLength = len(UTF16LEBOM)
	} else if bytes.HasPrefix(peekedBytes, UTF16BEBOM) {
		identifiedBomLength = len(UTF16BEBOM)
	}

	if identifiedBomLength > 0 {
		_, discardErr := reader.Discard(identifiedBomLength)
		if discardErr != nil {
			return 0, fmt.Errorf("error discarding BOM (length %d): %w", identifiedBomLength, discardErr)
		}

		return identifiedBomLength, nil
	}

	return 0, nil
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

	fileSize, err := handler.GetFileSize(ctx, request.Bucket, objectKey)
	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to fetch entity from AWS S3: %s, error: %v.", entityName, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}, customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds))
	}

	if fileSize == 0 {
		return nil, &framework.Error{
			Message: fmt.Sprintf("The file for entity %s is empty.", entityName),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
		}
	}

	s3HeaderStreamOutput, err := handler.GetObjectStream(ctx, request.Bucket, objectKey, nil)
	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to fetch entity from AWS S3: %s, error: %v.", entityName, err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}, customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds))
	}
	defer s3HeaderStreamOutput.Body.Close()

	headerBufReader := bufio.NewReader(s3HeaderStreamOutput.Body)

	bomLength, bomErr := stripBOM(headerBufReader)
	if bomErr != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to fetch entity from AWS S3: %s, error processing BOM: %v", entityName, bomErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}, customerror.WithRequestTimeoutMessage(bomErr, request.RequestTimeoutSeconds))
	}

	parsedHeaders, bytesReadForHeaderLine, err := CSVHeaders(headerBufReader)
	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Unable to parse CSV file headers: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}, customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds))
	}

	s3HeaderStreamOutput.Body.Close()

	firstDataByteOffset := int64(bomLength) + bytesReadForHeaderLine

	var startBytePos int64
	if request.Cursor != nil && request.Cursor.Cursor != nil {
		startBytePos = *request.Cursor.Cursor
	} else {
		startBytePos = firstDataByteOffset
	}

	if startBytePos >= fileSize {
		return &Response{StatusCode: 200, Objects: []map[string]any{}}, nil
	}

	rangeHeaderForData := fmt.Sprintf("bytes=%d-", startBytePos)
	s3DataStreamOutput, err := handler.GetObjectStream(ctx, request.Bucket, objectKey, &rangeHeaderForData)

	if err != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to fetch entity from AWS S3: %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}, customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds))
	}

	defer s3DataStreamOutput.Body.Close()

	dataBufReader := bufio.NewReader(s3DataStreamOutput.Body)

	var (
		objects                 []map[string]any
		hasNext                 bool
		bytesReadFromDataStream int64
		processErr              error
	)

	objects, bytesReadFromDataStream, hasNext, processErr = StreamingCSVToPage(
		dataBufReader,
		parsedHeaders,
		request.PageSize,
		request.AttributeConfig,
		MaxBytesToProcessPerPage,
	)

	if processErr != nil {
		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Failed to fetch entity from AWS S3: %s, error: %v.", entityName, processErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}, customerror.WithRequestTimeoutMessage(processErr, request.RequestTimeoutSeconds))
	}

	response := &Response{
		StatusCode: 200,
		Objects:    objects,
	}

	if hasNext {
		nextBytePos := startBytePos + bytesReadFromDataStream
		if nextBytePos < fileSize {
			response.NextCursor = &pagination.CompositeCursor[int64]{Cursor: &nextBytePos}
		}
	}

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
