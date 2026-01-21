// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const (
	emptyCSVFileCode             = 800
	headersOnlyCSVFileCode       = 801
	largeCSVFileCode             = -300
	largeFileHeaderIndicatorCode = -301
)

type mockS3Middleware struct {
	headStatusCode int
	getStatusCode  int
}

func (m *mockS3Middleware) ID() string {
	return "MockS3Middleware"
}

func (m *mockS3Middleware) HandleSerialize(
	ctx context.Context,
	in middleware.SerializeInput,
	next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput,
	metadata middleware.Metadata,
	err error,
) {
	switch params := in.Parameters.(type) {
	case *s3.HeadObjectInput:
		return m.mockHeadObject(ctx, params, next)
	case *s3.GetObjectInput:
		return m.mockGetObject(ctx, params, next)
	default:
		return next.HandleSerialize(ctx, in)
	}
}

func (m *mockS3Middleware) mockHeadObject(
	_ context.Context,
	_ *s3.HeadObjectInput,
	_ middleware.SerializeHandler,
) (
	out middleware.SerializeOutput,
	metadata middleware.Metadata,
	err error,
) {
	if m.headStatusCode == largeFileHeaderIndicatorCode {
		largeSize := int64(11 * 1024 * 1024)
		out.Result = &s3.HeadObjectOutput{
			ContentLength: &largeSize,
			ContentType:   aws.String("text/csv"),
			ETag:          aws.String("\"f8a7b3f9be0e4c3d2e1a0b9c8d7e6f5\""),
			LastModified:  aws.Time(time.Now()),
			VersionId:     aws.String("3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY"),
			Metadata:      map[string]string{"example-metadata-key": "example-metadata-value"},
		}

		return
	}

	switch m.headStatusCode {
	case http.StatusOK:
		out.Result = &s3.HeadObjectOutput{
			ContentLength: aws.Int64(int64(len(validCSVData))),
			ContentType:   aws.String("text/csv"),
			ETag:          aws.String("\"f8a7b3f9be0e4c3d2e1a0b9c8d7e6f5\""),
			LastModified:  aws.Time(time.Now()),
			VersionId:     aws.String("3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY"),
			Metadata:      map[string]string{"example-metadata-key": "example-metadata-value"},
		}
	case http.StatusMovedPermanently:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusMovedPermanently}},
			Err: errors.New("permanent redirect: The bucket you are attempting to access must be addressed " +
				"using the specified endpoint"),
		}
	case http.StatusForbidden:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusForbidden}},
			Err:      errors.New("AccessDenied: Access Denied"),
		}
	case http.StatusNotFound:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
			Err:      errors.New("not found: The specified key does not exist"),
		}
	default:
		err = fmt.Errorf("mockHeadObject: unexpected headStatusCode %d", m.headStatusCode)
	}

	return
}

func (m *mockS3Middleware) mockGetObject(
	_ context.Context,
	getObjectInput *s3.GetObjectInput,
	_ middleware.SerializeHandler,
) (
	out middleware.SerializeOutput,
	metadata middleware.Metadata,
	err error,
) {
	var fullDataString string

	switch m.getStatusCode {
	case http.StatusMovedPermanently:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusMovedPermanently}},
			Err: errors.New("permanent redirect: The bucket you are attempting to access must be addressed " +
				"using the specified endpoint"),
		}

		return
	case http.StatusForbidden:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusForbidden}},
			Err:      errors.New("access denied: Access Denied"),
		}

		return
	case http.StatusNotFound:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
			Err:      errors.New("no such key: The specified key does not exist"),
		}

		return
	case http.StatusRequestedRangeNotSatisfiable:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusRequestedRangeNotSatisfiable}},
			Err:      errors.New("range not satisfiable"),
		}

		return
	}

	switch m.getStatusCode {
	case emptyCSVFileCode:
		fullDataString = ""
	case headersOnlyCSVFileCode:
		fullDataString = headersOnlyCSVData
	case largeCSVFileCode:
		fullDataString = generateLargeCSVData()
	case -200:
		fullDataString = corruptCSVData
	case http.StatusOK:
		fullDataString = validCSVData
	default:
		err = fmt.Errorf("mockGetObject: unexpected getStatusCode %d for body generation", m.getStatusCode)

		return
	}

	startByte := int64(0)
	servedData := fullDataString

	if getObjectInput.Range != nil && *getObjectInput.Range != "" {
		rangeStr := *getObjectInput.Range
		if strings.HasPrefix(rangeStr, "bytes=") {
			rangeValue := strings.TrimPrefix(rangeStr, "bytes=")
			parts := strings.SplitN(rangeValue, "-", 2)

			if parsedStart, parseErr := strconv.ParseInt(parts[0], 10, 64); parseErr == nil {
				startByte = parsedStart
			}

			if startByte < 0 {
				startByte = 0
			}

			if startByte >= int64(len(fullDataString)) {
				servedData = ""
			} else {
				servedData = fullDataString[startByte:]
			}
		}
	}

	out.Result = &s3.GetObjectOutput{
		Body:          io.NopCloser(strings.NewReader(servedData)),
		ContentLength: aws.Int64(int64(len(servedData))),
		ContentType:   aws.String("text/csv"),
		ETag:          aws.String("\"f8a7b3f9be0e4c3d2e1a0b9c8d7e6f5\""),
		LastModified:  aws.Time(time.Now()),
		VersionId:     aws.String("3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY"),
		Metadata:      map[string]string{"example-metadata-key": "example-metadata-value"},
	}

	return
}

func mockS3Config(headStatusCode, getStatusCode int) *aws.Config {
	mockMiddleware := &mockS3Middleware{
		headStatusCode: headStatusCode,
		getStatusCode:  getStatusCode,
	}

	return &aws.Config{
		Region: "us-west-2",
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				return s.Serialize.Add(mockMiddleware, middleware.After)
			},
		},
	}
}

// rangeTrackingMiddleware wraps mockS3Middleware and captures range headers for testing.
type rangeTrackingMiddleware struct {
	mockS3Middleware
	CapturedRanges []string
}

func (m *rangeTrackingMiddleware) HandleSerialize(
	ctx context.Context,
	in middleware.SerializeInput,
	next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput,
	metadata middleware.Metadata,
	err error,
) {
	if getObjectInput, ok := in.Parameters.(*s3.GetObjectInput); ok {
		rangeHeader := ""
		if getObjectInput.Range != nil {
			rangeHeader = *getObjectInput.Range
		}

		m.CapturedRanges = append(m.CapturedRanges, rangeHeader)
	}

	return m.mockS3Middleware.HandleSerialize(ctx, in, next)
}

func newRangeTrackingConfig(headStatusCode, getStatusCode int) (*aws.Config, *rangeTrackingMiddleware) {
	trackingMiddleware := &rangeTrackingMiddleware{
		mockS3Middleware: mockS3Middleware{
			headStatusCode: headStatusCode,
			getStatusCode:  getStatusCode,
		},
		CapturedRanges: []string{},
	}

	config := &aws.Config{
		Region: "us-west-2",
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				return s.Serialize.Add(trackingMiddleware, middleware.After)
			},
		},
	}

	return config, trackingMiddleware
}
